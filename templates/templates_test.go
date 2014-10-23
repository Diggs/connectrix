package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJsonEvent(t *testing.T) {

	type test struct {
		Foo string
	}

	data, err := Template(test{Foo: "hello"}, "{{.Foo}}")
	assert.Nil(t, err)
	assert.Equal(t, "hello", data)
}
