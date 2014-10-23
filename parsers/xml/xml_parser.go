package xml

import (
	"encoding/xml"
)

type XmlParser struct {
}

func (XmlParser) ParseContent(data *[]byte) (interface{}, error) {

	var object interface{}
	data_ := *data
	err := xml.Unmarshal(data_, &object)
	if err != nil {
		return nil, err
	}

	return object, nil
}
