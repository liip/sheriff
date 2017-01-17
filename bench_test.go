package sheriff

import (
	"encoding/json"
	"testing"
)

type BenchmarkModel struct {
	AString string            `json:"a_string"`
	AInt    int               `json:"a_int"`
	ABool   bool              `json:"a_bool"`
	AArray  []string          `json:"a_array"`
	AMap    map[string]string `json:"a_map"`

	BString string            `json:"b_string"`
	BInt    int               `json:"b_int"`
	BBool   bool              `json:"b_bool"`
	BArray  []string          `json:"b_array"`
	BMap    map[string]string `json:"b_map"`

	CString string            `json:"c_string"`
	CInt    int               `json:"c_int"`
	CBool   bool              `json:"c_bool"`
	CArray  []string          `json:"c_array"`
	CMap    map[string]string `json:"c_map"`

	DString string            `json:"d_string"`
	DInt    int               `json:"d_int"`
	DBool   bool              `json:"d_bool"`
	DArray  []string          `json:"d_array"`
	DMap    map[string]string `json:"d_map"`
}

func testData() *BenchmarkModel {
	return &BenchmarkModel{
		AString: "str",
		AInt:    1123,
		ABool:   false,
		AArray:  []string{"a", "b", "c"},
		AMap:    map[string]string{"a": "b", "c": "d", "e": "f"},

		BString: "str",
		BInt:    1123,
		BBool:   false,
		BArray:  []string{"a", "b", "c"},
		BMap:    map[string]string{"a": "b", "c": "d", "e": "f"},

		CString: "str",
		CInt:    1123,
		CBool:   false,
		CArray:  []string{"a", "b", "c"},
		CMap:    map[string]string{"a": "b", "c": "d", "e": "f"},

		DString: "str",
		DInt:    1123,
		DBool:   false,
		DArray:  []string{"a", "b", "c"},
		DMap:    map[string]string{"a": "b", "c": "d", "e": "f"},
	}
}

func BenchmarkModelsMarshaller_Marshal_NativeJSON(b *testing.B) {
	s := testData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(s)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkModelsMarshaller_Marshal(b *testing.B) {
	s := testData()
	o := NewOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := Marshal(o, s)
		if err != nil {
			b.Fatal(err)
		}
		_, err = json.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
