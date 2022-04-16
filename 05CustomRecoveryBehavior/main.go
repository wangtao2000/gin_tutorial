/**
    @author: wangtao
    @date: 2022/4/16
    @note: 自定义恢复行为
**/

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	router.GET("/panic", func(c *gin.Context) {
		panic("a panic example")
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello")
	})

	_ = router.Run()
}
