package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/BON4/gofeed/internal/accounts/ports"
	"github.com/gin-gonic/gin"
)

type AccountsServer struct {
}

// Deletes an account
// (DELETE /account)
func (s *AccountsServer) DeleteAccount(c *gin.Context, params ports.DeleteAccountParams) {
	panic("not implemented") // TODO: Implement
}

// Get account by username
// (GET /account)
func (s *AccountsServer) GetAccount(c *gin.Context, params ports.GetAccountParams) {
	panic("not implemented") // TODO: Implement
}

// Creates an account
// (POST /account)
func (s *AccountsServer) CreateAccount(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// Updates an account
// (PUT /account)
func (s *AccountsServer) UpdateAccount(c *gin.Context, params ports.UpdateAccountParams) {
	panic("not implemented") // TODO: Implement
}

// Get list of accounts
// (GET /account/list)
func (s *AccountsServer) ListAccount(c *gin.Context, params ports.ListAccountParams) {
	panic("not implemented") // TODO: Implement
}

// Logs user into the system
// (POST /auth/login)
func (s *AccountsServer) LoginUser(c *gin.Context, params ports.LoginUserParams) {
	panic("not implemented") // TODO: Implement
}

// Registers user into the system
// (POST /auth/register)
func (s *AccountsServer) RegisterUser(c *gin.Context, params ports.RegisterUserParams) {
	panic("not implemented") // TODO: Implement
}

func main() {
	r := gin.Default()

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
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
