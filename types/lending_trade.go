package types

import (
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/tomochain/tomox-sdk/utils/math"

	"github.com/globalsign/mgo/bson"
)

const (
	LendingTradeStatusPending = "PENDING"
	LendingTradeStatusSuccess = "SUCCESS"
	LengdingTradeStatusError  = "ERROR"
)

// LendingTrade struct holds arguments corresponding to a "Taker Order"
// To be valid an accept by the matching engine (and ultimately the exchange smart-contract),
// the trade signature must be made from the trader Maker account
type LendingTrade struct {
	ID               bson.ObjectId  `json:"id,omitempty" bson:"_id"`
	BorrowingOwner   common.Address `bson:"borrowingOwner" json:"borrowingOwner"`
	InvestingOwner   common.Address `bson:"investingOwner" json:"investingOwner"`
	LendingToken     common.Address `bson:"lendingToken" json:"lendingToken"`
	CollateralToken  common.Address `bson:"collateralToken" json:"collateralToken"`
	BorrowingHash    common.Hash    `bson:"borrowingHash" json:"borrowingHash"`
	InvestingHash    common.Hash    `bson:"investingHash" json:"investingHash"`
	BorrowingRelayer common.Address `bson:"borrowingRelayer" json:"borrowingRelayer"`
	InvestingRelayer common.Address `bson:"investingRelayer" json:"investingRelayer"`
	Term             uint64         `bson:"term" json:"term"`
	Interest         uint64         `bson:"interest" json:"interest"`
	CollateralPrice  *big.Int       `bson:"collateralPrice" json:"collateralPrice"`
	LiquidationPrice *big.Int       `bson:"liquidationPrice" json:"liquidationPrice"`
	Amount           *big.Int       `bson:"amount" json:"amount"`
	BorrowingFee     *big.Int       `bson:"borrowingFee" json:"borrowingFee"`
	InvestingFee     *big.Int       `bson:"investingFee" json:"investingFee"`
	Status           string         `bson:"status" json:"status"`
	TakerOrderSide   string         `bson:"takerOrderSide" json:"takerOrderSide"`
	TakerOrderType   string         `bson:"takerOrderType" json:"takerOrderType"`
	Hash             common.Hash    `bson:"hash" json:"hash"`
	TxHash           common.Hash    `bson:"txHash" json:"txHash"`
	ExtraData        string         `bson:"extraData" json:"extraData"`
	CreatedAt        time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time      `bson:"updatedAt" json:"updatedAt"`
}

// LendingTradeSpec for query
type LendingTradeSpec struct {
	CollateralToken string
	LendingToken    string
	DateFrom        int64
	DateTo          int64
}

// LendingTradeRes response api
type LendingTradeRes struct {
	Total         int             `json:"total" bson:"total"`
	LendingTrades []*LendingTrade `json:"trades" bson:"orders"`
}

// LendingTradeRecord struct item database
type LendingTradeRecord struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	BorrowingOwner   string        `bson:"borrowingOwner" json:"borrowingOwner"`
	InvestingOwner   string        `bson:"investingOwner" json:"investingOwner"`
	LendingToken     string        `bson:"lendingToken" json:"lendingToken"`
	CollateralToken  string        `bson:"collateralToken" json:"collateralToken"`
	BorrowingHash    string        `bson:"borrowingHash" json:"borrowingHash"`
	InvestingHash    string        `bson:"investingHash" json:"investingHash"`
	BorrowingRelayer string        `bson:"borrowingRelayer" json:"borrowingRelayer"`
	InvestingRelayer string        `bson:"investingRelayer" json:"investingRelayer"`
	Term             string        `bson:"term" json:"term"`
	Interest         string        `bson:"interest" json:"interest"`
	CollateralPrice  string        `bson:"collateralPrice" json:"collateralPrice"`
	LiquidationPrice string        `bson:"liquidationPrice" json:"liquidationPrice"`
	Amount           string        `bson:"amount" json:"amount"`
	BorrowingFee     string        `bson:"borrowingFee" json:"borrowingFee"`
	InvestingFee     string        `bson:"investingFee" json:"investingFee"`
	Status           string        `bson:"status" json:"status"`
	TakerOrderSide   string        `bson:"takerOrderSide" json:"takerOrderSide"`
	TakerOrderType   string        `bson:"takerOrderType" json:"takerOrderType"`
	Hash             string        `bson:"hash" json:"hash"`
	TxHash           string        `bson:"txHash" json:"txHash"`
	ExtraData        string        `bson:"extraData" json:"extraData"`
	CreatedAt        time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time     `bson:"updatedAt" json:"updatedAt"`
}

// MarshalJSON returns the json encoded byte array representing the trade struct
func (t *LendingTrade) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// UnmarshalJSON creates a trade object from a json byte string
func (t *LendingTrade) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, t)
	if err != nil {
		return err
	}
	return nil
}

// GetBSON insert to mongodb
func (t *LendingTrade) GetBSON() (interface{}, error) {
	tr := LendingTradeRecord{
		ID:              t.ID,
		BorrowingOwner:  t.BorrowingOwner.Hex(),
		InvestingOwner:  t.InvestingOwner.Hex(),
		CollateralToken: t.CollateralToken.Hex(),
		LendingToken:    t.LendingToken.Hex(),
		BorrowingHash:   t.BorrowingHash.Hex(),
		Hash:            t.Hash.Hex(),
		TxHash:          t.TxHash.Hex(),
		InvestingHash:   t.InvestingHash.Hex(),
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
		Interest:        strconv.FormatUint(t.Interest, 10),
		Status:          t.Status,
		Amount:          t.Amount.String(),
		BorrowingFee:    t.BorrowingFee.String(),
		InvestingFee:    t.InvestingFee.String(),
		TakerOrderSide:  t.TakerOrderSide,
	}

	return tr, nil
}

// SetBSON get data from database
func (t *LendingTrade) SetBSON(raw bson.Raw) error {
	decoded := new(LendingTradeRecord)

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	t.ID = decoded.ID

	t.BorrowingOwner = common.HexToAddress(decoded.BorrowingOwner)
	t.InvestingOwner = common.HexToAddress(decoded.InvestingOwner)
	t.CollateralToken = common.HexToAddress(decoded.CollateralToken)
	t.LendingToken = common.HexToAddress(decoded.LendingToken)
	t.BorrowingHash = common.HexToHash(decoded.BorrowingHash)
	t.InvestingHash = common.HexToHash(decoded.InvestingHash)
	t.BorrowingRelayer = common.HexToAddress(decoded.BorrowingRelayer)
	t.InvestingRelayer = common.HexToAddress(decoded.InvestingRelayer)
	t.Hash = common.HexToHash(decoded.Hash)
	t.TxHash = common.HexToHash(decoded.TxHash)
	t.Status = decoded.Status
	t.Amount = math.ToBigInt(decoded.Amount)
	t.Interest, _ = strconv.ParseUint(decoded.Interest, 10, 64)
	t.Term, _ = strconv.ParseUint(decoded.Term, 10, 64)
	t.BorrowingFee = math.ToBigInt(decoded.BorrowingFee)
	t.InvestingFee = math.ToBigInt(decoded.InvestingFee)

	t.CreatedAt = decoded.CreatedAt
	t.UpdatedAt = decoded.UpdatedAt
	t.TakerOrderSide = decoded.TakerOrderSide
	return nil
}

// ComputeHash returns hashes the trade
// The OrderHash, Amount, Taker and TradeNonce attributes must be
// set before attempting to compute the trade hash
func (t *LendingTrade) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(t.BorrowingOwner.Bytes())
	sha.Write(t.InvestingOwner.Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// LengdingTradeChangeEvent event for changing mongo watch
type LendingTradeChangeEvent struct {
	ID                interface{}   `bson:"_id"`
	OperationType     string        `bson:"operationType"`
	FullDocument      *LendingTrade `bson:"fullDocument,omitempty"`
	Ns                evNamespace   `bson:"ns"`
	DocumentKey       M             `bson:"documentKey"`
	UpdateDescription *updateDesc   `bson:"updateDescription,omitempty"`
}
