package app

import (
	"fmt"
	"net/http"
	"shorturl/lib/dao"
)

func Run() {
	res := dao.GetRow("qwerty")
	fmt.Println("查询结果", res)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello, short url")
	})

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	panic(err)
}
