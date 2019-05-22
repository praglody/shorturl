package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"shorturl/config"
)

var DB *sql.DB

type MapUrl struct {
	id          int
	shorturl    string
	url         string
	md5         string
	click       int
	create_time int
	ShortUrl    string `json:"shorturl"`
	Url         string `json:"url"`
}

//func init() {
//	var err error
//	DB, err = sql.Open("mysql", config.MysqlDSN)
//	if err != nil {
//		panic(err)
//	}
//}

func GetRow(shorturl string) (*MapUrl, error) {
	var url MapUrl
	var err error

	DB, err = sql.Open("mysql", config.MysqlDSN)
	if err != nil {
		panic(err)
	} else {
		defer DB.Close()
	}
	err = DB.QueryRow("select * from tbl_mapurl where shorturl = ?", shorturl).
		Scan(&url.id, &url.shorturl, &url.url, &url.md5, &url.click, &url.create_time)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, err
		default:
			return nil, err
		}
	}

	url.ShortUrl = url.shorturl
	url.Url = url.url
	return &url, nil
}
