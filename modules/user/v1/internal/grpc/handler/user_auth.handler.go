// Package handler is a handler for user handler gRPC module
package handler

import (
	"context"
	"net/http"

	"google.golang.org/grpc/status"

	userv1 "grpc-starter/api/user/v1"
	"grpc-starter/common/config"
	"grpc-starter/common/constant"
	"grpc-starter/common/errors"
	"grpc-starter/modules/user/v1/service"
)

// UserHandler is a gRPC handler for the user auth service.
type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	config         config.Config
	userFinderSvc  service.UserFinderUseCase
	userCreatorSvc service.UserCreatorUseCase
	userUpdaterSvc service.UserUpdaterUseCase
	userDeleterSvc service.UserDeleterUseCase
}

// NewUserHandler returns a new UserHandler.
func NewUserHandler(
	config config.Config,
	userFinderSvc service.UserFinderUseCase,
	userCreatorSvc service.UserCreatorUseCase,
	userUpdaterSvc service.UserUpdaterUseCase,
	userDeleterSvc service.UserDeleterUseCase,
) *UserHandler {
	return &UserHandler{
		config:         config,
		userFinderSvc:  userFinderSvc,
		userCreatorSvc: userCreatorSvc,
		userUpdaterSvc: userUpdaterSvc,
		userDeleterSvc: userDeleterSvc,
	}
}

// Login define gRPC handler login for user modules
func (ah *UserHandler) Login(ctx context.Context, request *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, token, err := ah.userFinderSvc.Login(ctx, request.Email, request.Password)

	if err != nil {
		parseError := errors.ParseError(err)
		return nil, status.Errorf(
			parseError.Code,
			parseError.Message,
		)
	}

	return &userv1.LoginResponse{
		Code:    http.StatusOK,
		Message: constant.SuccessMessage,
		Data: &userv1.TokenData{
			UserId: user.ID.String(),
			Token:  token,
		},
	}, nil
}
