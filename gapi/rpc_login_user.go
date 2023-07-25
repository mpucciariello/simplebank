package gapi

import (
	"context"
	"database/sql"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/pb"
	"github.com/micaelapucciariello/simplebank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "error invalid user: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "error getting user: %s", err)
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized user: %s", err)
	}

	accessToken, accessPayload, err := s.token.CreateToken(req.Username, s.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating  access token: %s", err)
	}

	refreshToken, refreshPayload, err := s.token.CreateToken(req.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating refresh token: %s", err)
	}

	mtdt := s.extractMetadata(ctx)
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    sql.NullTime{Time: refreshPayload.ExpiredAt},
	})

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		User:                  convertUser(user),
	}

	return rsp, nil
}
