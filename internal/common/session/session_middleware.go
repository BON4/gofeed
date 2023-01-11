package session

import (
	"fmt"
	"strings"

	"github.com/BON4/gofeed/internal/common/errors"
	"github.com/BON4/gofeed/internal/common/server/httperr"
	"github.com/BON4/gofeed/internal/common/session/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SessionMiddleware struct {
	store     domain.Store
	tokenFc   *domain.TokenVerifier
	logger    *logrus.Entry
	headerKey string
}

func NewSessionMiddleware(
	store domain.Store,
	tokenFc *domain.TokenVerifier,
	logger *logrus.Entry,
	headerKey string,
) *SessionMiddleware {
	return &SessionMiddleware{
		store:     store,
		tokenFc:   tokenFc,
		logger:    logger,
		headerKey: headerKey,
	}
}

const authorizationTypeBearer = "bearer"

func (s *SessionMiddleware) Middleware() gin.HandlerFunc {
	// TODO: Add mock auth for testing purpeses
	// 	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); mockAuth {
	// 		return mock_auth_middleware
	// }
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(s.headerKey)

		if len(authorizationHeader) == 0 {
			httperr.GinRespondWithSlugError(errors.NewAuthorizationError("authorization header is not provided", "no-token-header"), ctx)
			return

		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			httperr.GinRespondWithSlugError(errors.NewAuthorizationError("invalid authorization header format", "invalid-token-header"), ctx)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			httperr.GinRespondWithSlugError(err, ctx)
			return
		}

		accessToken := fields[1]
		payload, err := s.tokenFc.VerifyToken(accessToken)
		if err != nil {
			httperr.GinRespondWithSlugError(err, ctx)
			return
		}

		if err := payload.Valid(); err != nil {
			httperr.GinRespondWithSlugError(err, ctx)
			return
		}

		ctx.Set("payload", payload)
		ctx.Next()
	}
}
