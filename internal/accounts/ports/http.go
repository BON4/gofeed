package ports

import (
	"time"

	"github.com/BON4/gofeed/internal/accounts/app"
	"github.com/BON4/gofeed/internal/accounts/app/usecase"
	"github.com/BON4/gofeed/internal/common/errors"
	"github.com/BON4/gofeed/internal/common/server/httperr"
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
	Username int64           `json:"username"`
	Email    string          `json:"email" format:"email"`
	Role     string          `json:"role" enums:"admin, basic"`
	Users    []*userResponse `json:"users"`
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

	err := h.app.LoginAccount.Handle(ctx.Request.Context(), usecase.LoginCommand{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
	}

	// TODO: create access and refresh token. Save it to go-cache
}

type registerParams struct {
	Username string `json:"username" minLength:"1" validate:"required"`
	Email    string `json:"email" validate:"required" format:"password"`
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
