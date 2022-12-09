package healthcheck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"grpc-starter/common/healthcheck"
)

var (
	testContext = context.Background()
)

type HealthCheckerExecutor struct {
	service *healthcheck.HealthChecker
}

func TestNewHealthChecker(t *testing.T) {
	t.Run("successfully create an instance of HealthChecker", func(t *testing.T) {
		exec := createHealthCheckerExecutor()
		assert.NotNil(t, exec.service)
	})
}

func TestHealthChecker_Check(t *testing.T) {
	t.Run("all systems are well", func(t *testing.T) {
		exec := createHealthCheckerExecutor()

		err := exec.service.Check(testContext)

		assert.Nil(t, err)
	})
}

func createHealthCheckerExecutor() *HealthCheckerExecutor {
	u := healthcheck.NewHealthChecker()
	return &HealthCheckerExecutor{
		service: u,
	}
}
