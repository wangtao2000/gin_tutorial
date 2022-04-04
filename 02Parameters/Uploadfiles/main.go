/**
    @author: wangtao
    @date: 2022/4/4
    @note: 文件上传相关接口学习
**/

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"path"
)

// generateName 随机生成文件名
func generateName() (string, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

func main() {
	router := gin.Default()

	// 为router设置一个较低的内存限制用于文件上传(默认是32MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(context *gin.Context) {
		file, err := context.FormFile("file")
		if err != nil {
			context.String(500, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		log.Println(file.Filename)

		// 上传文件到指定目录

		filename, err := generateName()
		if err != nil {
			context.String(500, fmt.Sprintf("generate name err: %s", err.Error()))
			return
		}
		// 获取文件后缀类型
		ext := path.Ext(file.Filename)

		// 保存文件
		err = context.SaveUploadedFile(file, "./Uploads/"+filename+ext)
		if err != nil {
			context.String(500, fmt.Sprintf("Uploads file err: %s", err.Error()))
			return
		}

		context.String(http.StatusOK, "Uploads success")
	})

	err := router.Run()
	if err != nil {
		return
	}
}
