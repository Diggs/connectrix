package json

import (
	"encoding/json"
)

type JsonParser struct {
}

func (JsonParser) ParseContent(data *[]byte) (interface{}, error) {

	var object interface{}
	data_ := *data
	err := json.Unmarshal(data_, &object)
	if err != nil {
		return nil, err
	}

	return object, nil
}
