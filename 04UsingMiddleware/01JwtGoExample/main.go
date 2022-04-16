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
	"errors"
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

func resultWithSuccessMessage(message string, c *gin.Context) {
	result(SUCCESS, message, map[string]interface{}{}, c)
}

func resultWithSuccessMessageData(message string, data interface{}, c *gin.Context) {
	result(SUCCESS, message, data, c)
}

var (
	signKey *rsa.PrivateKey
)

var (
	prvkey []byte
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

	actionsGroup := router.Group("")
	actionsGroup.Use(JwtAuth)
	{
		actionsGroup.GET("/hello", hello)
	}

	_ = router.Run()
}

func hello(c *gin.Context) {
	claims, ok := c.Get("claims")
	if !ok {
		resultWithErrorMessage("认证失败", c)
		return
	}
	infoClaims, ok := claims.(*userInfoClaims)
	if !ok {
		resultWithErrorMessage("token有误", c)
		return
	}
	resultWithSuccessMessage("你经过认证了"+infoClaims.userInfo.Name, c)
}

// JwtAuth
// Jwt认证中间件
func JwtAuth(c *gin.Context) {
	token := c.Request.Header.Get("x-token")
	if token == "" {
		resultWithErrorMessage("未登录或非法访问", c)
		c.Abort()
		return
	}

	claims, err := validToken(token)
	if err != nil {
		resultWithErrorMessage("token有误!", c)
		c.Abort()
		return
	}

	tokenStorageInterface, ok := tokenCache.Load(claims.userInfo.Name)
	if !ok {
		resultWithErrorMessage("token不存在此工作负载上!", c)
		c.Abort()
		return
	}
	tokenStorage := tokenStorageInterface.(string)
	if tokenStorage != token {
		resultWithErrorMessage("token不一致!", c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
	c.Next()
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if len(username) == 0 || len(password) == 0 {
		resultWithErrorMessage("账号或者密码不合法", c)
		return
	}
	if _, ok := tokenCache.Load(username); ok {
		tokenCache.Delete(username)
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

func validToken(tokenString string) (*userInfoClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &userInfoClaims{}, func(token *jwt.Token) (interface{}, error) {
		return prvkey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*userInfoClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("claims Error")
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

}
