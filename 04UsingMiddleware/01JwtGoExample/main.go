/**
    @author: wangtao
    @date: 2022/4/13
    @note:
**/

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"sync"
	"time"
)

const (
	ERROR   = 7
	SUCCESS = 0
)

type response struct {
	Code    int         `bson:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func result(code int, message string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, response{Data: data, Message: message, Code: code})
}

func resultWithErrorMessage(message string, c *gin.Context) {
	result(ERROR, message, map[string]interface{}{}, c)
}

func resultWithSuccessMessageData(message string, data interface{}, c *gin.Context) {
	result(SUCCESS, message, data, c)
}

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
)

var (
	prvkey []byte
	pubkey []byte
)

type userInfo struct {
	Name string
}

type userInfoClaims struct {
	*jwt.StandardClaims
	TokenType string
	userInfo
}

// 储存token的cache
var tokenCache sync.Map

func main() {
	initKey()

	router := gin.Default()

	loginGroup := router.Group("")
	{
		loginGroup.POST("/login", login)
	}

	_ = router.Run()
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if len(username) == 0 || len(password) == 0 {
		resultWithErrorMessage("账号或者密码不合法", c)
		return
	}
	if _, ok := tokenCache.Load(username); ok {
		resultWithErrorMessage("您已经在别处登录", c)
		return
	}

	token, err := createToken(username)
	if err != nil {
		resultWithErrorMessage("生成token失败: "+err.Error(), c)
		return
	}
	tokenCache.Store(username, token)
	resultWithSuccessMessageData("登录成功", struct {
		User  string `json:"user"`
		Token string `json:"token"`
	}{username, token}, c)

}

//todo: token验证

func createToken(user string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, userInfoClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
		"level1",
		userInfo{
			user,
		},
	})
	return t.SignedString(prvkey)
}

func initKey() {
	var err error
	signKey, err = rsa.GenerateKey(rand.Reader, 256)
	if err != nil {
		panic("can't generate key ")
	}
	derStream := x509.MarshalPKCS1PrivateKey(signKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	verifyKey = &signKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(verifyKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)

}
