package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"grpc-starter/common/config"
	"grpc-starter/common/postgres"
)

func TestNewPool(t *testing.T) {
	cfg := &config.Postgres{
		Host:            "localhost",
		Port:            "5432",
		Name:            "postgres",
		User:            "user",
		Password:        "password",
		MaxOpenConns:    "10",
		MaxConnLifetime: "10m",
		MaxIdleLifetime: "5m",
	}

	t.Run("fail build sql client", func(t *testing.T) {
		client, err := postgres.NewPool(cfg)

		assert.NotNil(t, err)
		assert.Nil(t, client)
	})
}
