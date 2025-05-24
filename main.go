package main

import (
	"YunDisk/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 路由分组简化
	fileGroup := r.Group("/file")
	{
		fileGroup.GET("/upload", gin.WrapF(handler.UpLoadView))
		fileGroup.POST("/upload", gin.WrapF(handler.UpLoadHandler))
		fileGroup.GET("/upload/suc", gin.WrapF(handler.UpLoadSucHandler))
		fileGroup.GET("/meta", gin.WrapF(handler.GetFileMetaHandler))
		fileGroup.GET("/download", gin.WrapF(handler.DownloadHandler))
		fileGroup.GET("/update", gin.WrapF(handler.FileMetaUpdateHandler))
		fileGroup.GET("/delete", gin.WrapF(handler.FileDeleteHandler))
	}
	// 启动服务
	err := r.Run(":8080")
	if err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
