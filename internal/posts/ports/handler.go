package ports

import "github.com/gin-gonic/gin"

func MountHandlers(server *HttpServer, router *gin.RouterGroup) {
	router.POST("/", server.CreatePost)
	router.DELETE("/:post_id", server.DeletePost)
	router.GET("/list", server.ListPosts)
}
