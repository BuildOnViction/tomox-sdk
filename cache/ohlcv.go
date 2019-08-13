package cache

import (
	"github.com/go-redis/redis"
)

const ohlcvKey = "ohlcv"

// OHLCVCacheClient redis client
type OHLCVCacheClient struct {
	client *redis.Client
}

// NewOHLCVCacheClient init market cache client
func NewOHLCVCacheClient(addr string, password string, db int) *OHLCVCacheClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &OHLCVCacheClient{client: client}
}
