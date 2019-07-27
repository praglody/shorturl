package middlewares

import (
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"shorturl/commons"
	"time"
)

/**
 * 记录请求日志，加入traceId
 */
func Request() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			traceId   = commons.GetUUID()
			startTime = time.Now()
		)
		// before request
		logs.Info("request_start: %s %s", traceId, startTime)

		ctx.Next()
		var (
			now      = time.Now().Format("2006-01-02 15:04:05.000")
			duration = int(time.Now().Sub(startTime) / 1e6) //单位毫秒
			request  = ctx.Request.RequestURI
			host     = ctx.Request.Host
			clientIp = ctx.ClientIP()
			code     = ctx.Writer.Status()
			ua       = ctx.Request.UserAgent()
		)
		// after request
		logs.Info("request_end: %s %s %d %s %s %s %d %s", traceId, now, duration, request, host, clientIp, code, ua)
	}
}
