package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"net/http"
)

type (
	createTransferReq struct {
		FromAccountID int64 `json:"from_account_id" binding:"required"`
		ToAccountID   int64 `json:"to_account_id" binding:"required"`
		Amount        int64 `json:"amount" binding:"required,min=1"`
	}
)

func (s *Server) createTranfer(ctx *gin.Context) {
	var req createTransferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// gin converts key-value error into a json
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	account, err := s.store.CreateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	} else {
		ctx.JSON(http.StatusOK, account)
	}
}
