package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

type BaseModel struct {
	ID        uint  `gorm:"primary_key"`
	CreatedAt uint  `json:"created_at"`
	UpdatedAt uint  `json:"updated_at"`
	Status    uint8 `json:"status"`
}

// DB连接
var DB *gorm.DB

// redis连接
var Redis *redis.Pool

func BaseInit() {
	initDB()
	initRedis()
}

// 初始化mysql
func initDB() {
	var db *gorm.DB
	url := Conf.DBUser + ":" + Conf.DBPass + "@tcp(" + Conf.DBHost + ":" + Conf.DBPort + ")/" + Conf.DBDatabase +
		"?charset=" + Conf.DBCharset
	db, err := gorm.Open(Conf.DBDialect, url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	db.DB().SetMaxIdleConns(Conf.DBMaxIdleConns)
	db.DB().SetMaxOpenConns(Conf.DBMaxOpenConns)
	DB = db
}

// 初始化redis
func initRedis() {
	server := Conf.RedisHost + ":" + Conf.RedisPort
	redisPool := &redis.Pool{
		MaxIdle:     Conf.RedisMaxIdle,
		MaxActive:   Conf.RedisMaxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialPassword(Conf.RedisPass))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	Redis = redisPool
}
