package jwt

import (
	"context"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-starter/common/tools"
)

const (
	// userServiceLogin is the name of the service that is used for login.
	userServiceLogin = "/starter.user.v1.UserService/Login"
	// userServiceRegister is the name of the service that is used for register.
	userServiceRegister = "/starter.user.v1.UserService/Register"
	//nolint // userServiceForgotPassword is the name of the service that is used for forgot password.
	userServiceForgotPassword = "/starter.user.v1.UserService/ForgotPassword"
)

var ignoreMethod = []string{
	userServiceLogin,
	userServiceRegister,
	userServiceForgotPassword,
}

// CustomClaims define available data in JWT
type CustomClaims struct {
	ExpiresAt int64     `json:"exp,omitempty"`
	ID        string    `json:"jti,omitempty"`
	IssuedAt  int64     `json:"iat,omitempty"`
	NotBefore int64     `json:"nbf,omitempty"`
	Subject   uuid.UUID `json:"sub,omitempty"`
	Issuer    string    `json:"iss,omitempty"`
	jwt.StandardClaims
}

// Authorize is used by a middleware to authenticate requests
func Authorize(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	for _, v := range ignoreMethod {
		if method == v {
			return ctx, nil
		}
	}
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	userClaims, err := verifyClaims(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Sesi anda telah berakhir. Silahkan login kembali")
	}

	newCtx := context.WithValue(ctx, tools.ContextKeySubjectID, userClaims.ID)

	return newCtx, nil
}

// verifyClaims verifies the claims in the token.
func verifyClaims(accessToken string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
