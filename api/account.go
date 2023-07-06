package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/token"
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

	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	if authPayload.UserName != req.Owner {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("owner doesn't belong to the authenticated user"))
		return
	}
	arg := db.CreateAccountParams{
		Owner:    authPayload.UserName,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
			}

		}
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
	}
	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	if authPayload.UserName != account.Owner {
		err = fmt.Errorf("owner doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
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

	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	params := db.ListAccountsParams{
		Owner:  authPayload.UserName,
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

	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}
	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	if authPayload.UserName != account.Owner {
		err = fmt.Errorf("owner doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	err = s.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	} else {
		ctx.JSON(http.StatusOK, gin.H{"account": req.ID})
	}
}
