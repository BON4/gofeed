package httperr

import (
	"net/http"

	"github.com/BON4/gofeed/internal/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func InternalError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Internal server error", http.StatusInternalServerError)
}

func Unauthorised(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Unauthorised", http.StatusUnauthorized)
}

// func DoesNotExists(slug string, err error, w http.ResponseWriter, r *http.Request) {
// 	httpRespondWithError(err, slug, w, r, "Does Not Exists", http.StatusBadRequest)
// }

// func AlreadyExists(slug string, err error, w http.ResponseWriter, r *http.Request) {
// 	httpRespondWithError(err, slug, w, r, "Already Exists", http.StatusBadRequest)
// }

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Bad request", http.StatusBadRequest)
}

func RespondWithSlugError(err error, w http.ResponseWriter, r *http.Request) {
	slugError, ok := err.(errors.SlugError)
	if !ok {
		InternalError("internal-server-error", err, w, r)
		return
	}

	switch slugError.ErrorType() {
	case errors.ErrorTypeInvalidCred:
		Unauthorised(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeAlredyExists:
		BadRequest(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeDoesNotExists:
		BadRequest(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeAuthorization:
		Unauthorised(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeIncorrectInput:
		BadRequest(slugError.Slug(), slugError, w, r)
	default:
		InternalError(slugError.Slug(), slugError, w, r)
	}
}

func GinRespondWithSlugError(err error, ctx *gin.Context) {
	slugError, ok := err.(errors.SlugError)
	if !ok {
		InternalError("internal-server-error", err, ctx.Writer, ctx.Request)
		return
	}

	switch slugError.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorised(slugError.Slug(), slugError, ctx.Writer, ctx.Request)
	case errors.ErrorTypeIncorrectInput:
		BadRequest(slugError.Slug(), slugError, ctx.Writer, ctx.Request)
	default:
		InternalError(slugError.Slug(), slugError, ctx.Writer, ctx.Request)
	}
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, logMSg string, status int) {
	// TODO: log error
	resp := ErrorResponse{slug, status}

	if err := render.WriteJSON(w, resp); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
