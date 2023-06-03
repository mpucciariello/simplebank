package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"net/http"
)

type (
	createAccountReq struct {
		Owner    string `json:"owner" binding:"required"`
		Currency string `json:"currency" binding:"required" oneof:"USD, EUR, ARS"`
	}

	getAccountReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	getAccountsListReq struct {
		PageID   int32 `form:"page_id" binding:"required,min=1"`
		PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	}
)

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// gin converts key-value error into a json
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	} else {
		ctx.JSON(http.StatusOK, account)
	}
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}

	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	} else {
		ctx.JSON(http.StatusOK, account)
	}
}

// getAccountsList executes a paginated query
func (s *Server) getAccountsList(ctx *gin.Context) {
	var req getAccountsListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}

	params := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: req.PageID,
	}

	account, err := s.store.ListAccounts(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	} else {
		ctx.JSON(http.StatusOK, account)
	}
}
