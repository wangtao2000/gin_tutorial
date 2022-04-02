/**
    @author: wangtao
    @date: 2022/3/30
    @note: gin quick start
**/

package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		return
	}
}
