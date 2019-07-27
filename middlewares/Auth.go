package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shorturl/app"
	"shorturl/commons"
	"strconv"
	"strings"
	"time"
)

var authUsers = make(map[string]app.UserModel)

func Auth() gin.HandlerFunc {
	var userModel = app.UserModel{}
	users := userModel.GetUsers()
	for _, v := range users {
		authUsers[v.AppId] = v
	}
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": commons.UnAuthorized,
				"msg":  "token is empty",
				"data": "",
			})
			return
		}
		tokens := strings.Split(token, ".")
		if tokens == nil || len(tokens) != 3 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": commons.UnAuthorized,
				"msg":  "token is invalid",
				"data": "",
			})
			return
		} else {
			appId := tokens[0]
			appToken := tokens[1]
			timestamp := tokens[2]
			i, err := strconv.Atoi(timestamp)
			if err != nil || int(time.Now().Unix())-i > 100000 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code": commons.UnAuthorized,
					"msg":  "token is expired",
					"data": "",
				})
				return
			}
			if user, ok := authUsers[appId]; !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code": commons.UnAuthorized,
					"msg":  "token is wrong",
					"data": "",
				})
				return
			} else {
				expected := commons.MD5(user.AppSecret + timestamp)
				if expected != appToken {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"code": commons.UnAuthorized,
						"msg":  "token auth failed",
						"data": "",
					})
					return
				} else {
					c.Set("userId", user.Id)
				}
			}
		}
		c.Next()
	}
}
