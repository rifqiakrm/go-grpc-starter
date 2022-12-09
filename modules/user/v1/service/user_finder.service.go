package service

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"grpc-starter/common/config"
	"grpc-starter/common/constant"
	"grpc-starter/common/errors"
	commonJwt "grpc-starter/common/jwt"
	"grpc-starter/common/tools"
	"grpc-starter/modules/user/v1/entity"
	"grpc-starter/modules/user/v1/internal/repository"
)

// UserFinder responsible for finding user
type UserFinder struct {
	cfg                  config.Config
	userFinderRepository repository.UserFinderRepositoryUseCase
}

// UserFinderUseCase is use case for finding existing user
type UserFinderUseCase interface {
	// FindByID finds user by user id
	FindByID(ctx context.Context, refID uuid.UUID) (*entity.User, error)
	// Login finds user by email and password and generates token returns user and token
	Login(ctx context.Context, email string, password string) (*entity.User, string, error)
}

// NewUserFinder constructs new instance of UserFinder
func NewUserFinder(
	cfg config.Config,
	userFinderRepository repository.UserFinderRepositoryUseCase,
) *UserFinder {
	return &UserFinder{
		cfg:                  cfg,
		userFinderRepository: userFinderRepository,
	}
}

// FindByID finds user by user id
func (svc *UserFinder) FindByID(ctx context.Context, refID uuid.UUID) (*entity.User, error) {
	res, err := svc.userFinderRepository.FindByID(ctx, refID)

	if err != nil {
		log.Println("[UserFinder - FindByID] Error while finding user data :", err)
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrRecordNotFound.Error()
		}
		return nil, err
	}

	return res, nil
}

// Login finds user by email and password and generates token returns user and token
func (svc *UserFinder) Login(ctx context.Context, email string, password string) (*entity.User, string, error) {
	res, err := svc.userFinderRepository.FindByEmail(ctx, email)

	if err != nil {
		log.Println("[UserFinder - FindByEmailPassword] Error while finding user data :", err)
		if err == gorm.ErrRecordNotFound {
			return nil, "", errors.ErrRecordNotFound.Error()
		}
		return nil, "", err
	}

	verifyPassword := tools.BcryptVerifyHash(res.Password, password)

	if !verifyPassword {
		return nil, "", errors.ErrWrongLoginCredentials.Error()
	}

	claims := &commonJwt.CustomClaims{
		ExpiresAt: time.Now().Add(time.Hour * constant.TwentyFourHour * constant.DaysInOneYear).Unix(),
		ID:        uuid.New().String(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Subject:   res.ID,
		Issuer:    constant.MobileIssuer,
	}

	gen := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := gen.SignedString([]byte(svc.cfg.JWTConfig.SecretKey))

	if err != nil {
		log.Println("[UserFinder - Login] Error while generating token :", err)
		return nil, "", errors.ErrInternalServerError.Error()
	}

	return res, token, nil
}
