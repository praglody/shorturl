package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type MapUrl struct {
	id          int
	shorturl    string
	url         string
	md5         string
	click       int
	create_time int
}

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:adgjmptw123@tcp(123.207.144.90)/shorturl")
	if err != nil {
		panic(err)
	}

	fmt.Println(DB)
}

func GetRow(shorturl string) *MapUrl {

	fmt.Println("GetRow")
	fmt.Println(DB)
	var url MapUrl

	err := DB.QueryRow("select * from tbl_mapurl where shorturl = ?", shorturl).
		Scan(&url.id, &url.shorturl, &url.url, &url.md5, &url.click, &url.create_time)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			fmt.Println(err)
			return nil
		default:
			fmt.Println(err)
			return nil
		}
	}

	return &url
}
