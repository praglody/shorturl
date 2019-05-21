package app

import (
	"fmt"
	"net/http"
)

func Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello, short url")
	})

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	panic(err)
}
