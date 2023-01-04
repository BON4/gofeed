package ports

import (
	"net/http"
	"time"

	"github.com/BON4/gofeed/internal/accounts/app"
	"github.com/BON4/gofeed/internal/accounts/app/usecase"
	"github.com/BON4/gofeed/internal/common/errors"
	"github.com/BON4/gofeed/internal/common/server/httperr"
	"github.com/BON4/gofeed/internal/common/tokens"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	app *app.Application
}

func NewHttpServer(app *app.Application) *HttpServer {
	return &HttpServer{
		app: app,
	}
}

type loginParams struct {
	Username string `json:"username" minLength:"1" validate:"required"`
	Password string `json:"password" minLength:"4" validate:"required" format:"password"`
}

type userResponse struct {
	UUID    string `json:"uuid" format:"uuid"`
	IP      string `json:"ip" format:"ipv4"`
	Browser string `json:"browser"`
}

type accountResponse struct {
	Username string `json:"username"`
	Role     string `json:"role" enums:"admin, basic"`
}

type loginResponse struct {
	AccessToken           string          `json:"access_token"`
	AccessTokenExpiresAt  time.Time       `json:"access_token_expires_at" format:"date-time"`
	RefreshToken          string          `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time       `json:"refresh_token_expires_at" format:"date-time"`
	Account               accountResponse `json:"account"`
}

// @Summary      Logs user into the system
// @Description  logs in to account with user provided credantials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body  loginParams  true  "account login credentials"
// @Success      200     {object}  loginResponse
// @Failure      default {object}  httperr.ErrorResponse
// @Router       /api/login [post]
func (h *HttpServer) Login(ctx *gin.Context) {
	var req loginParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	resp, err := h.app.LoginAccount.Handle(ctx.Request.Context(), usecase.LoginQuery{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
	}

	ctx.JSON(http.StatusAccepted, loginResponse{
		AccessToken:           resp.AccessToken,
		AccessTokenExpiresAt:  resp.AccessTokenExpiresAt,
		RefreshToken:          resp.RefreshToken,
		RefreshTokenExpiresAt: resp.RefreshTokenExpiresAt,
		Account: accountResponse{
			Username: resp.Instance.Username,
			Role:     resp.Instance.Role,
		},
	})

	// TODO: handle error
	h.app.CreateSession.Handle(ctx, usecase.CreateSessionCommand{
		ID:           resp.RefreshTokenId,
		Refreshtoken: resp.RefreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    resp.RefreshTokenExpiresAt,
		Instance: tokens.InstanceCredentials{
			Username: resp.Instance.Username,
			Role:     resp.Instance.Role,
		},
	})
}

type registerParams struct {
	Username string `json:"username" minLength:"1" validate:"required"`
	Email    string `json:"email" validate:"required" format:"email"`
	Password string `json:"password" minLength:"4" validate:"required" format:"password"`
}

// @Summary      Registers user into the system
// @Description  registers new account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body      registerParams  true  "register credantials"
// @Success      201
// @Failure      default {object}  httperr.ErrorResponse
// @Router       /api/register [post]
func (h *HttpServer) Register(ctx *gin.Context) {
	var req registerParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	err := h.app.RegisterAccount.Handle(ctx.Request.Context(), usecase.RegisterCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
	}

}

// TODO: maby separate Account CRUD to another service

func (h *HttpServer) GetAccount(ctx *gin.Context) {
	httperr.GinRespondWithSlugError(errors.NewNotImplementedError("Endpopint not ipmlemented", "edp-not-implemented"), ctx)
}

func (h *HttpServer) CreateAccount(ctx *gin.Context) {
	httperr.GinRespondWithSlugError(errors.NewNotImplementedError("Endpopint not ipmlemented", "edp-not-implemented"), ctx)
}

func (h *HttpServer) UpdatePassword(ctx *gin.Context) {
	httperr.GinRespondWithSlugError(errors.NewNotImplementedError("Endpopint not ipmlemented", "edp-not-implemented"), ctx)
}

func (h *HttpServer) UpdateUsername(ctx *gin.Context) {
	httperr.GinRespondWithSlugError(errors.NewNotImplementedError("Endpopint not ipmlemented", "edp-not-implemented"), ctx)
}

func (h *HttpServer) DeleteAccount(ctx *gin.Context) {
	httperr.GinRespondWithSlugError(errors.NewNotImplementedError("Endpopint not ipmlemented", "edp-not-implemented"), ctx)
}

func (h *HttpServer) ListAccounts(ctx *gin.Context) {
	httperr.GinRespondWithSlugError(errors.NewNotImplementedError("Endpopint not ipmlemented", "edp-not-implemented"), ctx)
}
