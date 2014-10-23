package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomHeaders(t *testing.T) {
	args := map[string]string{"Headers": "Authorization:Basic ZGliZ3M6WnVzdDNyNDQ= , Content-Type:application/json,Foo:bar"}
	headers := getCustomHeaders(args)
	assert.Len(t, headers, 3)

	auth := headers["Authorization"]
	assert.Equal(t, "Basic ZGliZ3M6WnVzdDNyNDQ=", auth)

	content := headers["Content-Type"]
	assert.Equal(t, "application/json", content)

	foo := headers["Foo"]
	assert.Equal(t, "bar", foo)
}

func TestSubChannelArgValidation(t *testing.T) {

	httpChannel := HttpChannel{}
	err := httpChannel.ValidateSubChannelArgs(map[string]string{"URL": "http://foo.com", "Self Signed Cert": "false"})
	assert.Nil(t, err)

	err = httpChannel.ValidateSubChannelArgs(map[string]string{"URL": "hfvsadfasdd", "Self Signed Cert": "false"})
	assert.NotNil(t, err)

	err = httpChannel.ValidateSubChannelArgs(map[string]string{"URL": "http://foo.com", "Self Signed Cert": "imnotabool"})
	assert.NotNil(t, err)
}
