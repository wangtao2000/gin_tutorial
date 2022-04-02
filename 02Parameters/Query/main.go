/**
    @author: wangtao
    @date: 2022/4/3
    @note:
**/

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/welcome", func(context *gin.Context) {

		// 如果不传firstname则默认为Guest
		firstName := context.DefaultQuery("firstname", "Guest")
		// 无默认值
		lastName := context.Query("lastname")
		context.String(http.StatusOK, "Hello , %s %s ! ", firstName, lastName)
	})

	// 请求 /echo?list=1&list=2&list=3  通过queryArray拿到一个切片
	// 响应 [1 2 3] : len : 3
	router.GET("/echo", func(context *gin.Context) {
		list := context.QueryArray("list")
		context.String(http.StatusOK, "%v : len : %d", list, len(list))
	})

	// 请求 /query_map?map%5Ba%5D=aaa&map%5Bb%5D=bbb -->%5B [  %5D ]
	// 响应 {"a":"aaa","b":"bbb"}
	router.GET("/query_map", func(context *gin.Context) {
		m := context.QueryMap("map")
		context.JSON(http.StatusOK, m)
	})

	err := router.Run()
	if err != nil {
		return
	}
}
