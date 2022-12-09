package healthcheck

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// RegisterHealthHandler initializes gRPC health check modules.
func RegisterHealthHandler(server *grpc.Server) {
	checker := NewHealthChecker()
	health := NewHealthHandler(checker)
	grpc_health_v1.RegisterHealthServer(server, health)
}
