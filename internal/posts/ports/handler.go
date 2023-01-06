package ports

import "github.com/gin-gonic/gin"

func MountHandlers(server *HttpServer, router *gin.RouterGroup) {
	authorized := router.Group("/")
	authorized.Use(server.SessionMiddleware.Middleware())
	{
		authorized.POST("/", server.CreatePost)
		authorized.DELETE("/:post_id", server.DeletePost)
	}

	router.GET("/list", server.ListPosts)
}
