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
		Currency string `json:"currency" binding:"required,currency"`
	}

	getAccountReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	getAccountsListReq struct {
		PageID   int32 `form:"page_id" binding:"required,min=1"`
		PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	}

	deleteAccountReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
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
		return
	}

	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
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
		return
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

func (s *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	err := s.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	} else {
		ctx.JSON(http.StatusOK, gin.H{"account": req.ID})
	}
}
