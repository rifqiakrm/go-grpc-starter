package redis_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

	"grpc-starter/common/config"
	"grpc-starter/common/redis"
)

func TestNewPool(t *testing.T) {
	t.Run("fail create redis pool", func(t *testing.T) {
		server, _ := miniredis.Run()

		cfg := &config.Redis{
			Address: server.Addr(),
		}

		server.Close()
		pool := redis.NewPool(cfg.Address, "")

		ctx := context.Background()
		_, err := pool.GetContext(ctx)

		assert.NotNil(t, err)
	})

	t.Run("success create redis pool", func(t *testing.T) {
		server, _ := miniredis.Run()
		defer server.Close()

		cfg := &config.Redis{
			Address: server.Addr(),
		}

		pool := redis.NewPool(cfg.Address, "")
		ctx := context.Background()
		_, err := pool.GetContext(ctx)

		assert.Nil(t, err)
	})
}
