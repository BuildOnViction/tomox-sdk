package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/utils"

	"github.com/globalsign/mgo/bson"
)

const (
	TradeStatusOpen       = "OPEN"
	TradeStatusClosed     = "CLOSED"
	TradeStatusLiquidated = "LIQUIDATED"
)

// LendingTrade lending trade struct
type LendingTrade struct {
	Borrower               common.Address `bson:"borrower" json:"borrower"`
	Investor               common.Address `bson:"investor" json:"investor"`
	LendingToken           common.Address `bson:"lendingToken" json:"lendingToken"`
	CollateralToken        common.Address `bson:"collateralToken" json:"collateralToken"`
	BorrowingOrderHash     common.Hash    `bson:"borrowingOrderHash" json:"borrowingOrderHash"`
	InvestingOrderHash     common.Hash    `bson:"investingOrderHash" json:"investingOrderHash"`
	BorrowingRelayer       common.Address `bson:"borrowingRelayer" json:"borrowingRelayer"`
	InvestingRelayer       common.Address `bson:"investingRelayer" json:"investingRelayer"`
	Term                   uint64         `bson:"term" json:"term"`
	Interest               uint64         `bson:"interest" json:"interest"`
	CollateralPrice        *big.Int       `bson:"collateralPrice" json:"collateralPrice"`
	LiquidationPrice       *big.Int       `bson:"liquidationPrice" json:"liquidationPrice"`
	CollateralLockedAmount *big.Int       `bson:"collateralLockedAmount" json:"collateralLockedAmount"`
	LiquidationTime        uint64         `bson:"liquidationTime" json:"liquidationTime"`
	DepositRate            *big.Int       `bson:"depositRate" json:"depositRate"`
	Amount                 *big.Int       `bson:"amount" json:"amount"`
	BorrowingFee           *big.Int       `bson:"borrowingFee" json:"borrowingFee"`
	InvestingFee           *big.Int       `bson:"investingFee" json:"investingFee"`
	Status                 string         `bson:"status" json:"status"`
	TakerOrderSide         string         `bson:"takerOrderSide" json:"takerOrderSide"`
	TakerOrderType         string         `bson:"takerOrderType" json:"takerOrderType"`
	MakerOrderType         string         `bson:"makerOrderType" json:"makerOrderType"`
	TradeID                bson.ObjectId  `bson:"tradeID" json:"tradeID"`
	Hash                   common.Hash    `bson:"hash" json:"hash"`
	TxHash                 common.Hash    `bson:"txHash" json:"txHash"`
	ExtraData              string         `bson:"extraData" json:"extraData"`
	CreatedAt              time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt              time.Time      `bson:"updatedAt" json:"updatedAt"`
}

// MarshalJSON returns the json encoded byte array representing the trade struct
func (t *LendingTrade) MarshalJSON() ([]byte, error) {

	trade := map[string]interface{}{
		"borrower":           t.Borrower,
		"investor":           t.Investor,
		"borrowingOrderHash": t.BorrowingOrderHash,
		"investingOrderHash": t.InvestingOrderHash,
		"borrowingRelayer":   t.BorrowingRelayer,
		"investingRelayer":   t.InvestingRelayer,
		"term":               strconv.FormatUint(t.Term, 10),
		"interest":           strconv.FormatUint(t.Interest, 10),
		"collateralPrice":    t.CollateralPrice.String(),
		"liquidationPrice":   t.LiquidationPrice.String(),
		"liquidationTime":    strconv.FormatUint(t.LiquidationTime, 10),
		"depositRate":        t.DepositRate.String(),
		"amount":             t.Amount.String(),
		"borrowingFee":       t.BorrowingFee.String(),
		"investingFee":       t.InvestingFee.String(),
		"status":             t.Status,
		"takerOrderSide":     t.TakerOrderSide,
		"takerOrderType":     t.TakerOrderType,
		"hash":               t.Hash,
		"createdAt":          t.CreatedAt.Format(time.RFC3339Nano),
		"updatedAt":          t.UpdatedAt.Format(time.RFC3339Nano),
	}

	if (t.CollateralToken != common.Address{}) {
		trade["collateralToken"] = t.CollateralToken.Hex()
	}
	if (t.LendingToken != common.Address{}) {
		trade["lendingToken"] = t.LendingToken.Hex()
	}
	return json.Marshal(trade)
}

// UnmarshalJSON creates a trade object from a json byte string
func (t *LendingTrade) UnmarshalJSON(b []byte) error {
	trade := map[string]interface{}{}

	err := json.Unmarshal(b, &trade)
	if err != nil {
		return err
	}
	if trade["collateralToken"] == nil {
		return errors.New("collateralToken Hash is not set")
	}
	t.CollateralToken = common.HexToAddress(trade["collateralToken"].(string))
	if trade["lendingToken"] == nil {
		return errors.New("lendingToken Hash is not set")
	}
	t.LendingToken = common.HexToAddress(trade["lendingToken"].(string))

	if trade["borrower"] == nil {
		return errors.New("borrower Hash is not set")
	}
	t.Borrower = common.HexToAddress(trade["borrower"].(string))

	if trade["investor"] == nil {
		return errors.New("investor is not set")
	}
	t.Investor = common.HexToAddress(trade["investor"].(string))

	if trade["borrowingOrderHash"] == nil {
		return errors.New("borrowingOrderHash is not set")
	}
	t.BorrowingOrderHash = common.HexToHash(trade["borrowingOrderHash"].(string))

	if trade["hash"] == nil {
		return errors.New("Hash is not set")
	}
	t.Hash = common.HexToHash(trade["hash"].(string))

	if trade["investingOrderHash"] == nil {
		return errors.New("investingOrderHash is not set")
	}
	t.InvestingOrderHash = common.HexToHash(trade["investingOrderHash"].(string))

	if trade["borrowingRelayer"] == nil {
		return errors.New("borrowingRelayer is not set")
	}
	t.BorrowingRelayer = common.HexToAddress(trade["borrowingRelayer"].(string))

	if trade["investingRelayer"] == nil {
		return errors.New("investingRelayer is not set")
	}
	t.InvestingRelayer = common.HexToAddress(trade["investingRelayer"].(string))

	if trade["term"] == nil {
		return errors.New("term is not set")
	}
	t.Term, _ = strconv.ParseUint(trade["term"].(string), 10, 64)

	if trade["interest"] == nil {
		return errors.New("interest is not set")
	}
	t.Interest, _ = strconv.ParseUint(trade["interest"].(string), 10, 64)

	if trade["collateralPrice"] != nil {
		t.CollateralPrice = new(big.Int)
		t.CollateralPrice, _ = t.CollateralPrice.SetString(trade["collateralPrice"].(string), 10)
	}
	if trade["liquidationPrice"] != nil {
		t.LiquidationPrice = new(big.Int)
		t.LiquidationPrice, _ = t.LiquidationPrice.SetString(trade["liquidationPrice"].(string), 10)
	}
	if trade["amount"] != nil {
		t.Amount = new(big.Int)
		t.Amount, _ = t.Amount.SetString(trade["amount"].(string), 10)
	}
	if trade["borrowingFee"] != nil {
		t.BorrowingFee = new(big.Int)
		t.BorrowingFee, _ = t.BorrowingFee.SetString(trade["borrowingFee"].(string), 10)
	}
	if trade["investingFee"] != nil {
		t.InvestingFee = new(big.Int)
		t.InvestingFee, _ = t.InvestingFee.SetString(trade["investingFee"].(string), 10)
	}
	if trade["status"] != nil {
		t.Status = trade["status"].(string)
	}
	if trade["takerOrderSide"] != nil {
		t.TakerOrderSide = trade["takerOrderSide"].(string)
	}
	if trade["takerOrderType"] != nil {
		t.TakerOrderType = trade["takerOrderType"].(string)
	}
	if trade["createdAt"] != nil {
		tm, _ := time.Parse(time.RFC3339Nano, trade["createdAt"].(string))
		t.CreatedAt = tm
	}
	if trade["updateAt"] != nil {
		tm, _ := time.Parse(time.RFC3339Nano, trade["updateAt"].(string))
		t.UpdatedAt = tm
	}

	return nil
}

// LendingTradeBSON lending trade mongo
type LendingTradeBSON struct {
	Borrower               string        `bson:"borrower" json:"borrower"`
	Investor               string        `bson:"investor" json:"investor"`
	LendingToken           string        `bson:"lendingToken" json:"lendingToken"`
	CollateralToken        string        `bson:"collateralToken" json:"collateralToken"`
	BorrowingOrderHash     string        `bson:"borrowingOrderHash" json:"borrowingOrderHash"`
	InvestingOrderHash     string        `bson:"investingOrderHash" json:"investingOrderHash"`
	BorrowingRelayer       string        `bson:"borrowingRelayer" json:"borrowingRelayer"`
	InvestingRelayer       string        `bson:"investingRelayer" json:"investingRelayer"`
	Term                   string        `bson:"term" json:"term"`
	Interest               string        `bson:"interest" json:"interest"`
	CollateralPrice        string        `bson:"collateralPrice" json:"collateralPrice"`
	LiquidationPrice       string        `bson:"liquidationPrice" json:"liquidationPrice"`
	LiquidationTime        string        `bson:"liquidationTime" json:"liquidationTime"`
	CollateralLockedAmount string        `bson:"collateralLockedAmount" json:"collateralLockedAmount"`
	DepositRate            string        `bson:"depositRate" json:"depositRate"`
	Amount                 string        `bson:"amount" json:"amount"`
	BorrowingFee           string        `bson:"borrowingFee" json:"borrowingFee"`
	InvestingFee           string        `bson:"investingFee" json:"investingFee"`
	Status                 string        `bson:"status" json:"status"`
	TakerOrderSide         string        `bson:"takerOrderSide" json:"takerOrderSide"`
	TakerOrderType         string        `bson:"takerOrderType" json:"takerOrderType"`
	MakerOrderType         string        `bson:"makerOrderType" json:"makerOrderType"`
	TradeID                bson.ObjectId `bson:"tradeID" json:"tradeID"`
	Hash                   string        `bson:"hash" json:"hash"`
	TxHash                 string        `bson:"txHash" json:"txHash"`
	ExtraData              string        `bson:"extraData" json:"extraData"`
	UpdatedAt              time.Time     `bson:"updatedAt" json:"updatedAt"`
}

// GetBSON for monggo insert
func (t *LendingTrade) GetBSON() (interface{}, error) {
	return bson.M{
		"$setOnInsert": bson.M{
			"createdAt": t.CreatedAt,
		},
		"$set": LendingTradeBSON{
			Borrower:               t.Borrower.Hex(),
			Investor:               t.Investor.Hex(),
			LendingToken:           t.LendingToken.Hex(),
			CollateralToken:        t.CollateralToken.Hex(),
			BorrowingOrderHash:     t.BorrowingOrderHash.Hex(),
			InvestingOrderHash:     t.InvestingOrderHash.Hex(),
			BorrowingRelayer:       t.BorrowingRelayer.Hex(),
			InvestingRelayer:       t.InvestingRelayer.Hex(),
			Term:                   strconv.FormatUint(t.Term, 10),
			Interest:               strconv.FormatUint(t.Interest, 10),
			CollateralPrice:        t.CollateralPrice.String(),
			LiquidationPrice:       t.LiquidationPrice.String(),
			LiquidationTime:        strconv.FormatUint(t.LiquidationTime, 10),
			CollateralLockedAmount: t.CollateralLockedAmount.String(),
			DepositRate:            t.DepositRate.String(),
			Amount:                 t.Amount.String(),
			BorrowingFee:           t.BorrowingFee.String(),
			InvestingFee:           t.InvestingFee.String(),
			Status:                 t.Status,
			TakerOrderSide:         t.TakerOrderSide,
			TakerOrderType:         t.TakerOrderType,
			MakerOrderType:         t.MakerOrderType,
			TradeID:                t.TradeID,
			Hash:                   t.Hash.Hex(),
			TxHash:                 t.TxHash.Hex(),
			ExtraData:              t.ExtraData,
			UpdatedAt:              t.UpdatedAt,
		},
	}, nil
}

// SetBSON get monggo record
func (t *LendingTrade) SetBSON(raw bson.Raw) error {
	decoded := new(LendingTradeBSON)

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	t.TradeID = decoded.TradeID
	t.Borrower = common.HexToAddress(decoded.Borrower)
	t.Investor = common.HexToAddress(decoded.Investor)
	t.LendingToken = common.HexToAddress(decoded.LendingToken)
	t.CollateralToken = common.HexToAddress(decoded.CollateralToken)
	t.BorrowingOrderHash = common.HexToHash(decoded.BorrowingOrderHash)
	t.InvestingOrderHash = common.HexToHash(decoded.InvestingOrderHash)
	t.BorrowingRelayer = common.HexToAddress(decoded.BorrowingRelayer)
	t.InvestingRelayer = common.HexToAddress(decoded.InvestingRelayer)
	term, err := strconv.ParseInt(decoded.Term, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse lendingItem.term. Err: %v", err)
	}
	t.Term = uint64(term)
	interest, err := strconv.ParseInt(decoded.Interest, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse lendingItem.interest. Err: %v", err)
	}
	t.Interest = uint64(interest)
	t.CollateralPrice = utils.ToBigInt(decoded.CollateralPrice)
	t.LiquidationPrice = utils.ToBigInt(decoded.LiquidationPrice)
	liquidationTime, err := strconv.ParseInt(decoded.LiquidationTime, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse lendingItem.LiquidationTime. Err: %v", err)
	}
	t.LiquidationTime = uint64(liquidationTime)
	t.CollateralLockedAmount = utils.ToBigInt(decoded.CollateralLockedAmount)
	t.DepositRate = utils.ToBigInt(decoded.DepositRate)
	t.Amount = utils.ToBigInt(decoded.Amount)
	t.BorrowingFee = utils.ToBigInt(decoded.BorrowingFee)
	t.InvestingFee = utils.ToBigInt(decoded.InvestingFee)
	t.Status = decoded.Status
	t.TakerOrderSide = decoded.TakerOrderSide
	t.TakerOrderType = decoded.TakerOrderType
	t.MakerOrderType = decoded.MakerOrderType
	t.ExtraData = decoded.ExtraData
	t.Hash = common.HexToHash(decoded.Hash)
	t.TxHash = common.HexToHash(decoded.TxHash)
	t.UpdatedAt = decoded.UpdatedAt

	return nil
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

// ComputeHash returns hashes the trade
// The OrderHash, Amount, Taker and TradeNonce attributes must be
// set before attempting to compute the trade hash
func (t *LendingTrade) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(t.Borrower.Bytes())
	sha.Write(t.Investor.Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// LendingTradeChangeEvent event for changing mongo watch
type LendingTradeChangeEvent struct {
	ID                interface{}   `bson:"_id"`
	OperationType     string        `bson:"operationType"`
	FullDocument      *LendingTrade `bson:"fullDocument,omitempty"`
	Ns                evNamespace   `bson:"ns"`
	DocumentKey       M             `bson:"documentKey"`
	UpdateDescription *updateDesc   `bson:"updateDescription,omitempty"`
}
