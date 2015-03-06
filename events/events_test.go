package events

import (
	"github.com/diggs/connectrix/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatedEventContent(t *testing.T) {

	// content should be equal to the string representation of the raw content when no template is specified
	object := ""
	eventData := []byte(TestData)
	eventType := &config.EventType{
		Template: "",
		Type:     "test",
	}

	content, err := makeTemplatedEventContent(object, eventType, &eventData)

	assert.Nil(t, err)
	assert.Equal(t, TestData, content)
}

var TestData = `{
  "ref": "refs/heads/gh-pages",
  "after": "4d2ab4e76d0d405d17d1a0f2b8a6071394e3ab40",
  "before": "993b46bdfc03ae59434816829162829e67c4d490",
  "created": false}`
