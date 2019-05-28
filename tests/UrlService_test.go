package tests

import (
	"github.com/stretchr/testify/assert"
	"shorturl/services"
	"testing"
)

func TestGenCode(t *testing.T) {
	expected, err := services.UrlService{}.GenShortUrl("http://www.baidu.com")
	assert.Nil(t, err)
	assert.NotEmpty(t, expected)

	expected, err = services.UrlService{}.GenShortUrl("http://www.iwork.com")
	assert.Nil(t, err)
	assert.NotEmpty(t, expected)
}

func TestTransToCode(t *testing.T) {
	expected := services.TransToCode(10000)
	assert.Equal(t, "bCv", expected)

	expected = services.TransToCode(20000)
	assert.Equal(t, "DqN", expected)
}
