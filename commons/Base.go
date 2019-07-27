package commons

import (
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	. "shorturl/configs"
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
var Redis *redis.Client

func init() {
	initDB()
	initRedis()
}

// 初始化mysql
func initDB() {
	var db *gorm.DB
	var (
		username = Conf.String("DB_USERNAME")
		pass     = Conf.String("DB_PASSWORD")
		host     = Conf.String("DB_HOST")
		port     = Conf.String("DB_PORT")
		database = Conf.String("DB_DATABASE")
		charset  = Conf.String("DB_CHARSET")
		dialect  = Conf.String("DB_Dialect")
	)
	dsn := username + ":" + pass + "@tcp(" + host + ":" + port + ")/" + database + "?charset=" + charset
	db, err := gorm.Open(dialect, dsn)
	if err != nil {
		log.Fatalf("init DB connect failed, error: %s", err.Error())
	}
	err = db.DB().Ping()
	if Conf.String("APP_ENV") != "prod" {
		db.LogMode(true)
	}
	if err != nil {
		log.Fatalf("init DB connect failed, error: %s", err.Error())
	} else {
		DB = db
		log.Println("init DB connect success")
	}
}

// 初始化redis
func initRedis() {
	var (
		host        = Conf.String("REDIS_HOST")
		port        = Conf.String("REDIS_PORT")
		pass        = Conf.String("REDIS_PASS")
		idleConn, _ = Conf.Int("REDIS_MIN_IDLE")
	)
	Redis = redis.NewClient(&redis.Options{
		Addr:         host + ":" + port,
		Password:     pass,
		MaxRetries:   3,
		MinIdleConns: idleConn,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		log.Println("init Redis Pool failed")
	} else {
		log.Println("init Redis Pool success")
	}
}
