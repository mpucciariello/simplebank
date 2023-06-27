package client

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/utils"
	"net/http"
)

type (
	createUserReq struct {
		UserName string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	createUserRsp struct {
		UserName string `json:"username"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}

	getUserReq struct {
		UserName string `uri:"username" binding:"required,alphanum"`
	}

	deleteUserReq struct {
		UserName string `uri:"username" binding:"required,alphanum"`
	}
)

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	arg := db.CreateUserParams{
		UserName:       req.UserName,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
			}
			return

		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	} else {
		rsp := createUserRsp{
			UserName: user.UserName,
			FullName: user.FullName,
			Email:    user.Email,
		}
		ctx.JSON(http.StatusOK, rsp)
	}
}

func (s *Server) getUser(ctx *gin.Context) {
	var req getUserReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	} else {
		rsp := createUserRsp{
			UserName: user.UserName,
			FullName: user.FullName,
			Email:    user.Email,
		}
		ctx.JSON(http.StatusOK, rsp)
	}
}

//func (s *Server) deleteUser(ctx *gin.Context) {
//	var req deleteUserReq
//	if err := ctx.ShouldBindUri(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, errResponse(err))
//		return
//	}
//
//	err := s.store.DeleteUser(ctx, req.UserName)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, errResponse(err))
//	} else {
//		ctx.JSON(http.StatusOK, gin.H{"user": req.UserName})
//	}
//}
