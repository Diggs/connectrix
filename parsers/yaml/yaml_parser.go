package yaml

import (
	"gopkg.in/yaml.v1"
)

type YamlParser struct {
}

func (YamlParser) ParseContent(data *[]byte) (interface{}, error) {

	var object interface{}
	data_ := *data
	err := yaml.Unmarshal(data_, &object)
	if err != nil {
		return nil, err
	}

	return object, nil
}
