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

	loginUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	loginUserResponse struct {
		AccessToken  string        `json:"access_token"`
		UserMetadata createUserRsp `json:"user_metadata"`
	}
)

func parseUserInfo(user db.User) createUserRsp {
	return createUserRsp{
		UserName: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
}

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
		Username:       req.UserName,
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
		rsp := parseUserInfo(user)
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
		rsp := parseUserInfo(user)
		ctx.JSON(http.StatusOK, rsp)
	}
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	} else {
		rsp := createUserRsp{
			UserName: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		}
		ctx.JSON(http.StatusOK, rsp)
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, err := s.token.CreateToken(req.Username, s.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	rsp := loginUserResponse{
		AccessToken:  accessToken,
		UserMetadata: parseUserInfo(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
