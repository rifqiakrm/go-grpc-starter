package builder

import (
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"grpc-starter/common/config"
	commonredis "grpc-starter/common/redis"
	"grpc-starter/modules/user/v1/internal/grpc/handler"
	"grpc-starter/modules/user/v1/internal/repository"
	"grpc-starter/modules/user/v1/service"
)

// BuildUserHandler builds grpc user handler with services and repositories
func BuildUserHandler(
	cfg config.Config,
	db *gorm.DB,
	redisPool *redis.Pool,
	grpcConn *grpc.ClientConn,
) *handler.UserHandler {
	// Cache
	cache := commonredis.NewClient(redisPool)

	// Repositories
	userFinderRepo := repository.NewUserFinderRepository(db, cache)
	userCreatorRepo := repository.NewUserCreatorRepository(db, cache)
	userUpdaterRepo := repository.NewUserUpdaterRepository(db, cache)
	userDeleterRepo := repository.NewUserDeleterRepository(db, cache)

	// Services
	userFinderSvc := service.NewUserFinder(cfg, userFinderRepo)
	userCreatorSvc := service.NewUserCreator(cfg, userCreatorRepo)
	userUpdaterSvc := service.NewUserUpdater(cfg, userUpdaterRepo)
	userDeleterSvc := service.NewUserDeleter(cfg, userDeleterRepo)

	return handler.NewUserHandler(
		cfg,
		userFinderSvc,
		userCreatorSvc,
		userUpdaterSvc,
		userDeleterSvc,
	)
}
