package ports

import "github.com/gin-gonic/gin"

func MountHandlers(server *HttpServer, router *gin.RouterGroup) {
	router.POST("/login", server.Login)
	router.POST("/register", server.Register)
}
