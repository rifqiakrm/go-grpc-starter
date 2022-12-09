package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logsettable "github.com/grpc-ecosystem/go-grpc-middleware/logging/settable"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonJwt "grpc-starter/common/jwt"
	"grpc-starter/server/interceptor"
)

const (
	// connProtocol is the protocol used to listen for incoming connections.
	connProtocol = "tcp"
	// envProduction is the environment name for production environment.
	envProduction = "production"
	// tracingSampleProbability is the sample probability for tracing.
	tracingSampleProbability = 0.01
)

const (
	// maxMsgSize is the maximum message size allowed.
	maxMsgSize = 1024 * 1024 * 150
)

// Grpc is responsible to act as gRPC server.
// It composes grpc.Server.
type Grpc struct {
	*grpc.Server
	listener net.Listener
	port     string
}

// NewGrpc creates an instance of Grpc.
func NewGrpc(port string, options ...grpc.ServerOption) *Grpc {
	options = append(options, grpc.MaxSendMsgSize(maxMsgSize))
	options = append(options, grpc.MaxRecvMsgSize(maxMsgSize))
	srv := grpc.NewServer(options...)
	return &Grpc{
		Server: srv,
		port:   port,
	}
}

// NewDevelopmentGrpc creates an instance of Grpc for used in development environment.
//
// These are list of interceptors that are attached (from innermost to outermost):
// 	- Metrics, using Prometheus.
// 	- Logging, using zap logger.
// 	- Recoverer, using grpc_recovery.
func NewDevelopmentGrpc(port string) *Grpc {
	options := grpc_middleware.WithUnaryServerChain(defaultUnaryServerInterceptors()...)
	srv := NewGrpc(port, options)
	grpc_prometheus.Register(srv.Server)
	return srv
}

// NewProductionGrpc creates an instance of Grpc with default production options attached.
// Actually, it can be used for non-production environment (such as staging or sandbox) as long as the environment satisfies all prerequisites.
//
// These are list of interceptors that are attached (from innermost to outermost):
// 	- Metrics, using Prometheus.
// 	- Logging, using zap logger.
// 	- Recoverer, using grpc_recovery.
// 	- Error Reporter, using Google Cloud Error Reporter.
//
// It also activates some auxiliaries:
// 	- Profiler, using Google Cloud Profiler.
// 	- Tracing, using Google Cloud Stackdriver Trace. The sample probability is 1% for production environment. Otherwise, it is 100%.
func NewProductionGrpc(env, serviceName, gcpProjectID, grpcPort string) (*Grpc, error) {
	if err := activateProfiling(gcpProjectID, serviceName); err != nil {
		return nil, err
	}
	if err := activateTracing(gcpProjectID, env); err != nil {
		return nil, err
	}

	reporter, err := createErrorReporter(serviceName, gcpProjectID)
	if err != nil {
		return nil, err
	}

	midds := []grpc.UnaryServerInterceptor{interceptor.ErrorReporting(reporter)}
	midds = append(midds, defaultUnaryServerInterceptors()...)

	options := grpc_middleware.WithUnaryServerChain(midds...)
	srv := NewGrpc(grpcPort, grpc.StatsHandler(&ocgrpc.ServerHandler{}), options)
	grpc_prometheus.Register(srv.Server)

	return srv, nil
}

// Run runs the server.
// It basically runs grpc.Server.Serve and is a blocking.
func (g *Grpc) Run() error {
	var err error
	g.listener, err = net.Listen(connProtocol, fmt.Sprintf(":%s", g.port))
	if err != nil {
		return err
	}

	go g.serve()
	log.Printf("grpc server is running on port %s\n", g.port)
	return nil
}

// AwaitTermination blocks the server and wait for termination signal.
// The termination signal must be one of SIGINT or SIGTERM.
// Once it receives one of those signals, the gRPC server will perform graceful stop and close the listener.
func (g *Grpc) AwaitTermination() error {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	g.GracefulStop()
	return g.listener.Close()
}

func (g *Grpc) serve() {
	if err := g.Serve(g.listener); err != nil {
		panic(err)
	}
}

// createErrorReporter creates an instance of Google Cloud Error Reporter.
func createErrorReporter(serviceName, gcpProjectID string) (*errorreporting.Client, error) {
	return errorreporting.NewClient(context.Background(), gcpProjectID, errorreporting.Config{
		ServiceName: serviceName,
	})
}

// defaultUnaryServerInterceptors returns a list of default unary server interceptors.
func defaultUnaryServerInterceptors() []grpc.UnaryServerInterceptor {
	logger, _ := zap.NewProduction() // error is impossible, hence ignored.
	grpc_zap.SetGrpcLoggerV2(grpc_logsettable.ReplaceGrpcLoggerV2(), logger)
	grpc_prometheus.EnableHandlingTimeHistogram()

	options := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryHandler)),
		grpc_zap.UnaryServerInterceptor(logger),
		grpc_auth.UnaryServerInterceptor(commonJwt.Authorize),
		grpc_prometheus.UnaryServerInterceptor,
	}
	return options
}

// recoveryHandler is a recovery handler for grpc_recovery.
func recoveryHandler(p interface{}) error {
	return status.Errorf(codes.Unknown, "%v", p)
}

// activateProfiling activates Google Cloud Profiler.
func activateProfiling(projectID, serviceName string) error {
	cfg := profiler.Config{
		Service:           serviceName,
		ProjectID:         projectID,
		EnableOCTelemetry: true,
	}
	return profiler.Start(cfg)
}

// activateTracing activates Google Cloud Stackdriver Trace.
func activateTracing(projectID, env string) error {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		return err
	}

	trace.RegisterExporter(exporter)
	if env == envProduction {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(tracingSampleProbability)})
	} else {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	return view.Register(ocgrpc.DefaultServerViews...)
}
