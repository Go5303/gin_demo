package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var prefix string

// Init sets the Redis client and key prefix for the cache package.
func Init(client *redis.Client, keyPrefix string) {
	rdb = client
	prefix = keyPrefix
}

func buildKey(key string) string {
	if prefix != "" {
		return prefix + ":" + key
	}
	return key
}

func ctx() context.Context {
	return context.Background()
}

// Get returns the string value of key, or empty string if not found.
func Get(key string) (string, error) {
	return rdb.Get(ctx(), buildKey(key)).Result()
}

// Set stores key with value and expiration. Pass 0 for no expiration.
func Set(key string, value any, expiration time.Duration) error {
	return rdb.Set(ctx(), buildKey(key), value, expiration).Err()
}

// SetNX sets key only if it does not exist. Returns true if the key was set.
func SetNX(key string, value any, expiration time.Duration) (bool, error) {
	return rdb.SetNX(ctx(), buildKey(key), value, expiration).Result()
}

// Del deletes one or more keys.
func Del(keys ...string) error {
	fullKeys := make([]string, len(keys))
	for i, k := range keys {
		fullKeys[i] = buildKey(k)
	}
	return rdb.Del(ctx(), fullKeys...).Err()
}

// Incr increments the integer value of key by 1.
func Incr(key string) (int64, error) {
	return rdb.Incr(ctx(), buildKey(key)).Result()
}

// IncrBy increments the integer value of key by n.
func IncrBy(key string, n int64) (int64, error) {
	return rdb.IncrBy(ctx(), buildKey(key), n).Result()
}

// Decr decrements the integer value of key by 1.
func Decr(key string) (int64, error) {
	return rdb.Decr(ctx(), buildKey(key)).Result()
}

// DecrBy decrements the integer value of key by n.
func DecrBy(key string, n int64) (int64, error) {
	return rdb.DecrBy(ctx(), buildKey(key), n).Result()
}

// Expire sets a timeout on key.
func Expire(key string, expiration time.Duration) error {
	return rdb.Expire(ctx(), buildKey(key), expiration).Err()
}

// TTL returns the remaining time to live of a key.
func TTL(key string) (time.Duration, error) {
	return rdb.TTL(ctx(), buildKey(key)).Result()
}

// Exists returns true if the key exists.
func Exists(key string) (bool, error) {
	n, err := rdb.Exists(ctx(), buildKey(key)).Result()
	return n > 0, err
}
