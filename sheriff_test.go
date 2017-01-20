package sheriff

import (
	"encoding/json"
	"testing"

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
