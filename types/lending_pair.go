package types

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/globalsign/mgo/bson"
)

// LendingPair struct is used to model the lendingPair data in the system and DB
type LendingPair struct {
	ID                   bson.ObjectId  `json:"-" bson:"_id"`
	Term                 uint64         `json:"term,omitempty" bson:"term"`
	LendingTokenSymbol   string         `json:"lendingTokenSymbol,omitempty" bson:"lendingTokenSymbol"`
	LendingTokenAddress  common.Address `json:"lendingTokenAddress,omitempty" bson:"lendingTokenAddress"`
	LendingTokenDecimals int            `json:"lendingTokenDecimals,omitempty" bson:"lendingTokenDecimals"`
	CreatedAt            time.Time      `json:"-" bson:"createdAt"`
	UpdatedAt            time.Time      `json:"-" bson:"updatedAt"`
}

// LendingPairRecord struct for database
type LendingPairRecord struct {
	ID                   bson.ObjectId `json:"id" bson:"_id"`
	Term                 string        `json:"term" bson:"term"`
	LendingTokenSymbol   string        `json:"lendingTokenSymbol" bson:"lendingTokenSymbol"`
	LendingTokenAddress  string        `json:"lendingTokenAddress" bson:"lendingTokenAddress"`
	LendingTokenDecimals int           `json:"lendingTokenDecimals" bson:"lendingTokenDecimals"`
	CreatedAt            time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// UnmarshalJSON umarshal JSON
func (p *LendingPair) UnmarshalJSON(b []byte) error {
	lendingPair := map[string]interface{}{}

	err := json.Unmarshal(b, &lendingPair)
	if err != nil {
		return err
	}

	if lendingPair["term"] != nil {
		p.Term, _ = strconv.ParseUint(lendingPair["term"].(string), 10, 64)
	}

	if lendingPair["lendingTokenAddress"] != nil {
		p.LendingTokenAddress = common.HexToAddress(lendingPair["lendingTokenAddress"].(string))
	}

	if lendingPair["lendingTokenSymbol"] != nil {
		p.LendingTokenSymbol = lendingPair["lendingTokenSymbol"].(string)
	}

	if lendingPair["lendingTokenDecimals"] != nil {
		p.LendingTokenDecimals = lendingPair["lendingTokenDecimals"].(int)
	}
	return nil
}

// MarshalJSON marshal json byte
func (p *LendingPair) MarshalJSON() ([]byte, error) {
	lendingPair := map[string]interface{}{
		"term":                 strconv.FormatUint(p.Term, 10),
		"lendingTokenAddress":  p.LendingTokenAddress,
		"lendingTokenSymbol":   p.LendingTokenSymbol,
		"lendingTokenDecimals": p.LendingTokenDecimals,
	}

	return json.Marshal(lendingPair)
}

// SetBSON get lending pair object from database
func (p *LendingPair) SetBSON(raw bson.Raw) error {
	decoded := &LendingPairRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	p.ID = decoded.ID
	p.Term, _ = strconv.ParseUint(decoded.Term, 10, 64)
	p.LendingTokenSymbol = decoded.LendingTokenSymbol
	p.LendingTokenAddress = common.HexToAddress(decoded.LendingTokenAddress)
	p.LendingTokenDecimals = decoded.LendingTokenDecimals
	p.LendingTokenSymbol = decoded.LendingTokenSymbol
	p.CreatedAt = decoded.CreatedAt
	p.UpdatedAt = decoded.UpdatedAt
	return nil
}

// GetBSON insert record to database
func (p *LendingPair) GetBSON() (interface{}, error) {
	return &LendingPairRecord{
		ID:                   p.ID,
		Term:                 strconv.FormatUint(p.Term, 10),
		LendingTokenAddress:  p.LendingTokenAddress.Hex(),
		LendingTokenDecimals: p.LendingTokenDecimals,
		LendingTokenSymbol:   p.LendingTokenSymbol,
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
	}, nil
}

// Name name of lending pair
func (p *LendingPair) Name() string {
	name := strconv.FormatUint(p.Term, 10) + "/" + p.LendingTokenSymbol
	return name
}
