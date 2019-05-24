package models

import (
	"github.com/astaxie/beego/config"
)

var Conf *App

func init() {
	Conf = &App{}
}

type App struct {
	AppPort string
	AppEnv  string
	AppUrl  string

	DBDialect      string
	DBHost         string
	DBPort         string
	DBDatabase     string
	DBUser         string
	DBPass         string
	DBCharset      string
	DBMaxIdleConns int
	DBMaxOpenConns int
	DBTimeZone     string

	RedisHost      string
	RedisPort      string
	RedisPass      string
	RedisMaxIdle   int
	RedisMaxActive int
}

func (App) InitConfig(file string) {
	conf, err := config.NewConfig("ini", file)
	if err != nil {
		panic(err)
	}

	Conf.AppPort = conf.DefaultString("APP_PORT", "8080")
	Conf.AppEnv = conf.DefaultString("APP_ENV", "dev")
	Conf.AppUrl = conf.DefaultString("APP_URL", "http://127.0.0.1:8080/")

	Conf.DBDialect = conf.DefaultString("DB_Dialect", "mysql")
	Conf.DBHost = conf.DefaultString("DB_HOST", "127.0.0.1")
	Conf.DBPort = conf.DefaultString("DB_PORT", "3306")
	Conf.DBDatabase = conf.DefaultString("DB_DATABASE", "short")
	Conf.DBUser = conf.DefaultString("DB_USERNAME", "root")
	Conf.DBPass = conf.DefaultString("DB_PASSWORD", "123456")
	Conf.DBCharset = conf.DefaultString("DB_CHARSET", "utf8mb4")
	Conf.DBMaxIdleConns = conf.DefaultInt("DB_MAX_IDLE_CONNS", 5)
	Conf.DBMaxOpenConns = conf.DefaultInt("DB_MAX_OPEN_CONNS", 10)
	Conf.DBTimeZone = conf.DefaultString("DB_TIMEZONE", "+08:00")

	Conf.RedisHost = conf.DefaultString("REDIS_HOST", "127.0.0.1")
	Conf.RedisPort = conf.DefaultString("REDIS_PORT", "6379")
	Conf.RedisPass = conf.DefaultString("REDIS_PASS", "")
	Conf.RedisMaxIdle = conf.DefaultInt("REDIS_MAX_IDLE", 3)
	Conf.RedisMaxActive = conf.DefaultInt("REDIS_MAX_ACTIVE", 5)

	BaseInit()
}
