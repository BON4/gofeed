package ports

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BON4/gofeed/internal/common/errors"
	"github.com/BON4/gofeed/internal/common/server/httperr"
	"github.com/BON4/gofeed/internal/common/session"
	"github.com/BON4/gofeed/internal/common/tokens"
	"github.com/BON4/gofeed/internal/posts/app"
	"github.com/BON4/gofeed/internal/posts/app/usecase"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	app               *app.Application
	SessionMiddleware *session.SessionMiddleware
}

func NewHttpServer(app *app.Application, SessionMiddleware *session.SessionMiddleware) *HttpServer {
	return &HttpServer{
		app:               app,
		SessionMiddleware: SessionMiddleware,
	}
}

type createPostRequest struct {
	Content string `json:"content"`
}

// @Summary     Create Post
// @Security    JWT
// @Description Creates post if user have permision
// @Tags        posts
// @Produce     json
// @Param       input   body  createPostRequest  true  "information for post creation"
// @Success     201
// @Failure     default {object}  httperr.ErrorResponse
// @Router      /api [post]
func (h *HttpServer) CreatePost(ctx *gin.Context) {
	fmt.Println("Got req")
	payload, err := tokens.GetPayloadFromContext(ctx, "payload")
	if err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	var req createPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	_, err = h.app.CreatePost.Handle(ctx.Request.Context(), usecase.CreatePostQuery{
		Content: req.Content,
		Account: payload.Instance.Username,
	})

	if err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

// @Summary     Upwote Post
// @Security    JWT
// @Description Upwotes Post and increments its score.
// @Tags        posts
// @Produce     json
// @Param       post_id   path  int64  true  "account id"
// @Success     200
// @Failure     default {object}  httperr.ErrorResponse
// @Router      /api/up/{post_id} [put]
func (h *HttpServer) UpwotePost(ctx *gin.Context) {
	payload, err := tokens.GetPayloadFromContext(ctx, "payload")
	if err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	reqId, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		httperr.GinRespondWithSlugError(errors.NewIncorrectInputError(err.Error(), "error-parsing-param"), ctx)
		return
	}

	if err := h.app.RatePost.Handle(ctx.Request.Context(), usecase.PostRateParams{
		PostId:  reqId,
		Account: payload.Instance.Username,
		Rate:    1,
	}); err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// @Summary     Downwote Post
// @Security    JWT
// @Description Downwotes Post and increments its score.
// @Tags        posts
// @Produce     json
// @Param       post_id   path  int64  true  "account id"
// @Success     200
// @Failure     default {object}  httperr.ErrorResponse
// @Router      /api/down/{post_id} [put]
func (h *HttpServer) DownwotePost(ctx *gin.Context) {
	payload, err := tokens.GetPayloadFromContext(ctx, "payload")
	if err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	reqId, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		httperr.GinRespondWithSlugError(errors.NewIncorrectInputError(err.Error(), "error-parsing-param"), ctx)
		return
	}

	if err := h.app.RatePost.Handle(ctx.Request.Context(), usecase.PostRateParams{
		PostId:  reqId,
		Account: payload.Instance.Username,
		Rate:    -1,
	}); err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// @Summary     Delete Post
// @Security    JWT
// @Description Deletes post if user have permision
// @Tags        posts
// @Produce     json
// @Param       post_id   path  int64  true  "account id"
// @Success     200
// @Failure     default {object}  httperr.ErrorResponse
// @Router      /api/{post_id} [delete]
func (h *HttpServer) DeletePost(ctx *gin.Context) {
	payload, err := tokens.GetPayloadFromContext(ctx, "payload")
	if err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	// TODO: this is stupid. Role is defined in accounts.domain
	if payload.Instance.Role != "admin" {
		httperr.GinRespondWithSlugError(errors.NewAuthorizationError("only admins can delete post", "post-delete-not-allowed"), ctx)
		return
	}

	reqId, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		httperr.GinRespondWithSlugError(errors.NewIncorrectInputError(err.Error(), "error-parsing-param"), ctx)
		return
	}

	if err := h.app.DeletePost.Handle(ctx.Request.Context(), usecase.DeletePostCommand{
		PostId: reqId,
	}); err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type getPostResponse struct {
	PostId   int64     `json:"post_id"`
	Content  string    `json:"content"`
	PostedOn time.Time `json:"posted_on"`
	PostedBy string    `json:"posted_by"`
	Score    int       `json:"score"`
}

// @Summary     List
// @Security     JWT
// @Description Retrives list of json formated objects
// @Tags        posts
// @Produce     json
// @Param       page_size         query     int              true "page size"
// @Param       page_number       query     int              true "page number"
// @Success     200     {array}   getPostResponse
// @Failure     default {object}  httperr.ErrorResponse
// @Router      /api/list [get]
func (h *HttpServer) ListPosts(ctx *gin.Context) {
	form := usecase.FindPostParams{}
	if err := ctx.Bind(&form); err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	posts, err := h.app.ListPost.Handle(ctx, form)
	if err != nil {
		httperr.GinRespondWithSlugError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}
