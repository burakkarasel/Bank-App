package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrSessionBlocked = errors.New("blocked session")
var ErrSessionUserIsInvalid = errors.New("incorrect session user")
var ErrInvalidToken = errors.New("mismatched session token")
var ErrExpiredSession = errors.New("expired session")

// renewAccessTokenRequest holds the params of the request's
type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token " binding:"required"`
}

// renewAccessTokenResponse holds access token and user
type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// renewAccessToken renews the acces token after checking it
func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// here we get the payload of the refresh token
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// after checking bindings we check for user from DB
	session, err := server.store.GetSession(ctx, refreshPayload.ID)

	if err != nil {
		// if err is sql.ErrNoRows refresh token ID is invalid
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// we check if the session is blocked or not
	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrSessionBlocked))
		return
	}

	// we check if the username in the session matches with the session in the DB
	if session.Username != refreshPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrSessionUserIsInvalid))
		return
	}

	// we check if the token from session matches with token in the request
	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrInvalidToken))
		return
	}

	// then we check if the session token is expired or not
	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrExpiredSession))
		return
	}

	// then we create an access token for this logged in user
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.RefreshTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// and then we send access token and it's expiration time
	resp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
