package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/burakkarasel/Bank-App/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware is a middleware that checks if a request is from an authorized user
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// here we first check for authorizationHeaderKey
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		// if len of the header is 0 then no header is avaiable, authorization fails
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// if len of the header is less than 2 then header is not a valid format
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])

		// if first piece of header is not our authorizationType than authorization fails
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)

		// if we cant verify token (invalid token, expired token...) and get a payload authorization fails
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// finally after authorization completes succesfully we put authorization key to context and move forward
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
