package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/tomochain/tomox-sdk/types"
)

const fiatKey = "fiat:"

// FiatCacheClient redis client
type FiatCacheClient struct {
	client *redis.Client
}

// NewFiatCacheClient init market cache client
func NewFiatCacheClient(addr string, password string, db int) *FiatCacheClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &FiatCacheClient{client: client}
}

// AddFiat add fiat of symbol to cache
func (m *FiatCacheClient) AddFiat(symbol string, fiatItem *types.FiatPriceItem) error {
	data, err := json.Marshal(fiatItem)
	if err != nil {
		return err
	}
	now := time.Now().Unix() * 1000
	diff := now - fiatItem.Timestamp
	err = m.client.ZAdd(fiatKey+symbol, &redis.Z{
		Score:  float64(diff),
		Member: data,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetCurrentFiat get fiat of specific symbol
func (m *FiatCacheClient) GetCurrentFiat(symbol string) (*types.FiatPriceItem, error) {
	val, err := m.client.ZRangeByScore(fiatKey+symbol, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		return nil, err
	}

	var data types.FiatPriceItem
	err = json.Unmarshal([]byte(val[0]), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetFiatRange get fiat array in hours
func (m *FiatCacheClient) GetFiatRange(symbol string, hour float64) ([]*types.FiatPriceItem, error) {
	begin := hour * 60 * 60 * 1000
	val, err := m.client.ZRangeByScore(fiatKey+symbol, &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%f", begin),
	}).Result()
	if err != nil {
		return nil, err
	}

	var res []*types.FiatPriceItem
	for _, r := range val {
		fmt.Print(r)
		d := &types.FiatPriceItem{}
		err = json.Unmarshal([]byte(r), d)
		if err == nil {
			res = append(res, d)
		}

	}

	return res, nil
}
