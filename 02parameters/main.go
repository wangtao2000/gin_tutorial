/**
    @author: wangtao
    @date: 2022/4/1
    @note:
**/

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	// 会匹配/user/name 不会单独匹配/user 或者 /name
	router.GET("/user/:name", func(context *gin.Context) {
		name := context.Param("name")
		context.String(http.StatusOK, "hello,%s!", name)
	})
	// GET /user/wangtao HTTP/1.1
	// Host: 127.0.0.1:8080
	// Response: hello,wangtao!
	//
	// GET /user HTTP/1.1
	// Host: 127.0.0.1:8080
	// Response: 404 page not found

	// 会匹配/user/name/ 以及/user/name/send
	// 解析分别拿到/和/send
	router.GET("/user/:name/*action", func(context *gin.Context) {
		name := context.Param("name")
		action := context.Param("action")
		message := name + " is " + action
		context.String(http.StatusOK, message)
	})

	router.POST("/user/:name/*action", func(context *gin.Context) {
		b := context.FullPath() == "/user/:name/*action" // true
		context.String(http.StatusOK, "%t", b)
	})

	// 确切的路由会比参数路由先匹配
	router.GET("/user/group", func(context *gin.Context) {
		context.String(http.StatusOK, "The available groups are [...]")
	})

	err := router.Run()
	if err != nil {
		return
	}
}
