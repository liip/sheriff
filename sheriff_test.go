package sheriff

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

type TestGroupsModel struct {
	DefaultMarshal            string   `json:"default_marshal"`
	NeverMarshal              string   `json:"-"`
	OnlyGroupTest             string   `json:"only_group_test" groups:"test"`
	OnlyGroupTestNeverMarshal string   `json:"-" groups:"test"`
	OnlyGroupTestOther        string   `json:"only_group_test_other" groups:"test-other"`
	GroupTestAndOther         string   `json:"group_test_and_other" groups:"test,test-other"`
	OmitEmpty                 string   `json:"omit_empty,omitempty"`
	OmitEmptyGroupTest        string   `json:"omit_empty_group_test,omitempty" groups:"test"`
	SliceString               []string `json:"slice_string,omitempty" groups:"test"`
}

func (data *TestGroupsModel) Marshal(options *Options) (interface{}, error) {
	return Marshal(options, data)
}

type SliceString []string

func (data SliceString) Marshal(options *Options) (interface{}, error) {
	list := make([]interface{}, len(data))
	for i, item := range data {
		target, err := Marshal(options, item)
		if err != nil {
			return nil, err
		}
		list[i] = target
	}
	return list, nil
}

func TestMarshal_GroupsValidGroup(t *testing.T) {
	testModel := &TestGroupsModel{
		DefaultMarshal:     "DefaultMarshal",
		NeverMarshal:       "NeverMarshal",
		OnlyGroupTest:      "OnlyGroupTest",
		OnlyGroupTestOther: "OnlyGroupTestOther",
		GroupTestAndOther:  "GroupTestAndOther",
		OmitEmpty:          "OmitEmpty",
		OmitEmptyGroupTest: "OmitEmptyGroupTest",
		SliceString:        []string{"test", "bla"},
	}

	o := &Options{
		Groups: []string{"test"},
	}

	actualMap, err := Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"only_group_test":       "OnlyGroupTest",
		"omit_empty_group_test": "OmitEmptyGroupTest",
		"group_test_and_other":  "GroupTestAndOther",
		"slice_string":          []string{"test", "bla"},
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

func TestMarshal_GroupsValidGroupOmitEmpty(t *testing.T) {
	testModel := &TestGroupsModel{
		DefaultMarshal:     "DefaultMarshal",
		NeverMarshal:       "NeverMarshal",
		OnlyGroupTest:      "OnlyGroupTest",
		OnlyGroupTestOther: "OnlyGroupTestOther",
		GroupTestAndOther:  "GroupTestAndOther",
		OmitEmpty:          "OmitEmpty",
		SliceString:        []string{"test", "bla"},
	}

	o := &Options{
		Groups: []string{"test"},
	}

	actualMap, err := Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"only_group_test":      "OnlyGroupTest",
		"group_test_and_other": "GroupTestAndOther",
		"slice_string":         []string{"test", "bla"},
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

func TestMarshal_GroupsInvalidGroup(t *testing.T) {
	testModel := &TestGroupsModel{
		DefaultMarshal:     "DefaultMarshal",
		NeverMarshal:       "NeverMarshal",
		OnlyGroupTest:      "OnlyGroupTest",
		OnlyGroupTestOther: "OnlyGroupTestOther",
		GroupTestAndOther:  "GroupTestAndOther",
		OmitEmpty:          "OmitEmpty",
		OmitEmptyGroupTest: "OmitEmptyGroupTest",
	}

	o := &Options{
		Groups: []string{"foo"},
	}

	actualMap, err := Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]string{})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

func TestMarshal_GroupsNoGroups(t *testing.T) {
	testModel := &TestGroupsModel{
		DefaultMarshal:     "DefaultMarshal",
		NeverMarshal:       "NeverMarshal",
		OnlyGroupTest:      "OnlyGroupTest",
		OnlyGroupTestOther: "OnlyGroupTestOther",
		GroupTestAndOther:  "GroupTestAndOther",
		OmitEmpty:          "OmitEmpty",
		OmitEmptyGroupTest: "OmitEmptyGroupTest",
	}

	o := &Options{}

	actualMap, err := Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"default_marshal":       "DefaultMarshal",
		"only_group_test":       "OnlyGroupTest",
		"only_group_test_other": "OnlyGroupTestOther",
		"group_test_and_other":  "GroupTestAndOther",
		"omit_empty":            "OmitEmpty",
		"omit_empty_group_test": "OmitEmptyGroupTest",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type TestVersionsModel struct {
	DefaultMarshal string `json:"default_marshal"`
	NeverMarshal   string `json:"-"`
	Until20        string `json:"until_20" until:"2"`
	Until21        string `json:"until_21" until:"2.1"`
	Since20        string `json:"since_20" since:"2"`
	Since21        string `json:"since_21" since:"2.1"`
}

func TestMarshal_Versions(t *testing.T) {
	testModel := &TestVersionsModel{
		DefaultMarshal: "DefaultMarshal",
		NeverMarshal:   "NeverMarshal",
		Until20:        "Until20",
		Until21:        "Until21",
		Since20:        "Since20",
		Since21:        "Since21",
	}

	v1, err := version.NewVersion("1.0.0")
	assert.NoError(t, err)

	o := &Options{
		ApiVersion: v1,
	}

	actualMap, err := Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]string{
		"default_marshal": "DefaultMarshal",
		"until_20":        "Until20",
		"until_21":        "Until21",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))

	// Api Version 2
	v2, err := version.NewVersion("2.0.0")
	assert.NoError(t, err)

	o = &Options{
		ApiVersion: v2,
	}

	actualMap, err = Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err = json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err = json.Marshal(map[string]string{
		"default_marshal": "DefaultMarshal",
		"until_20":        "Until20",
		"until_21":        "Until21",
		"since_20":        "Since20",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))

	// Api Version 2.1
	v21, err := version.NewVersion("2.1.0")
	assert.NoError(t, err)

	o = &Options{
		ApiVersion: v21,
	}

	actualMap, err = Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err = json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err = json.Marshal(map[string]string{
		"default_marshal": "DefaultMarshal",
		"until_21":        "Until21",
		"since_20":        "Since20",
		"since_21":        "Since21",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))

	// Api Version 3.0
	v3, err := version.NewVersion("3.0.0")
	assert.NoError(t, err)

	o = &Options{
		ApiVersion: v3,
	}

	actualMap, err = Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err = json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err = json.Marshal(map[string]string{
		"default_marshal": "DefaultMarshal",
		"since_20":        "Since20",
		"since_21":        "Since21",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type TestRecursiveModel struct {
	SomeData   string               `json:"some_data" groups:"test"`
	GroupsData SliceTestGroupsModel `json:"groups_data,omitempty" groups:"test"`
}

type SliceTestGroupsModel []*TestGroupsModel

func (data SliceTestGroupsModel) Marshal(options *Options) (interface{}, error) {
	list := make([]interface{}, len(data))
	for i, item := range data {
		target, err := Marshal(options, item)
		if err != nil {
			return nil, err
		}
		list[i] = target
	}
	return list, nil
}

func TestMarshal_Recursive(t *testing.T) {
	testModel := &TestGroupsModel{
		DefaultMarshal:     "DefaultMarshal",
		NeverMarshal:       "NeverMarshal",
		OnlyGroupTest:      "OnlyGroupTest",
		OnlyGroupTestOther: "OnlyGroupTestOther",
		GroupTestAndOther:  "GroupTestAndOther",
		OmitEmpty:          "OmitEmpty",
		OmitEmptyGroupTest: "OmitEmptyGroupTest",
		SliceString:        []string{"test", "bla"},
	}
	testRecursiveModel := &TestRecursiveModel{
		SomeData:   "SomeData",
		GroupsData: SliceTestGroupsModel{testModel},
	}

	o := &Options{
		Groups: []string{"test"},
	}

	actualMap, err := Marshal(o, testRecursiveModel)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"some_data": "SomeData",
		"groups_data": []map[string]interface{}{
			{
				"only_group_test":       "OnlyGroupTest",
				"omit_empty_group_test": "OmitEmptyGroupTest",
				"group_test_and_other":  "GroupTestAndOther",
				"slice_string":          []string{"test", "bla"},
			},
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

func TestMarshalInvalidTypeError_Error(t *testing.T) {
	o := &Options{}
	_, err := Marshal(o, "test")
	assert.EqualError(t, err, "marshaller: Unable to marshal type string. Struct required.")
}
