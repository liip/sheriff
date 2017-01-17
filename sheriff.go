package sheriff

import (
	"reflect"
	"strings"

	"github.com/hashicorp/go-version"
)

type Options struct {
	Groups     []string
	ApiVersion *version.Version
}

// OnlyGroups sets groups to marshal.
//
// If a group has been set, only fields with that specific group will be marshalled.
// A field with multiple groups (comma-separated) will result in marshalling of that field
// if one of their groups is specified.
func (m *Options) SetOnlyGroups(groups []string) {
	m.Groups = groups
}

// ApiVersion sets the API version to use when marshalling.
//
// The tags `since` and `until` use the API version setting.
// Specifying the API version as "1" and having an until setting
// of "2" will result in the field being marshalled.
// Specifying a since setting of "2" with the same API version specified,
// will result in the field not being marshalled.
// Only when switching the API version to "2", the field with `since` will be
// marshalled.
func (m *Options) SetApiVersion(apiVersion string) error {
	v, err := version.NewVersion(apiVersion)
	m.ApiVersion = v
	return err
}

// NewOptions creates new options with default settings.
func NewOptions() *Options {
	return &Options{}
}

// Marshaller is the interface models have to implement in order to conform to marshalling.
type Marshaller interface {
	Marshal(options *Options) (interface{}, error)
}

// Marshal encodes the passed data into JSON.
func Marshal(options *Options, data interface{}) (interface{}, error) {
	v := reflect.ValueOf(data)
	t := v.Type()

	dest := make(map[string]interface{})
	checkGroups := len(options.Groups) > 0

	if t.Kind() == reflect.Ptr {
		// follow pointer
		t = t.Elem()
	}
	if v.Kind() == reflect.Ptr {
		// follow pointer
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		if marshaller, ok := data.(Marshaller); ok {
			return marshaller.Marshal(options)
		}
		return data, nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i)

		jsonTag, jsonOpts := parseTag(field.Tag.Get("json"))

		if jsonTag == "-" {
			continue
		}
		if jsonOpts.Contains("omitempty") && isEmptyValue(val) {
			continue
		}

		if checkGroups {
			groups := strings.Split(field.Tag.Get("groups"), ",")

			shouldShow := listContains(groups, options.Groups)
			if !shouldShow || len(groups) == 0 {
				continue
			}
		}

		if since := field.Tag.Get("since"); since != "" {
			sinceVersion, err := version.NewVersion(since)
			if err != nil {
				return nil, err
			}
			if options.ApiVersion.LessThan(sinceVersion) {
				continue
			}
		}

		if until := field.Tag.Get("until"); until != "" {
			untilVersion, err := version.NewVersion(until)
			if err != nil {
				return nil, err
			}
			if options.ApiVersion.GreaterThan(untilVersion) {
				continue
			}
		}

		actualValue := val.Interface()

		if marshaller, ok := actualValue.(Marshaller); ok {
			var err error
			actualValue, err = marshaller.Marshal(options)
			if err != nil {
				return nil, err
			}
		}

		dest[jsonTag] = actualValue
	}

	return dest, nil
}

func contains(key string, list []string) bool {
	for _, innerKey := range list {
		if key == innerKey {
			return true
		}
	}
	return false
}

func listContains(a []string, b []string) bool {
	for _, key := range a {
		if contains(key, b) {
			return true
		}
	}
	return false
}
