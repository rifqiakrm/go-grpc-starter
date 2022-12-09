package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	maxValLength   = 15
	cutVal         = 12
	maxIdle        = 3
	maxIdleTimeout = 240
)

// NewPool creates new redis server pool.
func NewPool(server string, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: maxIdleTimeout * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				log.Println("there was an error while dialing redis :", err)
				return nil, err
			}

			// Uncomment this if redis use password
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					_ = c.Close()
					log.Println("there was an error while doing auth redis :", err)
					return nil, err
				}
			}

			// Uncomment this if redis use database (0-15)
			// if _, err := c.Do("SELECT", os.Getenv("REDIS_DATABASE")); err != nil {
			//	c.Close()
			//	checkError(err)
			//	return nil, err
			// }

			return c, nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")

			if err != nil {
				log.Println("there was an error while pinging redis server :", err)
				return err
			}

			return nil
		},
	}
}

// Client define struct for redis client
type Client struct {
	cachePool *redis.Pool
}

// NewClient create new redis client
func NewClient(cachePool *redis.Pool) *Client {
	return &Client{
		cachePool: cachePool,
	}
}

// Ping ping the redis server
func (r *Client) Ping() error {
	conn := r.cachePool.Get()
	defer func() {
		_ = conn.Close()
	}()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

// Get get data from redis by cache key
func (r *Client) Get(key string) ([]byte, error) {
	conn := r.cachePool.Get()
	defer func() {
		_ = conn.Close()
	}()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

// Set set data with defined cache key, value, and time-to-live (ttl) to store in redis
func (r *Client) Set(key string, value interface{}, ttl int) error {
	conn := r.cachePool.Get()
	defer func() {
		_ = conn.Close()
	}()

	val, err := json.Marshal(value)

	if err != nil {
		return fmt.Errorf("error while marshalling interface for cache")
	}

	_, err = conn.Do("SET", key, val)
	if err != nil {
		v := string(val)
		if len(v) > maxValLength {
			v = v[0:cutVal] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	_, err = conn.Do("EXPIRE", key, ttl)

	if err != nil {
		return fmt.Errorf("error set expire key %s", key)
	}

	return err
}

// SetWithExpireAt set key value and update expire using unix timestamp
func (r *Client) SetWithExpireAt(key string, value interface{}, ttl time.Time) error {
	conn := r.cachePool.Get()
	defer func() {
		_ = conn.Close()
	}()

	val, err := json.Marshal(value)

	if err != nil {
		return fmt.Errorf("error while marshalling interface for cache")
	}

	_, err = conn.Do("SET", key, val)
	if err != nil {
		v := string(val)
		if len(v) > maxValLength {
			v = v[0:cutVal] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	_, err = conn.Do("EXPIREAT", key, ttl.Unix())

	if err != nil {
		return fmt.Errorf("error set expireat key %s", key)
	}

	return err
}

// Exists check if key is exist in redis
func (r *Client) Exists(key string) (bool, error) {
	conn := r.cachePool.Get()
	defer func() {
		_ = conn.Close()
	}()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

// Remove remove cache by cache key
func (r *Client) Remove(key string) error {
	conn := r.cachePool.Get()
	defer func() {
		_ = conn.Close()
	}()

	_, err := conn.Do("DEL", key)
	return err
}

// BulkRemove remove cache by certain cache key pattern
func (r *Client) BulkRemove(pattern string) error {
	conn := r.cachePool.Get()
	defer func() {
		_ = conn.Close()
	}()

	val, getKeysErr := redis.Strings(conn.Do("KEYS", pattern))

	if getKeysErr != nil {
		return fmt.Errorf("error retrieving '%s' keys", pattern)
	}

	for i := range val {
		_, err := conn.Do("DEL", val[i])
		if err != nil {
			return fmt.Errorf("error deleting '%s' keys", val[i])
		}
	}

	return nil
}

// Scan scan all cache key with certain pattern
func (r *Client) Scan(pattern string) ([]string, error) {
	conn := r.cachePool.Get()
	defer func() {
		_ = conn.Close()
	}()

	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}
