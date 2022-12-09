package healthcheck_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"grpc-starter/common/healthcheck"
	mock_healthcheck "grpc-starter/test/mock/common/healthcheck"
)

var (
	testHealthCheckRequest = &grpc_health_v1.HealthCheckRequest{Service: "firestore"}
)

type HealthHandlerExecutor struct {
	handler *healthcheck.HealthHandler
	checker *mock_healthcheck.MockCheckHealth
}

func TestNewHealthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of HealthHandler", func(t *testing.T) {
		exec := createHealthHandlerExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestHealthHandler_Check(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createHealthHandlerExecutor(ctrl)

		resp, err := exec.handler.Check(testContext, nil)

		assert.NotNil(t, err)
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_UNKNOWN, resp.GetStatus())
	})

	t.Run("system is not healthy", func(t *testing.T) {
		exec := createHealthHandlerExecutor(ctrl)
		exec.checker.EXPECT().Check(testContext).Return(status.New(codes.Internal, "dep1 is not alive").Err())

		resp, err := exec.handler.Check(testContext, testHealthCheckRequest)

		assert.NotNil(t, err)
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
	})

	t.Run("system is healthy", func(t *testing.T) {
		exec := createHealthHandlerExecutor(ctrl)
		exec.checker.EXPECT().Check(testContext).Return(nil)

		resp, err := exec.handler.Check(testContext, testHealthCheckRequest)

		assert.Nil(t, err)
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.GetStatus())
	})
}

func createHealthHandlerExecutor(ctrl *gomock.Controller) *HealthHandlerExecutor {
	c := mock_healthcheck.NewMockCheckHealth(ctrl)
	h := healthcheck.NewHealthHandler(c)
	return &HealthHandlerExecutor{
		handler: h,
		checker: c,
	}
}
