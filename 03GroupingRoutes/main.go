/**
    @author: wangtao
    @date: 2022/4/7
    @note:
**/

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func hello(context *gin.Context) {
	context.String(http.StatusOK, "hello %s", context.FullPath())
}

func main() {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/hello", hello)
	}
	v2 := router.Group("v2")
	{
		v2.GET("/hello", hello)
	}

	_ = router.Run()
}
