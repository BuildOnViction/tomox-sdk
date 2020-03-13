package types

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/utils/math"
)

// LendingTick is the format in which mongo aggregate pipeline returns data when queried for OHLCV data
type LendingTick struct {
	LendingID LendingID `json:"lendingID,omitempty" bson:"lendingID"`
	Open      uint64    `json:"open,omitempty" bson:"open"`
	Close     uint64    `json:"close,omitempty" bson:"close"`
	High      uint64    `json:"high,omitempty" bson:"high"`
	Low       uint64    `json:"low,omitempty" bson:"low"`
	Volume    *big.Int  `json:"volume,omitempty" bson:"volume"`
	Count     *big.Int  `json:"count,omitempty" bson:"count"`
	Timestamp int64     `json:"timestamp,omitempty" bson:"timestamp"`
	Duration  int64     `json:"duration" bson:"duration"`
	Unit      string    `json:"unit" bson:"unit"`
}

// LendingID is the subdocument for aggregate grouping for OHLCV data
type LendingID struct {
	Name         string         `json:"name" bson:"name"`
	Term         uint64         `json:"term" bson:"term"`
	LendingToken common.Address `json:"lendingToken" bson:"lendingToken"`
}

// LendingTicks array of lending ticks
type LendingTicks []*LendingTick

// MarshalJSON returns the json encoded byte array representing the trade struct
func (t *LendingTick) MarshalJSON() ([]byte, error) {
	tick := map[string]interface{}{
		"lendingID": map[string]interface{}{
			"name":         t.LendingID.Name,
			"term":         strconv.FormatUint(t.LendingID.Term, 10),
			"lendingToken": t.LendingID.LendingToken.Hex(),
		},
		"timestamp": t.Timestamp,
	}

	tick["open"] = strconv.FormatUint(t.Open, 10)
	tick["high"] = strconv.FormatUint(t.High, 10)
	tick["low"] = strconv.FormatUint(t.Low, 10)
	tick["close"] = strconv.FormatUint(t.Close, 10)
	if t.Volume != nil {
		tick["volume"] = t.Volume.String()
	}

	if t.Count != nil {
		tick["count"] = t.Count.String()
	}
	tick["duration"] = t.Duration
	tick["unit"] = t.Unit

	bytes, err := json.Marshal(tick)
	return bytes, err
}

// UnmarshalJSON creates a trade object from a json byte string
func (t *LendingTick) UnmarshalJSON(b []byte) error {
	tick := map[string]interface{}{}
	err := json.Unmarshal(b, &tick)

	if err != nil {
		return err
	}

	if tick["lendingID"] != nil {
		lendingID := tick["lendingID"].(map[string]interface{})
		t.LendingID = LendingID{
			Name:         lendingID["name"].(string),
			LendingToken: common.HexToAddress(lendingID["lendingToken"].(string)),
		}
		t.LendingID.Term, _ = strconv.ParseUint(lendingID["term"].(string), 10, 64)
	}

	if tick["timestamp"] != nil {
		t.Timestamp = int64(tick["timestamp"].(float64))
	}
	t.Open, _ = strconv.ParseUint(tick["open"].(string), 10, 64)
	t.High, _ = strconv.ParseUint(tick["high"].(string), 10, 64)
	t.Low, _ = strconv.ParseUint(tick["low"].(string), 10, 64)
	t.Close, _ = strconv.ParseUint(tick["close"].(string), 10, 64)

	if tick["volume"] != nil {
		t.Volume = math.ToBigInt(tick["volume"].(string))
	}

	if tick["count"] != nil {
		t.Count = math.ToBigInt(tick["count"].(string))
	}
	if tick["unit"] != nil {
		t.Unit = tick["unit"].(string)
	}
	if tick["duration"] != nil {
		t.Duration = int64(tick["duration"].(float64))
	}
	return nil
}

// AddressCode generate code from pair
func (t *LendingTick) AddressCode() string {
	code := strconv.FormatUint(t.LendingID.Term, 10) + "::" + t.LendingID.LendingToken.Hex()
	return code
}
