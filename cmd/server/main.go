package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"grpc-starter/common/config"
	gormConn "grpc-starter/common/gorm"
	"grpc-starter/common/healthcheck"
	"grpc-starter/common/postgres"
	commonRedis "grpc-starter/common/redis"
	notificationModules "grpc-starter/modules/notification/v1"
	userModules "grpc-starter/modules/user/v1"
	pubsubSDK "grpc-starter/sdk/pubsub"
	"grpc-starter/server"
)

const (
	envDevelopment     = "development"
	maxCallRecvMsgSize = 20000000
	version            = "1.0.0"
)

// splash prints out the splash screen
func splash(cfg *config.Config) {
	colorReset := "\033[0m"
	colorBlue := "\033[34m"
	colorCyan := "\033[36m"

	fmt.Printf(`
                                        __                 __                
   _________________   ____     _______/  |______ ________/  |_  ___________ 
  / ___\_  __ \____ \_/ ___\   /  ___/\   __\__  \\_  __ \   __\/ __ \_  __ \
 / /_/  >  | \/  |_> >  \___   \___ \  |  |  / __ \|  | \/|  | \  ___/|  | \/
 \___  /|__|  |   __/ \___  > /____  > |__| (____  /__|   |__|  \___  >__|   
/_____/       |__|        \/       \/            \/                 \/       v%s
`, version)

	fmt.Println(colorBlue, fmt.Sprintf(`⇨ REST server started on :%s`, cfg.Port.REST))
	fmt.Println(colorCyan, fmt.Sprintf(`⇨ GRPC server started on :%s`, cfg.Port.GRPC))
	fmt.Println(colorReset, "")
}

func main() {
	cfg, cerr := config.NewConfig(".env")
	checkError(cerr)

	splash(cfg)

	pgpool, perr := postgres.NewPool(&cfg.Postgres)
	checkError(perr)

	db, gerr := gormConn.NewPostgresGormDB(pgpool)
	checkError(gerr)

	redisPool := buildRedisPool(cfg)

	grpcServer := createGrpcServer(cfg)

	grpcConn := server.InitGRPCConn(fmt.Sprintf("127.0.0.1:%v", cfg.Port.GRPC), false, "")

	registerGrpcHandlers(grpcServer.Server, *cfg, db, redisPool, grpcConn)

	// Reflection for Evans CLI for GRPC Debugging. DO NOT EXPOSE THE SERVER PORT!
	reflection.Register(grpcServer)

	restServer := createRestServer(cfg.Port.REST)
	registerRestHandlers(context.Background(), restServer.ServeMux, fmt.Sprintf(":%s", cfg.Port.GRPC), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Uncomment to enable pub sub
	// psClient := createPubSubClient(cfg.Google.ProjectID, cfg.Google.ServiceAccountFile)
	// psHandlers := registerPubSubHandlers(context.Background(), db, *cfg)
	//
	// _ = psClient.StartSubscriptions(psHandlers...)

	healthcheck.RegisterHealthHandler(grpcServer.Server)

	_ = grpcServer.Run()
	_ = restServer.Run()
	_ = grpcServer.AwaitTermination()
}

//nolint // createPubSubClient creates a pubsub client
func createPubSubClient(projectID, googleSaFile string) *pubsubSDK.PubSub {
	return pubsubSDK.NewPubSub(projectID, &googleSaFile)
}

// createGrpcServer creates a grpc server
func createGrpcServer(cfg *config.Config) *server.Grpc {
	if cfg.Env == envDevelopment {
		return server.NewDevelopmentGrpc(cfg.Port.GRPC)
	}
	srv, err := server.NewProductionGrpc(cfg.Env, cfg.ServiceName, cfg.Google.ProjectID, cfg.Port.GRPC)
	checkError(err)
	return srv
}

// createRestServer creates a rest server
func createRestServer(port string) *server.Rest {
	return server.NewProductionRest(port)
}

// registerGrpcHandlers registers all the grpc handlers
func registerGrpcHandlers(
	server *grpc.Server,
	cfg config.Config,
	db *gorm.DB,
	redisPool *redis.Pool,
	grpcConn *grpc.ClientConn,
) {
	// start register all module's gRPC handlers
	userModules.InitGrpc(server, cfg, db, redisPool, grpcConn)
	// end of register all module's gRPC handlers
}

// registerRestHandlers registers all the rest handlers
func registerRestHandlers(ctx context.Context, server *runtime.ServeMux, grpcPort string, options ...grpc.DialOption) {
	// start register all module's REST handlers
	options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize)))
	userModules.InitRest(ctx, server, grpcPort, options...)
	// end of register all module's REST handlers
}

//nolint // registerPubSubHandlers registers all the pubsub handlers
func registerPubSubHandlers(
	ctx context.Context,
	db *gorm.DB,
	config config.Config,
) []pubsubSDK.Subscriber {
	var handlers []pubsubSDK.Subscriber
	handlers = append(handlers, notificationModules.InitSendEmailSubscription(ctx, db, config))

	return handlers
}

// buildRedisPool builds a redis pool
func buildRedisPool(cfg *config.Config) *redis.Pool {
	cachePool := commonRedis.NewPool(cfg.Redis.Address, cfg.Redis.Password)

	ctx := context.Background()
	_, err := cachePool.GetContext(ctx)

	if err != nil {
		checkError(err)
	}

	log.Print("redis successfully connected!")
	return cachePool
}

// checkError checks if the error is not nil and prints it out
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
