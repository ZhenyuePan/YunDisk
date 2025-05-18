package main

import (
	"YunDisk/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/file/upload", func(c *gin.Context) {
		handler.UpLoadView(c.Writer, c.Request)
	})
	r.POST("/file/upload", func(c *gin.Context) {
		handler.UpLoadHandler(c.Writer, c.Request)
	})
	r.GET("/file/upload/suc", func(c *gin.Context) {
		handler.UpLoadHandler(c.Writer, c.Request)
	})
	r.GET("/file/meta", func(c *gin.Context) {
		handler.GetFileMetaHandler(c.Writer, c.Request)
	})
	// 启动服务
	err := r.Run(":8080")
	if err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
