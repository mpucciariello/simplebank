package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"net/http"
	"time"
)

type (
	renewAccessTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	renewAccessTokenResponse struct {
		RefreshToken          string    `json:"refresh_token"`
		RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	}
)

func (s *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := s.token.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	session, err := s.store.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if session.Username != payload.UserName {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	refreshToken, payload, err := s.token.CreateToken(payload.UserName, s.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	session, err = s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           payload.ID,
		Username:     payload.UserName,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    sql.NullTime{Time: payload.ExpiredAt},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	rsp := renewAccessTokenResponse{
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: payload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
