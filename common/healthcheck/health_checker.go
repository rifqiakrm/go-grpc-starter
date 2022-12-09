package healthcheck

import (
	"context"
)

// CheckHealth is the interface that defines the health check.
type CheckHealth interface {
	// Check checks the health of the system, including its dependencies.
	Check(ctx context.Context) error
}

// HealthChecker is responsible for doing the health check.
type HealthChecker struct {
}

// NewHealthChecker creates an instance of HealthChecker.
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

// Check checks the health of the system.
// It doesn't check the dependencies since the project already handles if the dependencies fail.
// It always return nil.
func (hc *HealthChecker) Check(_ context.Context) error {
	return nil
}
