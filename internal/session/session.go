// Package session provides Redis-backed admin session storage.
// Sessions survive server restarts because they live in Redis, not in
// application memory or the SQLite admin row's session_token column.
package session

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "admin_session:"

func key(token string) string {
	return keyPrefix + token
}

// Set stores adminID under the session token with the given TTL.
// Call this immediately after a successful admin login.
func Set(rdb *redis.Client, token string, adminID uint, ttl time.Duration) error {
	return rdb.Set(context.Background(), key(token), strconv.FormatUint(uint64(adminID), 10), ttl).Err()
}

// Get returns the adminID stored under token, or (0, nil) when the token is
// missing or expired. Any Redis transport error is returned as-is.
func Get(rdb *redis.Client, token string) (uint, error) {
	val, err := rdb.Get(context.Background(), key(token)).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("session get: %w", err)
	}

	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("session corrupt value: %w", err)
	}
	return uint(id), nil
}

// Delete removes the session token from Redis (logout / expiry cleanup).
func Delete(rdb *redis.Client, token string) error {
	return rdb.Del(context.Background(), key(token)).Err()
}
