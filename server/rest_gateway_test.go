package server_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"grpc-starter/server"
)

var (
	testRestPort = "8081"
)

func TestNewRest(t *testing.T) {
	t.Run("success create rest server", func(t *testing.T) {
		srv := server.NewRest(testRestPort)
		assert.NotNil(t, srv)
	})
}

func TestNewProductionRest(t *testing.T) {
	t.Run("success create production rest server", func(t *testing.T) {
		srv := server.NewRest(testRestPort)
		assert.NotNil(t, srv)
	})
}

func TestRest_EnablePrometheus(t *testing.T) {
	t.Run("success enable prometheus", func(t *testing.T) {
		srv := server.NewRest(testRestPort)
		err := srv.EnablePrometheus()
		assert.Nil(t, err)
	})
}

func TestRest_EnableHealth(t *testing.T) {
	t.Run("success enable health check", func(t *testing.T) {
		srv := server.NewRest(testRestPort)
		err := srv.EnableHealth()
		assert.Nil(t, err)
	})
}
