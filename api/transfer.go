package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
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
		// gin converts key-value error into a json
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}

	if !s.validCurrency(ctx, req.ToAccountID, req.Currency) || !s.validCurrency(ctx, req.FromAccountID, req.Currency) {
		return
	}
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	account, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	} else {
		ctx.JSON(http.StatusOK, account)
		return
	}
}

func (s *Server) validCurrency(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errResponse(err))
				return false
			}

			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return false
		} else {
			ctx.JSON(http.StatusOK, account)
		}
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%v] currency mismatched: account currency %v - transfer currency %v", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return false
	}

	return true
}
