package main

import (
	"os"
	"strings"

	_ "github.com/BON4/gofeed/internal/accounts/api/openapi"
	"github.com/BON4/gofeed/internal/accounts/config"
	"google.golang.org/grpc"

	"github.com/BON4/gofeed/internal/accounts/ports"
	"github.com/BON4/gofeed/internal/accounts/service"
	"github.com/BON4/gofeed/internal/common/genproto/accounts"
	"github.com/BON4/gofeed/internal/common/server"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middlewar
)

// @title           Telegram Subs API
// @version         1.0
// @description     This service provide functionality for storing and managing privat telegram channels with subscription based payments for acessing content.

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apiKey JWT
// @in header
// @name authorization
func main() {
	cfg, err := config.LoadServerConfig(".")
	if err != nil {
		panic(err)
	}

	application, appCleanup := service.NewApplication(cfg)

	defer func() {
		appCleanup()
	}()

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		server.RunHTTPServer(
			func(router *gin.RouterGroup) {
				ports.MountHandlers(
					ports.NewHttpServer(application),
					router,
				)
			},
			func(router *gin.RouterGroup) {
				router.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
			})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			srvc := ports.NewGrpcServer(application)
			accounts.RegisterAccountServiceServer(server, srvc)
		})
	}
}
