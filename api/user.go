package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/utils"
	"net/http"
	"time"
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
		SessionID             uuid.UUID     `json:"session_id"`
		RefreshToken          string        `json:"refresh_token"`
		RefreshTokenExpiresAt time.Time     `json:"refresh_token_expires_at"`
		AccessToken           string        `json:"access_token"`
		AccessTokenExpiresAt  time.Time     `json:"access_token_expires_at"`
		UserMetadata          createUserRsp `json:"user_metadata"`
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
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, accessPayload, err := s.token.CreateToken(req.Username, s.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	refreshToken, refreshPayload, err := s.token.CreateToken(req.Username, s.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    sql.NullTime{Time: refreshPayload.ExpiredAt},
	})

	rsp := loginUserResponse{
		SessionID:             session.ID,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		UserMetadata:          parseUserInfo(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
