package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/burakkarasel/Bank-App/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// addAuthorization creates a token and sets request's header with given authorizationType and the token
func addAuthorization(t *testing.T, r *http.Request, tokenMaker token.Maker, authorizationType, username string, duration time.Duration) {
	// here we create token with given input
	token, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// and here we set request's header with generated token, and given authorizationType
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	r.Header.Set(authorizationHeaderKey, authorizationHeader)
}

// TestAuthMiddleware tests authMiddleware
func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "No Authorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// without authorization
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unsupported Authorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// here i pass unsupported as authorizationType which will fail because unsupported Authorization Type
				addAuthorization(t, request, tokenMaker, "unsupported", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Authorization Format",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// here i pass no header so authorization field won't be long enough
				addAuthorization(t, request, tokenMaker, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Expired token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// here i pass -time.Minute
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// here we dont need store because we are just gonna test our middleware
			server := newTestServer(t, nil)

			// here we create a new route for the request
			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			// here we create a new recorder and a request
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			// then we setup our middleware and make a request
			tt.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}
