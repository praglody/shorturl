package main

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"shorturl/models"
	"shorturl/routers"
)

func main() {
	initConfig()

	if models.AppConfig.Env != "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	err := logs.SetLogger(logs.AdapterFile, `{"filename":"app.log"}`)

	if err != nil {
		log.Fatalln("Log init failed, error: " + err.Error())
	}

	routers.Route(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatalln("Server start failed, error: " + err.Error())
	}
}

func initConfig() {
	conf, err := config.NewConfig("ini", "configs/app.ini")
	if err != nil {
		panic(err)
	}
	env := conf.String("env")
	port := conf.String("port")
	baseUrl := conf.String("base_url")
	dir, _ := os.Getwd()
	models.AppConfig = &models.AppConf{Env: env, Port: port, BaseUrl: baseUrl, WorkDir: dir}
}
