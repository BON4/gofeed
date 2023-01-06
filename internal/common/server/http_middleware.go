package server

import "github.com/gin-gonic/gin"

type MiddlewareProvider interface {
	Middleware() gin.HandlerFunc
}
