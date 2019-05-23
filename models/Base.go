package models

import (
	"fmt"
	"github.com/astaxie/beego/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"os"
	"strconv"
	"time"
)

type BaseModel struct {
	ID        uint  `gorm:"primary_key"`
	CreatedAt uint  `json:"created_at"`
	UpdatedAt uint  `json:"updated_at"`
	Status    uint8 `json:"status"`
}

// DB 数据库连接
var DB *gorm.DB

// redis连接
var Redis *redis.Pool

func init() {
	initDB()
	initRedis()
}

// 初始化mysql
func initDB() {
	conf, err := config.NewConfig("ini", "configs/mysql.ini")
	if err != nil {
		panic(err)
	}
	dialect := getConfValue(conf, "dialect")
	name := getConfValue(conf, "name")
	host := getConfValue(conf, "host")
	port := getConfValue(conf, "port")
	user := getConfValue(conf, "user")
	pass := getConfValue(conf, "pass")
	charset := getConfValue(conf, "charset")

	maxIdleConns, err := strconv.Atoi(getConfValue(conf, "maxIdleConns"))
	if err != nil {
		maxIdleConns = 5
	}
	maxOpenConns, err := strconv.Atoi(getConfValue(conf, "maxOpenConns"))
	if err != nil {
		maxOpenConns = 20
	}
	url := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + name + "?charset" + charset

	db, err := gorm.Open(dialect, url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	DB = db
}

// 初始化redis
func initRedis() {
	conf, err := config.NewConfig("ini", "configs/redis.ini")
	if err != nil {
		panic(err)
	}
	host := getConfValue(conf, "host")
	port := getConfValue(conf, "port")
	pass := getConfValue(conf, "pass")
	maxIdle, err := strconv.Atoi(getConfValue(conf, "maxIdle"))
	if err != nil {
		maxIdle = 3
	}

	maxActive, err := strconv.Atoi(getConfValue(conf, "maxActive"))
	if err != nil {
		maxActive = 10
	}
	server := host + ":" + port

	var redisPool *redis.Pool

	redisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialPassword(pass))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	if err != nil {
		panic(err)
	}

	Redis = redisPool
}

func getConfValue(conf config.Configer, name string) string {
	return conf.String("product::" + name)
}
