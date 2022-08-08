package api

import (
	"database/sql"
	"net/http"

	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/burakkarasel/Bank-App/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// createAccountRequest holds the params of the request's and response's
type createAccountRequest struct {
	Currency string `json:"currency"  binding:"required,currency"`
}

// createAccount handles account creation requests, checks the binding, and finally if the account is succesfully inserted to DB
// returns http.StatusOK and account
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// here we prevent getting different owners to create account instead users can only create account with their username
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// getAccountByIdRequest holds the data from request's URI
type getAccountByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccountById checks for URI bindings, then gets the account from DB and checks for any error
func (server *Server) getAccountById(ctx *gin.Context) {
	var req getAccountByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// here we prevent users to check other user's accounts
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrAccountIsNotAuthenticatedUsers))
	}

	ctx.JSON(http.StatusOK, account)
}

// ListAccountsRequest holds the query params for the listAccounts handler
type ListAccountsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAccounts returns a accounts as specified in the query of the URI
func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccountsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// here we prevent user's to list other user's accounts
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
