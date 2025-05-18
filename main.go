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
		handler.UpLoadSucHandler(c.Writer, c.Request) // 修改这里
	})
	r.GET("/file/meta", func(c *gin.Context) {
		handler.GetFileMetaHandler(c.Writer, c.Request)
	})
	r.GET("/file/download", func(c *gin.Context) {
		handler.DownloadHandler(c.Writer, c.Request)
	})
	// 启动服务
	err := r.Run(":8080")
	if err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
