package ports

import (
	"time"

	"github.com/BON4/gofeed/internal/common/errors"
	"github.com/BON4/gofeed/internal/common/server/httperr"
	"github.com/BON4/gofeed/internal/posts/app"
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

type createPostRequest struct {
	Content string `json:"content"`
}

// @Summary     Create Post
// @Description Creates post if user have permision
// @Tags        posts
// @Produce     json
// @Param       input   body  createPostRequest  true  "information for post creation"
// @Success     201
// @Failure     default {object}  httperr.ErrorResponse
// @Router      /api [post]
func (h *HttpServer) CreatePost(ctx *gin.Context) {
	httperr.GinRespondWithSlugError(errors.NewNotImplementedError("Endpopint not ipmlemented", "edp-not-implemented"), ctx)
}

// @Summary     Delete Post
// @Description Deletes post if user have permision
// @Tags        posts
// @Produce     json
// @Param       post_id   path  int64  true  "account id"
// @Success     200
// @Failure     default {object}  httperr.ErrorResponse
// @Router      /api/{post_id} [delete]
func (h *HttpServer) DeletePost(ctx *gin.Context) {
	httperr.GinRespondWithSlugError(errors.NewNotImplementedError("Endpopint not ipmlemented", "edp-not-implemented"), ctx)
}

type getPostResponse struct {
	PostId   int64     `json:"post_id"`
	Content  string    `json:"content"`
	PostedOn time.Time `json:"posted_on"`
	PostedBy string    `json:"posted_by"`
	Score    int       `json:"score"`
}

// @Summary     List
// @Description Retrives list of json formated objects
// @Tags        posts
// @Produce     json
// @Param       page_size         query     int              true "page size"
// @Param       page_number       query     int              true "page number"
// @Success     200     {array}   getPostResponse
// @Failure     default {object}  httperr.ErrorResponse
// @Router      /api/list [get]
func (h *HttpServer) ListPosts(ctx *gin.Context) {
	httperr.GinRespondWithSlugError(errors.NewNotImplementedError("Endpopint not ipmlemented", "edp-not-implemented"), ctx)
}
