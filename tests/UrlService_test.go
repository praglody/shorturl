package tests

import (
	"github.com/stretchr/testify/assert"
	"shorturl/services"
	"testing"
)

func TestGenCode(t *testing.T) {
	code, err := services.UrlService{}.GenCode("http://www.baidu.com")
	assert.Nil(t, err)
	assert.NotEmpty(t, code)
}

func TestTransToCode(t *testing.T) {
	code := services.TransToCode(10000)
	assert.Equal(t, "", code)
}
