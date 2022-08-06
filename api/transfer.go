package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/burakkarasel/Bank-App/token"
	"github.com/gin-gonic/gin"
)

// createAccountRequest holds the params of the request's and response's
type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency"  binding:"required,currency"`
}

// createAccount handles account creation requests, checks the binding, and finally if the account is succesfully inserted to DB
// returns http.StatusOK and account
func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.Currency, req.FromAccountID)

	if !valid {
		return
	}

	// here after checking fromAccount is valid or not we check if the fromAccount and authenticated user is same
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesnt belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.Currency, req.ToAccountID)

	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// validAccount checks if a given currency is valid for given account id
func (server *Server) validAccount(ctx *gin.Context, currency string, accID int64) (db.Account, bool) {
	acc, err := server.store.GetAccount(ctx, accID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return acc, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return acc, false
	}

	if acc.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", acc.ID, acc.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return acc, false
	}

	return acc, true
}
