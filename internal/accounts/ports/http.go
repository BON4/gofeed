package ports

import (
	"github.com/BON4/gofeed/internal/accounts/app"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	app app.Application
}

type logginParams struct{}

type logginResp struct{}

// @Summary      Logs user into the system
// @Description  logs in to account with user provided credantials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body  loginParams  true  "account login credentials"
// @Success      200     {object}  logginResp
// @Failure      default {object}  common
// @Router       /login [post]
func (h *HttpServer) Login(ctx *gin.Context) {}

func (h *HttpServer) Register(ctx *gin.Context) {}

// TODO: maby separate Account CRUD to another service

func (h *HttpServer) GetAccount(ctx *gin.Context) {}

func (h *HttpServer) CreateAccount(ctx *gin.Context) {}

func (h *HttpServer) UpdatePassword(ctx *gin.Context) {}

func (h *HttpServer) UpdateUsername(ctx *gin.Context) {}

func (h *HttpServer) DeleteAccount(ctx *gin.Context) {}

func (h *HttpServer) ListAccounts(ctx *gin.Context) {}
