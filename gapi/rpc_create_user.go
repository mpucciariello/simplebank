package gapi

import (
	"context"
	"github.com/lib/pq"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/pb"
	"github.com/micaelapucciariello/simplebank/utils"
	"github.com/micaelapucciariello/simplebank/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// req.GetField() checks if the field is ok
	if violations := validateCreateUserReq(req); violations != nil {
		return nil, InvalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error hashing password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.Internal, "username already exists: %s", err)
			}

		}
		return nil, status.Errorf(codes.Internal, "db err while creating user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func validateCreateUserReq(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.Username); err != nil {
		violations = append(violations, ViolationErr("username", err.Error()))
	}
	if err := validator.ValidateEmail(req.Email); err != nil {
		violations = append(violations, ViolationErr("email", err.Error()))
	}
	if err := validator.ValidatePassword(req.Password); err != nil {
		violations = append(violations, ViolationErr("password", err.Error()))
	}
	if err := validator.ValidateFullName(req.FullName); err != nil {
		violations = append(violations, ViolationErr("full_name", err.Error()))
	}

	return violations
}
