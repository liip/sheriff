package sheriff

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

type AModel struct {
	AllGroups bool `json:"something" groups:"test"`
	TestGroup bool `json:"something_else" groups:"test-other"`
}

type TestGroupsModel struct {
	DefaultMarshal            string            `json:"default_marshal"`
	NeverMarshal              string            `json:"-"`
	OnlyGroupTest             string            `json:"only_group_test" groups:"test"`
	OnlyGroupTestNeverMarshal string            `json:"-" groups:"test"`
	OnlyGroupTestOther        string            `json:"only_group_test_other" groups:"test-other"`
	GroupTestAndOther         string            `json:"group_test_and_other" groups:"test,test-other"`
	OmitEmpty                 string            `json:"omit_empty,omitempty"`
	OmitEmptyGroupTest        string            `json:"omit_empty_group_test,omitempty" groups:"test"`
	SliceString               []string          `json:"slice_string,omitempty" groups:"test"`
	MapStringStruct           map[string]AModel `json:"map_string_struct,omitempty" groups:"test,test-other"`
	IncludeEmptyTag           string            `json:"include_empty_tag"`
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
		MapStringStruct:    map[string]AModel{"firstModel": {true, true}},
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
		"map_string_struct": map[string]map[string]bool{
			"firstModel": {
				"something": true,
			},
		},
		"slice_string": []string{"test", "bla"},
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
		MapStringStruct:    map[string]AModel{"firstModel": {true, true}},
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
		"map_string_struct": map[string]map[string]bool{
			"firstModel": {
				"something": true,
			},
		},
		"slice_string": []string{"test", "bla"},
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
		MapStringStruct:    map[string]AModel{"firstModel": {true, true}},
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
		IncludeEmptyTag:    "IncludeEmptyTag",
		MapStringStruct:    map[string]AModel{"firstModel": {true, true}},
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
		"map_string_struct": map[string]map[string]bool{
			"firstModel": {
				"something":      true,
				"something_else": true,
			},
		},
		"omit_empty":            "OmitEmpty",
		"omit_empty_group_test": "OmitEmptyGroupTest",
		"include_empty_tag":     "IncludeEmptyTag",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

func TestMarshal_GroupsValidGroupIncludeEmptyTag(t *testing.T) {
	testModel := &TestGroupsModel{
		DefaultMarshal:     "DefaultMarshal",
		OnlyGroupTest:      "OnlyGroupTest",
		GroupTestAndOther:  "GroupTestAndOther",
		OmitEmpty:          "OmitEmpty",
		OmitEmptyGroupTest: "",
		SliceString:        []string{"test", "bla"},
		IncludeEmptyTag:    "IncludeEmptyTag",
	}

	o := &Options{
		IncludeEmptyTag: true,
		Groups:          []string{"test"},
	}

	actualMap, err := Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"default_marshal":      "DefaultMarshal",
		"only_group_test":      "OnlyGroupTest",
		"group_test_and_other": "GroupTestAndOther",
		"omit_empty":           "OmitEmpty",
		"slice_string":         []string{"test", "bla"},
		"include_empty_tag":    "IncludeEmptyTag",
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

type IsMarshaller struct {
	ShouldMarshal string `json:"should_marshal" groups:"test"`
}

func (i IsMarshaller) Marshal(options *Options) (interface{}, error) {
	return Marshal(options, i)
}

type TestRecursiveModel struct {
	SomeData     string             `json:"some_data" groups:"test"`
	GroupsData   []*TestGroupsModel `json:"groups_data,omitempty" groups:"test"`
	IsMarshaller IsMarshaller       `json:"is_marshaller" groups:"test"`
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
		MapStringStruct:    map[string]AModel{"firstModel": {true, true}},
	}
	testRecursiveModel := &TestRecursiveModel{
		SomeData:     "SomeData",
		GroupsData:   []*TestGroupsModel{testModel},
		IsMarshaller: IsMarshaller{"test"},
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
				"map_string_struct": map[string]map[string]bool{
					"firstModel": {
						"something": true,
					},
				},
				"slice_string": []string{"test", "bla"},
			},
		},
		"is_marshaller": map[string]interface{}{
			"should_marshal": "test",
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type TestNoJSONTagModel struct {
	SomeData    string `groups:"test"`
	AnotherData string `groups:"test"`
}

func TestMarshal_NoJSONTAG(t *testing.T) {
	testModel := &TestNoJSONTagModel{
		SomeData:    "SomeData",
		AnotherData: "AnotherData",
	}

	o := &Options{
		Groups: []string{"test"},
	}

	actualMap, err := Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"SomeData":    "SomeData",
		"AnotherData": "AnotherData",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type UserInfo struct {
	UserPrivateInfo `groups:"private"`
	UserPublicInfo  `groups:"public"`
}
type UserPrivateInfo struct {
	Age string
}
type UserPublicInfo struct {
	ID    string
	Email string `groups:"private"`
}

func TestMarshal_ParentInherit(t *testing.T) {
	publicInfo := UserPublicInfo{ID: "F94", Email: "hello@hello.com"}
	privateInfo := UserPrivateInfo{Age: "20"}
	testModel := UserInfo{
		UserPrivateInfo: privateInfo,
		UserPublicInfo:  publicInfo,
	}

	o := &Options{
		Groups: []string{"public"},
	}

	actualMap, err := Marshal(o, testModel)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"ID": "F94",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))

}

type TimeHackTest struct {
	ATime time.Time `json:"a_time" groups:"test"`
}

func TestMarshal_TimeHack(t *testing.T) {
	hackCreationTime, err := time.Parse(time.RFC3339, "2017-01-20T18:11:00Z")
	assert.NoError(t, err)

	tht := TimeHackTest{
		ATime: hackCreationTime,
	}
	o := &Options{
		Groups: []string{"test"},
	}

	actualMap, err := Marshal(o, tht)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"a_time": "2017-01-20T18:11:00Z",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type EmptyMapTest struct {
	AMap map[string]string `json:"a_map" groups:"test"`
}

func TestMarshal_EmptyMap(t *testing.T) {
	emp := EmptyMapTest{
		AMap: make(map[string]string),
	}
	o := &Options{
		Groups: []string{"test"},
	}

	actualMap, err := Marshal(o, emp)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"a_map": make(map[string]interface{}),
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

func TestMarshal_EmptyMapJson(t *testing.T) {
	emp := EmptyMapTest{
		AMap: make(map[string]string),
	}
	o := &Options{
		Groups: []string{"test"},
	}

	actualMap, err := Marshal(o, emp)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(emp)
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type PointerTest struct {
	BoolPointer          *bool    `json:"boolPointer" groups:"test"`
	BoolPointerNil       *bool    `json:"boolPointerNil" groups:"test"`
	BoolPointerNilOmit   *bool    `json:"boolPointerNilOmit,omitempty" groups:"test"`
	FloatPointer         *float64 `json:"floatPointer" groups:"test"`
	FloatPointerNil      *float64 `json:"floatPointerNil" groups:"test"`
	FloatPointerNilOmit  *float64 `json:"floatPointerNilOmit" groups:"test"`
	IntPointer           *int     `json:"intPointer" groups:"test"`
	IntPointerNil        *int     `json:"intPointerNil" groups:"test"`
	IntPointerNilOmit    *int     `json:"intPointerNilOmit" groups:"test"`
	StringPointer        *string  `json:"stringPointer" groups:"test"`
	StringPointerNil     *string  `json:"stringPointerNil" groups:"test"`
	StringPointerNilOmit *string  `json:"stringPointerNilOmit" groups:"test"`
}

func TestMarshal_Pointer(t *testing.T) {
	boolValue := true
	intValue := -20000
	stringValue := "12%&/()§?-loaMEN"
	floatValue := 0.0
	emp := PointerTest{
		BoolPointer:   &boolValue,
		IntPointer:    &intValue,
		StringPointer: &stringValue,
		FloatPointer:  &floatValue,
	}
	o := &Options{
		Groups: []string{"test"},
	}

	actualMap, err := Marshal(o, emp)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(emp)
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type TestMarshal_Embedded struct {
	Foo string `json:"foo" groups:"test"`
}

type TestMarshal_NamedEmbedded struct {
	Qux string `json:"qux" groups:"test"`
}

// TestMarshal_EmbeddedCustom is used to test an embedded struct with a custom marshaler that is not a pointer.
type TestMarshal_EmbeddedCustom struct {
	Val int
	Set bool
}

func (t TestMarshal_EmbeddedCustom) MarshalJSON() ([]byte, error) {
	if t.Set {
		return []byte(fmt.Sprintf("%d", t.Val)), nil
	}

	return nil, nil
}

// TestMarshal_EmbeddedCustomPtr is used to test an embedded struct with a custom marshaler that is a pointer.
type TestMarshal_EmbeddedCustomPtr struct {
	Val int
	Set bool
}

func (t *TestMarshal_EmbeddedCustomPtr) MarshalJSON() ([]byte, error) {
	if t.Set {
		return []byte(fmt.Sprintf("%d", t.Val)), nil
	}

	return nil, nil
}

type TestMarshal_EmbeddedParent struct {
	*TestMarshal_Embedded
	*TestMarshal_NamedEmbedded `json:"embedded"`
	*TestMarshal_EmbeddedCustom    `json:"value"`
	*TestMarshal_EmbeddedCustomPtr `json:"value_ptr"`
	Bar                        string `json:"bar" groups:"test"`
}

func TestMarshal_EmbeddedField(t *testing.T) {
	v := TestMarshal_EmbeddedParent{
		&TestMarshal_Embedded{"Hello"},
		&TestMarshal_NamedEmbedded{"Big"},
		&TestMarshal_EmbeddedCustom{10, true},
		&TestMarshal_EmbeddedCustomPtr{20, true},
		"World",
	}
	o := &Options{Groups: []string{"test"}}

	actualMap, err := Marshal(o, v)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	t.Run("should match the original json marshal", func(t *testing.T) {
		expected, err := json.Marshal(v)
		assert.NoError(t, err)

		assert.JSONEq(t, string(expected), string(actual))
	})

	t.Run("should match the expected map", func(t *testing.T) {
		expectedMap, err := json.Marshal(map[string]interface{}{
			"bar": "World",
			"foo": "Hello",
			"value":     10,
			"value_ptr": 20,
			"embedded": map[string]interface{}{
				"qux": "Big",
			},
		})
		assert.NoError(t, err)
		assert.JSONEq(t, string(expectedMap), string(actual))
	})
}

type TestMarshal_EmbeddedEmpty struct {
	Foo string
}

type TestMarshal_EmbeddedParentEmpty struct {
	*TestMarshal_EmbeddedEmpty
	Bar string `json:"bar" groups:"test"`
}

func TestMarshal_EmbeddedFieldEmpty(t *testing.T) {
	v := TestMarshal_EmbeddedParentEmpty{
		&TestMarshal_EmbeddedEmpty{"Hello"},
		"World",
	}
	o := &Options{Groups: []string{"test"}}

	actualMap, err := Marshal(o, v)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"bar": "World",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type InterfaceableBeta struct {
	Integer int    `json:"integer" groups:"safe"`
	Secret  string `json:"secret"`
}
type InterfaceableCharlie struct {
	Integer int    `json:"integer" groups:"safe"`
	Secret  string `json:"secret"`
}
type ArrayOfInterfaceable []CanHazInterface
type CanHazInterface interface {
}
type InterfacerAlpha struct {
	Plaintext     string               `json:"plaintext" groups:"safe"`
	Secret        string               `json:"secret"`
	Nested        InterfaceableBeta    `json:"nested" groups:"safe"`
	Interfaceable ArrayOfInterfaceable `json:"interfaceable" groups:"safe"`
}

func TestMarshal_ArrayOfInterfaceable(t *testing.T) {
	a := InterfacerAlpha{
		"I am plaintext",
		"I am a secret",
		InterfaceableBeta{
			100,
			"Still a secret",
		},
		ArrayOfInterfaceable{
			InterfaceableBeta{200, "Still a secret good"},
			InterfaceableCharlie{300, "Still a secret exellect"},
		}}

	o := &Options{
		Groups: []string{"safe"},
	}

	actualMap, err := Marshal(o, a)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"interfaceable": []map[string]interface{}{
			map[string]interface{}{"integer": 200},
			map[string]interface{}{"integer": 300},
		},
		"nested":    map[string]interface{}{"integer": 100},
		"plaintext": "I am plaintext",
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type TestInlineStruct struct {
	// explicitly testing unexported fields
	// golangci-lint complains about it and that's ok to ignore.
	tableName        struct{ Test string } `json:"-"`   //nolint
	tableNameWithTag struct{ Test string } `json:"foo"` //nolint

	Field  string  `json:"field"`
	Field2 *string `json:"field2"`
}

func TestMarshal_InlineStruct(t *testing.T) {
	v := TestInlineStruct{
		tableName:        struct{ Test string }{"test"},
		tableNameWithTag: struct{ Test string }{"testWithTag"},
		Field:            "World",
		Field2:           nil,
	}
	o := &Options{}

	actualMap, err := Marshal(o, v)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"field":  "World",
		"field2": nil,
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

type TestInet struct {
	IPv4 net.IP `json:"ipv4"`
	IPv6 net.IP `json:"ipv6"`
}

func TestMarshal_Inet(t *testing.T) {
	v := TestInet{
		IPv4: net.ParseIP("0.0.0.0").To4(),
		IPv6: net.ParseIP("::").To16(),
	}
	o := &Options{}

	actualMap, err := Marshal(o, v)
	assert.NoError(t, err)

	actual, err := json.Marshal(actualMap)
	assert.NoError(t, err)

	expected, err := json.Marshal(map[string]interface{}{
		"ipv4": net.ParseIP("0.0.0.0").To4(),
		"ipv6": net.ParseIP("::").To16(),
	})
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

func TestMarshal_AliaString(t *testing.T) {
	type ResourceName string
	type ResourceList map[ResourceName]string
	v := struct {
		Resources ResourceList
	}{
		Resources: map[ResourceName]string{
			"cpu": "100",
		},
	}
	_, err := Marshal(&Options{}, &v)
	assert.NoError(t, err)
}

type EmptyInterfaceStruct struct {
	Data interface{} `json:"data"`
}

func TestMarshal_EmptyInterface(t *testing.T) {
	v := EmptyInterfaceStruct{}
	o := &Options{}

	_, err := Marshal(o, v)
	assert.NoError(t, err)
}

func TestMarshal_BooleanPtrMap(t *testing.T) {
	tru := true

	toMarshal := map[string]*bool{
		"example": &tru,
		"another": nil,
	}

	marshalMap, err := Marshal(&Options{}, toMarshal)
	assert.NoError(t, err)

	marshal, err := json.Marshal(marshalMap)
	assert.NoError(t, err)

	expect, err := json.Marshal(toMarshal)
	assert.NoError(t, err)

	assert.Equal(t, string(marshal), string(expect))
}

func TestMarshal_NilSlice(t *testing.T) {
	var stringSlice []string // nil slice

	marshalSlice, err := Marshal(&Options{}, stringSlice)
	assert.NoError(t, err)

	jsonResult, err := json.Marshal(marshalSlice)
	assert.NoError(t, err)

	expect := "null"

	assert.Equal(t, expect, string(jsonResult))
}

func TestMarshal_EmptySlice(t *testing.T) {
	var stringSlice = []string{} // empty slice

	marshalSlice, err := Marshal(&Options{}, stringSlice)
	assert.NoError(t, err)

	jsonResult, err := json.Marshal(marshalSlice)
	assert.NoError(t, err)

	expect := "[]"

	assert.Equal(t, expect, string(jsonResult))
}

func TestMarshal_NilPointer(t *testing.T) {
	var a *AModel
	v, err := Marshal(&Options{}, a)
	assert.Nil(t, v)
	assert.NoError(t, err)
}

func TestMarshal_User(t *testing.T) {
	type JsonStringTag struct {
		Test  int64   `json:"test,string"`
		TestB bool    `json:"testb,string"`
		TestF float64 `json:"testf,string"`
		TestS string  `json:"tests,string"`
	}
	j := JsonStringTag{
		Test:  12,
		TestB: true,
		TestF: 12.0,
		TestS: "test",
	}

	m, err := Marshal(&Options{}, j)
	assert.NoError(t, err)

	d, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.Equal(t, `{"test":"12","testb":"true","testf":"12","tests":"test"}`, string(d))
}

func TestMarshal_CustomFieldFilter(t *testing.T) {
	type testStruct struct {
		TestValue   string `json:"test"`
		SecretValue string `json:"secret" hidden:"true"`
	}
	v := testStruct{
		TestValue:   "teststring",
		SecretValue: "asecretvalue",
	}

	o := &Options{
		FieldFilter: func(field reflect.StructField) (bool, error) {
			return field.Tag.Get("hidden") == "", nil
		},
	}
	m, err := Marshal(o, v)
	assert.NoError(t, err)

	d, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.Equal(t, `{"test":"teststring"}`, string(d))
}
