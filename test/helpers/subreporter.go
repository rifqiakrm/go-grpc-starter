package helpers

import (
	"github.com/golang/mock/gomock"
)

// SubReporter implements the same interface required by gomock
type SubReporter struct {
	t []gomock.TestReporter
}

// NewSubReporter creates instance of SubReporter
func NewSubReporter(t gomock.TestReporter) *SubReporter {
	return &SubReporter{t: []gomock.TestReporter{t}}
}

// Add is responsible for add subtest
func (s *SubReporter) Add(t gomock.TestReporter) func() {
	s.t = append(s.t, t)
	return func() {
		s.t = s.t[:len(s.t)-1]
	}
}

// Errorf returns error of SubReporter
func (s *SubReporter) Errorf(format string, args ...interface{}) {
	s.t[len(s.t)-1].Errorf(format, args...)
}

// Fatalf returns fatal of SubReporter
func (s *SubReporter) Fatalf(format string, args ...interface{}) {
	s.t[len(s.t)-1].Fatalf(format, args...)
}
