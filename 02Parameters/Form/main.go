/**
    @author: wangtao
    @date: 2022/4/4
    @note: 有关Post Form相关的操作
**/

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.POST("/post_form", func(context *gin.Context) {
		message := context.PostForm("message")
		nick := context.DefaultPostForm("nick", "anonymous")

		context.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})

	})

	err := router.Run()
	if err != nil {
		return
	}
}
