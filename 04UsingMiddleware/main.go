/**
    @author: wangtao
    @date: 2022/4/12
    @note:
**/

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()

	notAuthorized := router.Group("")
	{
		notAuthorized.GET("/hey", func(context *gin.Context) {
			context.String(http.StatusOK, "这是一个无需认证即可使用的接口")
		})
	}

	/* 需要认证的路由 */
	authorized := router.Group("")
	authorized.Use(func(context *gin.Context) {
		token := context.Request.Header.Get("token")

		/* 一些鉴权的代码 */

		if token == "right" {
			context.Set("authorizedTime", time.Now())
			context.Next()
		} else {
			context.String(http.StatusForbidden, "没有权限访问！")
			context.Abort()
			return
		}
	})
	{
		authorized.GET("/hello", func(context *gin.Context) {
			val, exists := context.Get("authorizedTime")
			if !exists {
				context.String(http.StatusInternalServerError, "认证时间获取失败")
			}
			context.String(http.StatusOK, "当你看到这段话,代表你有权限,认证时间:%s", val)
		})
	}

	_ = router.Run()
}
