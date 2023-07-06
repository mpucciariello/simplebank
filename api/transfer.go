package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/token"
	"net/http"
)

type (
	createTransferReq struct {
		FromAccountID int64  `json:"from_account_id" binding:"required"`
		ToAccountID   int64  `json:"to_account_id" binding:"required"`
		Amount        int64  `json:"amount" binding:"required,min=1"`
		Currency      string `json:"currency" binding:"required,currency"`
	}
)

func (s *Server) createTranfer(ctx *gin.Context) {
	var req createTransferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, isValidFromAccount := s.validAccountCurrency(ctx, req.FromAccountID, req.Currency)
	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	if account.Owner != authPayload.UserName {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("from account doesn't belong to the authenticated user"))
	}

	_, isValidToAccount := s.validAccountCurrency(ctx, req.ToAccountID, req.Currency)

	if !isValidToAccount || !isValidFromAccount {
		return
	}
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	} else {
		ctx.JSON(http.StatusOK, transfer)
		return
	}
}

func (s *Server) validAccountCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errResponse(err))
				return account, false
			}

			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return account, false
		} else {
			ctx.JSON(http.StatusOK, account)
		}
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%v] currency mismatched: account currency %v - transfer currency %v", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return account, false
	}

	return account, true
}
