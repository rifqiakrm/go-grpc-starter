// Package fsm is finite state management module.
package fsm

import (
	"context"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	userv1 "grpc-starter/api/user/v1"
	"grpc-starter/common/config"
	"grpc-starter/modules/user/v1/internal/builder"
)

// InitGrpc initializes gRPC user modules.
func InitGrpc(server *grpc.Server, cfg config.Config, db *gorm.DB, redisPool *redis.Pool, grpcConn *grpc.ClientConn) {
	user := builder.BuildUserHandler(cfg, db, redisPool, grpcConn)
	userv1.RegisterUserServiceServer(server, user)
}

// InitRest initializes REST user modules.
// If any error occurs, it logs the error and continue the process.
func InitRest(ctx context.Context, server *runtime.ServeMux, grpcPort string, options ...grpc.DialOption) {
	if err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, server, grpcPort, options); err != nil {
		log.Printf("RegisterUserServiceHandlerFromEndpoint failed to be registered: %v", err)
	}
}
