package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func RunHTTPServer(mountHandlers func(router *gin.RouterGroup), mountDocs func(router *gin.RouterGroup)) {
	port := ""
	if port = os.Getenv("PORT"); port == "" {
		// TODO: change to log
		fmt.Println("Not found PORT, runing on default 8080")
		port = "8080"
	}

	RunHTTPServerOnAddr(":"+port, mountHandlers, mountDocs)
}

func RunHTTPServerOnAddr(addr string, mountHandlers func(router *gin.RouterGroup), mountDocs func(router *gin.RouterGroup)) {
	eng := gin.Default()
	rootRoutes := eng.Group("")

	apiRoutes := rootRoutes.Group("/api")

	mountHandlers(apiRoutes)

	docsRoutes := rootRoutes.Group("/docs")
	mountDocs(docsRoutes)

	srv := &http.Server{
		Handler: eng,
		Addr:    addr,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
			return
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}
