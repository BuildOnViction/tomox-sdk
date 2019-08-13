package cache

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/tomochain/tomox-sdk/types"
)

const pairDataKey = "pairdata"

// PairDataCacheClient redis client
type PairDataCacheClient struct {
	client *redis.Client
}

// NewPairDataCacheClient init market cache client
func NewPairDataCacheClient(addr string, password string, db int) *PairDataCacheClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &PairDataCacheClient{client: client}
}

// AddPairData executes the redis Set command
func (m *PairDataCacheClient) AddPairData(pairData types.PairData) error {
	data, err := json.Marshal(pairData)
	if err != nil {
		return err
	}
	err = m.client.Set(pairDataKey, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetPairData get cache pair data
func (m *PairDataCacheClient) GetPairData() (*types.PairData, error) {
	val, err := m.client.Get(pairDataKey).Result()
	if err != nil {
		return nil, err
	}
	var data types.PairData
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
