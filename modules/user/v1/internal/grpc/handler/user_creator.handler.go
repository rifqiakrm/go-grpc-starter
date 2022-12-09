package handler

import (
	"context"
	"net/http"

	"google.golang.org/grpc/status"

	userv1 "grpc-starter/api/user/v1"
	"grpc-starter/common/constant"
	"grpc-starter/common/errors"
)

// Register handles the request to register a new user.
func (ah *UserHandler) Register(ctx context.Context, request *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	user, token, err := ah.userCreatorSvc.Register(ctx, request.GetUsername(), request.GetEmail(), request.GetPassword(), request.GetPhoneNumber())

	if err != nil {
		parseError := errors.ParseError(err)
		return nil, status.Errorf(
			parseError.Code,
			parseError.Message,
		)
	}

	return &userv1.RegisterResponse{
		Code:    http.StatusOK,
		Message: constant.SuccessMessage,
		Data: &userv1.TokenData{
			UserId: user.ID.String(),
			Token:  token,
		},
	}, nil
}
