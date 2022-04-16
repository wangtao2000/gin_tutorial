/**
    @author: wangtao
    @date: 2022/4/17
    @note: 将数据和model绑定并验证
**/

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

// binding:"required" 是gin中tag，用于表示这个字段是否为必要的
// 一些验证的方法可以参见 https://pkg.go.dev/gopkg.in/go-playground/validator.v9

type loginReq struct {
	Username string `json:"username" xml:"username" form:"username" binding:"required"`
	Password string `json:"password" xml:"password" form:"password" binding:"required"`
}

func login(loginReqData *loginReq, c *gin.Context) {
	if loginReqData.Username == "aaa" && loginReqData.Password == "bbb" {
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in : " + loginReqData.Username})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "authorized error  : " + loginReqData.Username})
}

func main() {
	router := gin.Default()

	router.POST("/loginJSON", func(c *gin.Context) {
		var loginReqData loginReq
		// ShouldBind 和 MustBind的区别：ShouldBind错误后会返回error，可以进行进一步的处理
		// MustBind 发生错误后请求会使用c.AbortWithError(400, err).SetType(ErrorTypeBind)立即中断
		err := c.ShouldBindJSON(&loginReqData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		login(&loginReqData, c)

	})

	router.POST("/login", func(c *gin.Context) {
		var loginReqData loginReq
		// Bind方法会根据Content-Type来进行匹配
		err := c.Bind(&loginReqData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		login(&loginReqData, c)
	})

	router.POST("/loginMust", func(c *gin.Context) {
		var loginReqData loginReq
		// 必须是JSON类型
		err := c.MustBindWith(&loginReqData, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		login(&loginReqData, c)
	})

	_ = router.Run()
}
