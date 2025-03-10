package opt

import "gopkg.in/yaml.v3"

func (o Option[T]) MarshalYAML() (any, error) {
	if o.IsNone() {
		return nil, nil
	}

	return o.Unwrap(), nil
}

func (o *Option[T]) UnmarshalYAML(value *yaml.Node) error {
	var v T
	err := value.Decode(&v)
	if err != nil {
		return err
	}

	*o = Some(v)
	return nil
}
