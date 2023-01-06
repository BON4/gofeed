package session

import (
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

func (s *SessionMiddleware) Middleware() gin.HandlerFunc {
	// TODO: Add mock auth for testing purpeses
	// 	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); mockAuth {
	// 		return mock_auth_middleware
	// }
	return func(ctx *gin.Context) {
		s.logger.Info("Trying to get token from:", s.headerKey)
		token := ctx.GetHeader(s.headerKey)
		if len(token) == 0 {
			httperr.GinRespondWithSlugError(errors.NewAuthorizationError("token not provided", "mdwr-not-provided-token"), ctx)
			return
		}

		s.logger.Infof("got token: %s", token)
		payload, err := s.tokenFc.VerifyToken(token)
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
