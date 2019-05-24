package main

import (
	"net/http"
	"net/url"
	"strconv"
)

func main() {
	Add()
}

func Add() {
	for i := 0; i < 100000; i++ {
		client := http.Client{}
		response, _ := client.PostForm("http://127.0.0.1:8080/v1/create", url.Values{
			"url": []string{"http://www.baidu.com/" + strconv.Itoa(i)},
		})
		err := response.Body.Close()
		if err != nil {
			panic(err)
		}
	}
}
