package opt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestYAML(t *testing.T) {
	{
		type YAMLStruct struct {
			Val Option[int] `yaml:"val"`
		}

		some := Some(123)
		yamlStruct := &YAMLStruct{Val: some}

		marshal, err := yaml.Marshal(yamlStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, "val: 123\n", string(marshal))

		var unmarshalYAMLStruct YAMLStruct
		err = yaml.Unmarshal(marshal, &unmarshalYAMLStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, yamlStruct, &unmarshalYAMLStruct)
	}

	{
		type YAMLStruct struct {
			Val Option[int] `yaml:"val"`
		}

		none := None[int]()
		yamlStruct := &YAMLStruct{Val: none}

		marshal, err := yaml.Marshal(yamlStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, "val: null\n", string(marshal))

		var unmarshalYAMLStruct YAMLStruct
		err = yaml.Unmarshal(marshal, &unmarshalYAMLStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, yamlStruct, &unmarshalYAMLStruct)
	}
}
