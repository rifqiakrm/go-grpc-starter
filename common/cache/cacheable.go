// Package cache is a generated cache package.
package cache

import (
	"time"
)

// Cacheable is a mock of Cache interface.
type Cacheable interface {
	// Ping ping the redis server
	Ping() error
	// Get get data from redis by cache key
	Get(key string) ([]byte, error)
	// Set set data with defined cache key, value, and time-to-live (ttl) to store in redis
	Set(key string, value interface{}, ttl int) error
	// SetWithExpireAt set key value and update expire using unix timestamp
	SetWithExpireAt(key string, value interface{}, ttl time.Time) error
	// Exists check if key is exist in redis
	Exists(key string) (bool, error)
	// Remove remove cache by cache key
	Remove(key string) error
	// BulkRemove remove cache by certain cache key pattern
	BulkRemove(pattern string) error
	// Scan scan all cache key with certain pattern
	Scan(pattern string) ([]string, error)
}
