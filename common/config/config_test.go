package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"grpc-starter/common/config"
)

func TestNewConfig(t *testing.T) {
	t.Run("fail to create an instance of Config due to invalid env file", func(t *testing.T) {
		cfg, err := config.NewConfig("../../test/fixture/env.incomplete")
		assert.NotNil(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("fail to create an instance of Config due to file not found", func(t *testing.T) {
		cfg, err := config.NewConfig("somewhere/will/not/be/found")
		assert.NotNil(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("successfully read config", func(t *testing.T) {
		cfg, err := config.NewConfig("../../test/fixture/env.valid")
		assert.Nil(t, err)
		assert.NotNil(t, cfg)
	})
}
