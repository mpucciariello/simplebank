package client

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micaelapucciariello/simplebank/api/token"
	"net/http"
	"strings"
)

const _authorizationHeaderKey = "authorization"
const _authorizationTypeBearer = "Bearer"
const _authorizationPayloadKey = "authorization_payload"

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(_authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		authorizationType := fields[0]
		if authorizationType != _authorizationTypeBearer {
			err := fmt.Errorf("invalid authorization type: %v", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		ctx.Set(_authorizationPayloadKey, payload)
		ctx.Next()
	}
}
