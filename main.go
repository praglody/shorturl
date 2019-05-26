package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"shorturl/models"
	"shorturl/routers"
	"time"
)

func main() {
	if models.Conf.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	file, _ := os.Create("storage/logs/access.log")
	gin.DefaultWriter = file
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		//定制日志格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	err := logs.SetLogger(logs.AdapterFile, `{"filename":"storage/logs/app.log"}`)

	if err != nil {
		log.Fatalln("Log init failed, error: " + err.Error())
	}

	routers.Route(r)

	err = r.Run(":" + models.Conf.AppPort)
	if err != nil {
		log.Fatalln("Server start failed, error: " + err.Error())
	}
}

func init() {
	dir, _ := os.Getwd()
	file := dir + "/.env"

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Panicf("conf file [%s]  not found!", file)
	}
	models.Conf.InitConfig(file)
}
