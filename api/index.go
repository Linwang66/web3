package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())
	r.POST("/api/chat", NewChat)
	r.POST("/api/upload", UploadImage)
	r.Run("127.0.0.1:8806")
}

func NewChat(c *gin.Context) {
	message := c.Query("message")
	qianfan.GetConfig().AccessKey = "ALTAKG4dhH4LYotEzsaN9y4Zc4"
	qianfan.GetConfig().SecretKey = "26b9a05bb59d433bab7599bc3bbc9d46"

	chat := qianfan.NewChatCompletion(
		qianfan.WithModel("ERNIE-Bot"),
	)

	resp, err := chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(message),
			},
		},
	)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	fmt.Println(resp.Result)
	c.JSON(http.StatusOK, gin.H{"message": resp.Result})
}

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	filePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// 将路径中的反斜杠替换为正斜杠
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded", "path": filePath})
}
