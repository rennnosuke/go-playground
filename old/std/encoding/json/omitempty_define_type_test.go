package json

import (
	"encoding/json"
	"testing"
)

type Field struct {
	Name string `json:"name"`
}

type FieldNotIgnored struct {
	Field Field `json:"field,omitempty"`
}

type FieldIgnorable struct {
	Field *Field `json:"field,omitempty"`
}

// TestEncodingJson_OmitEmptyDefinedType
// encoding/jsonによるJSON Encodingで、定義型fieldにomitemptyを付与してもfieldが無視されないことを検証します
func TestEncodingJson_OmitEmptyDefinedType(t *testing.T) {
	// 1. 定義型fieldをomitemptyできないことを検証
	p := FieldNotIgnored{}
	b, err := json.Marshal(&p)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("`field` is not ignored: %s\n", b)

	// 2. 定義型参照fieldであればomitemptyできることを検証
	pfi := FieldIgnorable{}
	b, err = json.Marshal(&pfi)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("`field` can be ignored: %s\n", b)
}
