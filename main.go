package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"shorturl/app"
	. "shorturl/configs"
	"shorturl/routers"
	"syscall"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	err := logs.SetLogger(logs.AdapterFile, `{"filename":"storage/logs/app.log"}`)
	if err != nil {
		log.Fatalln("Log init failed, error: " + err.Error())
	}
	routers.Route(engine)
	// 启动服务器，grace restart
	server := endless.NewServer(":"+Conf.String("APP_PORT"), engine)
	// 注册程序终止信号
	var signals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	}
	for _, signal := range signals {
		_ = server.RegisterSignalHook(endless.PRE_SIGNAL, signal, func() {
			app.StopTheWorld()
		})
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server start failed, error: " + err.Error())
	}
	log.Println("Server start success")
}
