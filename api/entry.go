package api

import (
	"database/sql"
	"net/http"

	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/burakkarasel/Bank-App/token"
	"github.com/gin-gonic/gin"
)

// createEntryRequest holds the inputs that are needed to make entry request
type createEntryRequest struct {
	AccountID int64 `json:"account_id" binding:"required,min=1"`
	Amount    int64 `json:"amount" binding:"required"`
}

// createEntry takes account ID and amount to make an entry
func (server *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// here we check the authenticated user and accounID is associated or not

	err := server.getAuthenticationValidation(ctx, req.AccountID)

	if err != nil {
		if err == ErrAccountIsNotAuthenticatedUsers {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.EntryTxParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	}

	result, err := server.store.EntryTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// getEntryRequest hold the ID of the entry user wants to get
type getEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getEntry gets the specified entry from the DB
func (server *Server) getEntry(ctx *gin.Context) {
	var req getEntryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// here we get the entry from DB
	entry, err := server.store.GetEntry(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// here we check the authenticated user and accounID is associated or not

	err = server.getAuthenticationValidation(ctx, entry.AccountID)

	if err != nil {
		if err == ErrAccountIsNotAuthenticatedUsers {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

// listEntriesRequest holds the query values for listEntries handler
type listEntriesRequest struct {
	AccountID int64 `form:"account_id" binding:"required,min=1"`
	PageID    int64 `form:"page_id" binding:"required,min=1"`
	PageSize  int64 `form:"page_size" binding:"required,min=5,max=10"`
}

// listEntries list entries for the user for a specific account
func (server *Server) listEntries(ctx *gin.Context) {
	var req listEntriesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// here we check the authenticated user and accounID is associated or not

	err := server.getAuthenticationValidation(ctx, req.AccountID)

	if err != nil {
		if err == ErrAccountIsNotAuthenticatedUsers {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.ListEntriesParams{
		AccountID: req.AccountID,
		Limit:     int32(req.PageSize),
		Offset:    int32((req.PageID - 1) * req.PageSize),
	}

	entries, err := server.store.ListEntries(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

// getAuthenticationValidation checks if auhtenticated user and account matches
func (server *Server) getAuthenticationValidation(ctx *gin.Context, accountID int64) error {
	acc, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		return err
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if acc.Owner != authPayload.Username {
		return ErrAccountIsNotAuthenticatedUsers
	}

	return nil
}
