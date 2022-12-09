package healthcheck

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// HealthHandler handles HTTP/2 gRPC request for health checking.
type HealthHandler struct {
	grpc_health_v1.UnimplementedHealthServer
	checker CheckHealth
}

// NewHealthHandler creates an instance of HealthHandler.
func NewHealthHandler(checker CheckHealth) *HealthHandler {
	return &HealthHandler{checker: checker}
}

// Check checks the entire system health, including its dependecies.
func (hc *HealthHandler) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if request == nil {
		st := status.New(codes.InvalidArgument, "health check request is nil")
		return createHealthCheckResponse(grpc_health_v1.HealthCheckResponse_UNKNOWN), st.Err()
	}

	if err := hc.checker.Check(ctx); err != nil {
		return createHealthCheckResponse(grpc_health_v1.HealthCheckResponse_NOT_SERVING), err
	}
	return createHealthCheckResponse(grpc_health_v1.HealthCheckResponse_SERVING), nil
}

func createHealthCheckResponse(status grpc_health_v1.HealthCheckResponse_ServingStatus) *grpc_health_v1.HealthCheckResponse {
	return &grpc_health_v1.HealthCheckResponse{
		Status: status,
	}
}
