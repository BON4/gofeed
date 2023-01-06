package main

import (
	_ "github.com/BON4/gofeed/internal/posts/api/openapi"

	"github.com/BON4/gofeed/internal/common/server"
	"github.com/BON4/gofeed/internal/posts/config"
	"github.com/BON4/gofeed/internal/posts/ports"
	"github.com/BON4/gofeed/internal/posts/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middlewar
)

// @title           Telegram Subs API
// @version         1.0
// @description     This service provide functionality for storing and managing privat telegram channels with subscription based payments for acessing content.

// @host      localhost:8081
// @BasePath  /

// @securityDefinitions.apiKey JWT
// @in header
// @name authorization
func main() {
	cfg, err := config.LoadServerConfig(".")
	if err != nil {
		panic(err)
	}

	application := service.NewApplication(cfg)
	server.RunHTTPServer(
		func(router *gin.RouterGroup) {
			ports.MountHandlers(
				ports.NewHttpServer(
					application,
					service.NewSessionMiddleware(cfg),
				),
				router,
			)
		},
		func(router *gin.RouterGroup) {
			router.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		})
}
