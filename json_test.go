package opt

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionSerdeJSONForSomeValue(t *testing.T) {
	{
		type JSONStruct struct {
			Val Option[int] `json:"val"`
		}

		some := Some[int](123)
		jsonStruct := &JSONStruct{Val: some}

		marshal, err := json.Marshal(jsonStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, string(marshal), `{"val":123}`)

		var unmarshalJSONStruct JSONStruct
		err = json.Unmarshal(marshal, &unmarshalJSONStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
	}

	{
		type JSONStruct struct {
			Val Option[string] `json:"val"`
		}

		some := Some[string]("foobar")
		jsonStruct := &JSONStruct{Val: some}

		marshal, err := json.Marshal(jsonStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, string(marshal), `{"val":"foobar"}`)

		var unmarshalJSONStruct JSONStruct
		err = json.Unmarshal(marshal, &unmarshalJSONStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
	}

	{
		type JSONStruct struct {
			Val Option[bool] `json:"val"`
		}

		some := Some[bool](false)
		jsonStruct := &JSONStruct{Val: some}

		marshal, err := json.Marshal(jsonStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, string(marshal), `{"val":false}`)

		var unmarshalJSONStruct JSONStruct
		err = json.Unmarshal(marshal, &unmarshalJSONStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
	}

	{
		type Inner struct {
			B *bool `json:"b,omitempty"`
		}
		type JSONStruct struct {
			Val Option[Inner] `json:"val"`
		}

		{
			falsy := false
			some := Some[Inner](Inner{
				B: &falsy,
			})
			jsonStruct := &JSONStruct{Val: some}

			marshal, err := json.Marshal(jsonStruct)
			assert.NoError(t, err)
			assert.EqualValues(t, string(marshal), `{"val":{"b":false}}`)

			var unmarshalJSONStruct JSONStruct
			err = json.Unmarshal(marshal, &unmarshalJSONStruct)
			assert.NoError(t, err)
			assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
		}

		{
			some := Some[Inner](Inner{
				B: nil,
			})
			jsonStruct := &JSONStruct{Val: some}

			marshal, err := json.Marshal(jsonStruct)
			assert.NoError(t, err)
			assert.EqualValues(t, string(marshal), `{"val":{}}`)

			var unmarshalJSONStruct JSONStruct
			err = json.Unmarshal(marshal, &unmarshalJSONStruct)
			assert.NoError(t, err)
			assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
		}
	}
}

func TestOptionSerdeJSONForNoneValue(t *testing.T) {
	type JSONStruct struct {
		Val Option[int] `json:"val"`
	}
	some := None[int]()
	jsonStruct := &JSONStruct{Val: some}

	marshal, err := json.Marshal(jsonStruct)
	assert.NoError(t, err)
	assert.EqualValues(t, string(marshal), `{"val":null}`)

	var unmarshalJSONStruct JSONStruct
	err = json.Unmarshal(marshal, &unmarshalJSONStruct)
	assert.NoError(t, err)
	assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
}

func TestOption_UnmarshalJSON_withEmptyJSONString(t *testing.T) {
	type JSONStruct struct {
		Val Option[int] `json:"val"`
	}

	var unmarshalJSONStruct JSONStruct
	err := json.Unmarshal([]byte("{}"), &unmarshalJSONStruct)
	assert.NoError(t, err)
	assert.EqualValues(t, &JSONStruct{
		Val: None[int](),
	}, &unmarshalJSONStruct)
}

func TestOption_MarshalJSON_shouldReturnErrorWhenInvalidJSONStructInputHasCome(t *testing.T) {
	type JSONStruct struct {
		Val Option[chan interface{}] `json:"val"` // chan type is unsupported on json marshaling
	}

	ch := make(chan interface{})
	some := Some[chan interface{}](ch)
	jsonStruct := &JSONStruct{Val: some}
	_, err := json.Marshal(jsonStruct)
	assert.Error(t, err)
}

func TestOption_UnmarshalJSON_shouldReturnErrorWhenInvalidJSONStringInputHasCome(t *testing.T) {
	type JSONStruct struct {
		Val Option[int] `json:"val"`
	}

	var unmarshalJSONStruct JSONStruct
	err := json.Unmarshal([]byte(`{"val":"__STRING__"}`), &unmarshalJSONStruct)
	assert.Error(t, err)
}
