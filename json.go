package opt

import (
	"bytes"
	"encoding/json"
)

var jsonNull = []byte("null")

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return jsonNull, nil
	}

	marshal, err := json.Marshal(o.Unwrap())
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if len(data) <= 0 || bytes.Equal(data, jsonNull) {
		*o = None[T]()
		return nil
	}

	var v T
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*o = Some(v)

	return nil
}
