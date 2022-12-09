package service

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"grpc-starter/common/config"
	"grpc-starter/common/constant"
	commonError "grpc-starter/common/errors"
	commonJwt "grpc-starter/common/jwt"
	"grpc-starter/modules/user/v1/entity"
	"grpc-starter/modules/user/v1/internal/repository"
)

// UserCreator responsible for creating user
type UserCreator struct {
	cfg                   config.Config
	userCreatorRepository repository.UserCreatorRepositoryUseCase
}

// UserCreatorUseCase is use case for creating existing user
type UserCreatorUseCase interface {
	// Create creates user
	Create(ctx context.Context, user *entity.User) error
	// Register creates user and send email to user
	Register(ctx context.Context, username string, email string, password string, phoneNumber string) (*entity.User, string, error)
}

// NewUserCreator constructs new instance of UserCreator
func NewUserCreator(
	cfg config.Config,
	userCreatorRepository repository.UserCreatorRepositoryUseCase,
) *UserCreator {
	return &UserCreator{
		cfg:                   cfg,
		userCreatorRepository: userCreatorRepository,
	}
}

// Create creates user
func (svc *UserCreator) Create(ctx context.Context, user *entity.User) error {
	if err := svc.userCreatorRepository.Create(ctx, user); err != nil {
		log.Print("[UserCreator - Create] Error while creating user data :", err)
		return commonError.ErrInternalServerError.Error()
	}

	return nil
}

// Register creates user and send email to user returns user and token
func (svc *UserCreator) Register(ctx context.Context, username string, email string, password string, phoneNumber string) (*entity.User, string, error) {
	newUser := entity.NewUser(
		uuid.New(),
		username,
		email,
		password,
		phoneNumber,
		"system",
	)

	if err := svc.userCreatorRepository.Create(ctx, newUser); err != nil {
		log.Print("[UserCreator - Register] Error while creating user data :", err)
		return nil, "", commonError.ErrInternalServerError.Error()
	}

	claims := &commonJwt.CustomClaims{
		ExpiresAt: time.Now().Add(time.Hour * constant.TwentyFourHour * constant.DaysInOneYear).Unix(),
		ID:        uuid.New().String(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Subject:   newUser.ID,
		Issuer:    constant.MobileIssuer,
	}

	gen := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := gen.SignedString([]byte(svc.cfg.JWTConfig.SecretKey))

	if err != nil {
		log.Print("[UserCreator - Register] Error while generating token for user :", err)
		return nil, "", commonError.ErrInternalServerError.Error()
	}

	return newUser, token, nil
}
