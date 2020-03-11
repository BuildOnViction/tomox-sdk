package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/utils/math"
)

const (
	BORROW     = "BORROW"
	LEND       = "INVEST"
	TypeMarket = "MO"
	TypeLimit  = "LO"

	LendingStatusOpen          = "OPEN"
	LendingStatusPartialFilled = "PARTIAL_FILLED"
	LendingStatusFilled        = "FILLED"
	LendingStatusRejected      = "REJECTED"
	LendingStatusCancelled     = "CANCELLED"
)

// LendingOrder contains the data related to an lending sent by the user
type LendingOrder struct {
	ID              bson.ObjectId  `json:"id" bson:"_id"`
	Quantity        *big.Int       `bson:"quantity" json:"quantity"`
	Interest        uint64         `bson:"interest" json:"interest"`
	Term            uint64         `bson:"term" json:"term"`
	Side            string         `bson:"side" json:"side"`
	Type            string         `bson:"type" json:"type"`
	LendingToken    common.Address `bson:"lendingToken" json:"lendingToken"`
	CollateralToken common.Address `bson:"collateralToken" json:"collateralToken"`
	FilledAmount    *big.Int       `bson:"filledAmount" json:"filledAmount"`
	Status          string         `bson:"status" json:"status"`
	UserAddress     common.Address `bson:"userAddress" json:"userAddress"`
	RelayerAddress  common.Address `bson:"relayerAddress" json:"relayerAddress"`
	Signature       *Signature     `bson:"signature" json:"signature"`
	Hash            common.Hash    `bson:"hash" json:"hash"`
	TxHash          common.Hash    `bson:"txHash" json:"txHash"`
	Nonce           *big.Int       `bson:"nonce" json:"nonce"`
	CreatedAt       time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time      `bson:"updatedAt" json:"updatedAt"`
	LendingID       uint64         `bson:"lendingId" json:"lendingId"`
	ExtraData       string         `bson:"extraData" json:"extraData"`
	LendingTradeID  uint64         `bson:"tradeId" json:"tradeId"`
	AutoTopUp       uint64         `json:"autoTopUp" json:"autoTopUp"`
	Key             string         `json:"key" bson:"key"`
}

// LendingRes use for api
type LendingRes struct {
	Total        int             `json:"total" bson:"total"`
	LendingItems []*LendingOrder `json:"lendings" bson:"lendings"`
}

// LendingSpec contains field for filter
type LendingSpec struct {
	UserAddress     string
	CollateralToken string
	LendingToken    string
	Status          string
	Side            string
	Type            string
	DateFrom        int64
	DateTo          int64
	Hash            string
}

// Validate Verify userAddress, collateralToken, lendingToken, etc. conditions are working
func (o *LendingOrder) Validate() error {
	if o.RelayerAddress != common.HexToAddress(app.Config.Tomochain["exchange_address"]) {
		return errors.New("LendingOrder 'exchange_address' parameter is incorrect")
	}

	if (o.UserAddress == common.Address{}) {
		return errors.New("LendingOrder 'userAddress' parameter is required")
	}

	if o.Nonce == nil {
		return errors.New("LendingOrder 'nonce' parameter is required")
	}

	if (o.CollateralToken == common.Address{}) {
		return errors.New("LendingOrder 'collateralToken' parameter is required")
	}

	if (o.LendingToken == common.Address{}) {
		return errors.New("LendingOrder 'lendingToken' parameter is required")
	}

	if o.Quantity == nil {
		return errors.New("LendingOrder 'quantity' parameter is required")
	}
	if o.Term == 0 {
		return errors.New("LendingOrder 'term' parameter is required")
	}
	if o.Type == TypeLimit {
		if o.Interest == 0 {
			return errors.New("LendingOrder 'interest' parameter is required")
		}
	}
	if o.Side != LEND && o.Side != BORROW {
		return errors.New("LendingOrder 'side' should be 'LEND' or 'BORROW', but got: '" + o.Side + "'")
	}

	if o.Signature == nil {
		return errors.New("LendingOrder 'signature' parameter is required")
	}

	if math.IsStrictlySmallerThan(o.Nonce, big.NewInt(0)) {
		return errors.New("LendingOrder 'nonce' parameter should be positive")
	}

	if math.IsEqualOrSmallerThan(o.Quantity, big.NewInt(0)) {
		return errors.New("LendingOrder 'quantity' parameter should be strictly positive")
	}

	valid, err := o.VerifySignature()
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("LendingOrder 'signature' parameter is invalid")
	}

	return nil
}

// ComputeHash calculates the orderRequest hash
func (o *LendingOrder) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(o.RelayerAddress.Bytes())
	sha.Write(o.UserAddress.Bytes())
	sha.Write(o.CollateralToken.Bytes())
	sha.Write(o.LendingToken.Bytes())
	sha.Write(common.BigToHash(o.Quantity).Bytes())
	sha.Write(common.BigToHash(big.NewInt(int64(o.Term))).Bytes())
	if o.Type == TypeLimit {
		sha.Write(common.BigToHash(big.NewInt(int64(o.Interest))).Bytes())
	}
	sha.Write([]byte(o.Side))
	sha.Write([]byte(o.Status))
	sha.Write([]byte(o.Type))
	sha.Write(common.BigToHash(o.Nonce).Bytes())
	if o.Side == BORROW {
		autoTopUp := int64(o.AutoTopUp)
		sha.Write(common.BigToHash(big.NewInt(autoTopUp)).Bytes())
	}

	return common.BytesToHash(sha.Sum(nil))
}

// VerifySignature checks that the orderRequest signature corresponds to the address in the userAddress field
func (o *LendingOrder) VerifySignature() (bool, error) {
	o.Hash = o.ComputeHash()

	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		o.Hash.Bytes(),
	)

	address, err := o.Signature.Verify(common.BytesToHash(message))
	if err != nil {
		return false, err
	}
	logger.Debug(address.Hex())
	if address != o.UserAddress {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}

// Process pre-process data
func (o *LendingOrder) Process() error {
	if o.FilledAmount == nil {
		o.FilledAmount = big.NewInt(0)
	}
	if o.Type != TypeMarket && o.Type != TypeLimit {
		o.Type = TypeLimit
	}

	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

// PairCode get orderbook code
func (o *LendingOrder) PairCode() (string, error) {
	if (o.LendingToken == common.Address{}) {
		return "", errors.New("Pair name is required")
	}

	return fmt.Sprint(o.Term) + "::" + o.LendingToken.Hex(), nil
}

// MarshalJSON implements the json.Marshal interface
func (o *LendingOrder) MarshalJSON() ([]byte, error) {
	lending := map[string]interface{}{
		"relayerAddress":  o.RelayerAddress,
		"userAddress":     o.UserAddress,
		"collateralToken": o.CollateralToken,
		"lendingToken":    o.LendingToken,
		"side":            o.Side,
		"type":            o.Type,
		"status":          o.Status,
		"quantity":        o.Quantity.String(),
		"term":            strconv.FormatUint(o.Term, 10),
		"interest":        strconv.FormatUint(o.Interest, 10),
		"createdAt":       o.CreatedAt.Format(time.RFC3339Nano),
		"updatedAt":       o.UpdatedAt.Format(time.RFC3339Nano),
		"lendingId":       strconv.FormatUint(o.LendingID, 10),
		"tradeId":         strconv.FormatUint(o.LendingTradeID, 10),
		"autoTopUp":       strconv.FormatUint(o.AutoTopUp, 10),
		"key":             o.Key,
	}

	if o.FilledAmount != nil {
		lending["filledAmount"] = o.FilledAmount.String()
	}

	if o.Hash.Hex() != "" {
		lending["hash"] = o.Hash.Hex()
	}

	if o.Nonce != nil {
		lending["nonce"] = o.Nonce.String()
	}
	if o.Signature != nil {
		lending["signature"] = map[string]interface{}{
			"V": o.Signature.V,
			"R": o.Signature.R,
			"S": o.Signature.S,
		}
	}
	return json.Marshal(lending)
}

// UnmarshalJSON : write custom logic to unmarshal bytes to LendingOrder
func (o *LendingOrder) UnmarshalJSON(b []byte) error {
	lending := map[string]interface{}{}
	err := json.Unmarshal(b, &lending)
	if err != nil {
		return err
	}
	if lending["id"] != nil && bson.IsObjectIdHex(lending["id"].(string)) {
		o.ID = bson.ObjectIdHex(lending["id"].(string))
	}

	if lending["relayerAddress"] != nil {
		o.RelayerAddress = common.HexToAddress(lending["relayerAddress"].(string))
	}

	if lending["userAddress"] != nil {
		o.UserAddress = common.HexToAddress(lending["userAddress"].(string))
	}

	if lending["collateralToken"] != nil {
		o.CollateralToken = common.HexToAddress(lending["collateralToken"].(string))
	}

	if lending["lendingToken"] != nil {
		o.LendingToken = common.HexToAddress(lending["lendingToken"].(string))
	}

	if lending["quantity"] != nil {
		o.Quantity = math.ToBigInt(lending["quantity"].(string))
	}

	if lending["term"] != nil {
		t, err := strconv.ParseInt(string(lending["term"].(string)), 10, 64)
		if err == nil {
			o.Term = uint64(t)
		}
	}
	if lending["interest"] != nil {
		i, err := strconv.ParseInt(string(lending["interest"].(string)), 10, 64)
		if err == nil {
			o.Interest = uint64(i)
		}
	}

	if lending["filledAmount"] != nil {
		o.FilledAmount = math.ToBigInt(lending["filledAmount"].(string))
	}

	if lending["nonce"] != nil {
		o.Nonce = math.ToBigInt(lending["nonce"].(string))
	}

	if lending["hash"] != nil {
		o.Hash = common.HexToHash(lending["hash"].(string))
	}

	if lending["side"] != nil {
		o.Side = lending["side"].(string)
	}

	if lending["type"] != nil {
		o.Type = lending["type"].(string)
	}

	if lending["status"] != nil {
		o.Status = lending["status"].(string)
	}

	if lending["signature"] != nil {
		signature := lending["signature"].(map[string]interface{})
		o.Signature = &Signature{
			V: byte(signature["V"].(float64)),
			R: common.HexToHash(signature["R"].(string)),
			S: common.HexToHash(signature["S"].(string)),
		}
	}

	if lending["createdAt"] != nil {
		t, _ := time.Parse(time.RFC3339Nano, lending["createdAt"].(string))
		o.CreatedAt = t
	}

	if lending["updatedAt"] != nil {
		t, _ := time.Parse(time.RFC3339Nano, lending["updatedAt"].(string))
		o.UpdatedAt = t
	}

	if lending["lendingId"] != nil {
		lendingID, err := strconv.ParseInt(lending["lendingId"].(string), 10, 64)
		if err != nil {
			logger.Error(err)
		}
		o.LendingID = uint64(lendingID)
	}
	if lending["tradeId"] != nil {
		lendingTradeID, err := strconv.ParseInt(lending["tradeId"].(string), 10, 64)
		if err != nil {
			logger.Error(err)
		}
		o.LendingTradeID = uint64(lendingTradeID)
	}
	if lending["autoTopUp"] != nil {
		autoTopUp, err := strconv.ParseInt(lending["autoTopUp"].(string), 10, 64)
		if err != nil {
			logger.Error(err)
		}
		o.AutoTopUp = uint64(autoTopUp)
	}
	if lending["key"] != nil {
		o.Key = lending["key"].(string)
	}

	return nil
}

// GetBSON return bson
func (o *LendingOrder) GetBSON() (interface{}, error) {
	or := LendingRecord{
		RelayerAddress:  o.RelayerAddress.Hex(),
		UserAddress:     o.UserAddress.Hex(),
		CollateralToken: o.CollateralToken.Hex(),
		LendingToken:    o.LendingToken.Hex(),
		Status:          o.Status,
		Side:            o.Side,
		Type:            o.Type,
		Hash:            o.Hash.Hex(),
		Quantity:        o.Quantity.String(),
		Term:            strconv.FormatUint(o.Term, 10),
		Interest:        strconv.FormatUint(o.Interest, 10),
		Nonce:           o.Nonce.String(),
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
		LendingID:       strconv.FormatUint(o.LendingID, 10),
		Key:             o.Key,
	}

	if o.ID.Hex() == "" {
		or.ID = bson.NewObjectId()
	} else {
		or.ID = o.ID
	}

	if o.FilledAmount != nil {
		or.FilledAmount = o.FilledAmount.String()
	}

	if o.Signature != nil {
		or.Signature = &SignatureRecord{
			V: o.Signature.V,
			R: o.Signature.R.Hex(),
			S: o.Signature.S.Hex(),
		}
	}

	return or, nil
}

// SetBSON for database
func (o *LendingOrder) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID              bson.ObjectId    `json:"id,omitempty" bson:"_id"`
		RelayerAddress  string           `json:"relayerAddress" bson:"relayerAddress"`
		UserAddress     string           `json:"userAddress" bson:"userAddress"`
		CollateralToken string           `json:"collateralToken" bson:"collateralToken"`
		LendingToken    string           `json:"lendingToken" bson:"lendingToken"`
		Term            string           `json:"term" bson:"term"`
		Interest        string           `json:"interest" bson:"interest"`
		Status          string           `json:"status" bson:"status"`
		Side            string           `json:"side" bson:"side"`
		Type            string           `json:"type" bson:"type"`
		Hash            string           `json:"hash" bson:"hash"`
		Quantity        string           `json:"quantity" bson:"quantity"`
		FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
		Nonce           string           `json:"nonce" bson:"nonce"`
		Signature       *SignatureRecord `json:"signature" bson:"signature"`
		CreatedAt       time.Time        `json:"createdAt" bson:"createdAt"`
		UpdatedAt       time.Time        `json:"updatedAt" bson:"updatedAt"`
		LendingID       string           `json:"lendingId" bson:"lendingId"`
		Key             string           `json:"key" bson:"key"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		logger.Error(err)
		return err
	}

	o.ID = decoded.ID
	o.RelayerAddress = common.HexToAddress(decoded.RelayerAddress)
	o.UserAddress = common.HexToAddress(decoded.UserAddress)
	o.CollateralToken = common.HexToAddress(decoded.CollateralToken)
	o.LendingToken = common.HexToAddress(decoded.LendingToken)
	o.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	o.Nonce = math.ToBigInt(decoded.Nonce)
	o.Status = decoded.Status
	o.Side = decoded.Side
	o.Type = decoded.Type
	o.Hash = common.HexToHash(decoded.Hash)
	term, err := strconv.ParseUint(decoded.Term, 10, 64)
	if err == nil {
		o.Term = term
	}
	interest, err := strconv.ParseUint(decoded.Interest, 10, 64)
	if err == nil {
		o.Interest = interest
	}
	if decoded.Quantity != "" {
		o.Quantity = math.ToBigInt(decoded.Quantity)
	}

	if decoded.FilledAmount != "" {
		o.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	}

	if decoded.Quantity != "" {
		o.Quantity = math.ToBigInt(decoded.Quantity)
	}

	if decoded.Signature != nil {
		o.Signature = &Signature{
			V: byte(decoded.Signature.V),
			R: common.HexToHash(decoded.Signature.R),
			S: common.HexToHash(decoded.Signature.S),
		}
	}
	o.CreatedAt = decoded.CreatedAt
	o.UpdatedAt = decoded.UpdatedAt

	lendingID, err := strconv.ParseInt(decoded.LendingID, 10, 64)
	if err != nil {
		logger.Error(err)
	}
	o.LendingID = uint64(lendingID)
	o.Key = decoded.Key

	return nil
}

// LendingRecord is the object that will be saved in the database
type LendingRecord struct {
	ID              bson.ObjectId    `json:"id" bson:"_id"`
	UserAddress     string           `json:"userAddress" bson:"userAddress"`
	RelayerAddress  string           `json:"relayerAddress" bson:"relayerAddress"`
	CollateralToken string           `json:"collateralToken" bson:"collateralToken"`
	LendingToken    string           `json:"lendingToken" bson:"lendingToken"`
	Term            string           `json:"term" bson:"term"`
	Interest        string           `json:"interest" bson:"interest"`
	Status          string           `json:"status" bson:"status"`
	Side            string           `json:"side" bson:"side"`
	Type            string           `json:"type" bson:"type"`
	Hash            string           `json:"hash" bson:"hash"`
	Quantity        string           `json:"quantity" bson:"quantity"`
	FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
	Nonce           string           `json:"nonce" bson:"nonce"`
	Signature       *SignatureRecord `json:"signature,omitempty" bson:"signature"`
	CreatedAt       time.Time        `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt" bson:"updatedAt"`
	LendingID       string           `json:"lendingId,omitempty" bson:"lendingId"`
	NextOrder       string           `json:"nextOrder,omitempty" bson:"nextOrder"`
	PrevOrder       string           `json:"prevOrder,omitempty" bson:"prevOrder"`
	OrderList       string           `json:"orderList,omitempty" bson:"orderList"`
	Key             string           `json:"key" bson:"key"`
}

// LendingOrderChangeEvent data format for changing data records
type LendingOrderChangeEvent struct {
	ID                interface{}   `bson:"_id"`
	OperationType     string        `bson:"operationType"`
	FullDocument      *LendingOrder `bson:"fullDocument,omitempty"`
	Ns                evNamespace   `bson:"ns"`
	DocumentKey       M             `bson:"documentKey"`
	UpdateDescription *updateDesc   `bson:"updateDescription,omitempty"`
}

// LendingOrderCancel for cancelled lending order
type LendingOrderCancel struct {
	LendingHash    common.Hash    `json:"lendingHash"`
	Nonce          *big.Int       `json:"nonce"`
	Hash           common.Hash    `json:"hash"`
	LendingID      uint64         `json:"lendingId"`
	Status         string         `json:"status"`
	UserAddress    common.Address `json:"userAddress"`
	RelayerAddress common.Address `json:"relayerAddress"`
	Term           uint64         `json:"term"`
	Interest       uint64         `json:"interest"`
	Signature      *Signature     `json:"signature"`
}

// MarshalJSON returns the json encoded byte array representing the LendingOrderCancel struct
func (oc *LendingOrderCancel) MarshalJSON() ([]byte, error) {
	orderCancel := map[string]interface{}{
		"lendingHash": oc.LendingHash,
		"nonce":       oc.Nonce,
		"hash":        oc.Hash,
		"signature": map[string]interface{}{
			"V": oc.Signature.V,
			"R": oc.Signature.R,
			"S": oc.Signature.S,
		},
		"lendingId":      oc.LendingID,
		"userAddress":    oc.UserAddress,
		"relayerAddress": oc.RelayerAddress,
		"status":         oc.Status,
	}

	return json.Marshal(orderCancel)
}

// UnmarshalJSON creates an LendingOrderCancel object from a json byte string
func (oc *LendingOrderCancel) UnmarshalJSON(b []byte) error {
	parsed := map[string]interface{}{}

	err := json.Unmarshal(b, &parsed)
	if err != nil {
		return err
	}

	// if parsed["lendingHash"] == nil {
	// 	return errors.New("Lending Hash is missing")
	// }
	// oc.LendingHash = common.HexToHash(parsed["lendingHash"].(string))

	if parsed["hash"] == nil {
		return errors.New("Hash is missing")
	}
	oc.Hash = common.HexToHash(parsed["hash"].(string))

	if parsed["nonce"] == nil {
		return errors.New("Nonce is missing")
	}
	oc.Nonce = math.ToBigInt(parsed["nonce"].(string))

	if parsed["status"] == nil {
		return errors.New("Status is missing")
	}
	oc.Status = parsed["status"].(string)

	if parsed["lendingId"] == nil {
		return errors.New("lendingId is missing")
	}
	lendingID, err := strconv.ParseUint(parsed["lendingId"].(string), 10, 64)
	if err != nil {
		return err
	}
	oc.LendingID = lendingID

	if parsed["userAddress"] == nil {
		return errors.New("userAddress is missing")
	}
	oc.UserAddress = common.HexToAddress(parsed["userAddress"].(string))

	if parsed["relayerAddress"] == nil {
		return errors.New("relayerAddress is missing")
	}
	oc.RelayerAddress = common.HexToAddress(parsed["relayerAddress"].(string))

	sig := parsed["signature"].(map[string]interface{})
	oc.Signature = &Signature{
		V: byte(sig["V"].(float64)),
		R: common.HexToHash(sig["R"].(string)),
		S: common.HexToHash(sig["S"].(string)),
	}

	return nil
}
