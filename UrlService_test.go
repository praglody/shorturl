package main

import (
	"github.com/stretchr/testify/assert"
	"shorturl/app/services"
	"testing"
)

func TestGenCode(t *testing.T) {
	code, err := services.UrlService{}.GenCode("http://www.baidu.com")

	assert.Nil(t, err)

	assert.NotEmpty(t, code)
}
