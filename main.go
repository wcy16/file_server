package main

import (
	"file_server/api"
	"file_server/middleware"
	"github.com/gin-gonic/gin"
)

func FileServerRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()

	router.StaticFile("not_found", "./default/default.jpg")
	router.StaticFile("file_not_found", "./default/default.jpg")
	router.StaticFile("error", "./default/error.png")

	router.POST("/file", middleware.CheckUploadPermission(), api.Upload)
	router.GET("/file/:filename", middleware.CheckDownloadPermission(), api.Download)

	return router
}

func main() {
	r := FileServerRouter()
	r.Run(":8081")
}
