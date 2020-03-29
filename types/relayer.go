package types

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/tomochain/tomox-sdk/utils/math"
)

// Relayer corresponds to a single Ethereum address. It contains a list of token balances for that address
type Relayer struct {
	ID         bson.ObjectId  `json:"-" bson:"_id"`
	RID        int            `json:"rid" bson:"rid"`
	Owner      common.Address `json:"owner" bson:"owner"`
	Deposit    *big.Int       `json:"deposit" bson:"deposit"`
	Address    common.Address `json:"address" bson:"address"`
	Domain     string         `json:"domain" bson:"domain"`
	MakeFee    *big.Int       `json:"makeFee,omitempty" bson:"makeFee,omitempty"`
	TakeFee    *big.Int       `json:"takeFee,omitempty" bson:"makeFee,omitempty"`
	LendingFee *big.Int       `json:"lendingFee,omitempty" bson:"lendingFee,omitempty"`
	CreatedAt  time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt" bson:"updatedAt"`
}

// GetBSON implements bson.Getter
func (a *Relayer) GetBSON() (interface{}, error) {
	ar := RelayerRecord{
		RID:       a.RID,
		Owner:     a.Owner.Hex(),
		Deposit:   a.Deposit.String(),
		Domain:    a.Domain,
		Address:   a.Address.Hex(),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	if a.ID.Hex() == "" {
		ar.ID = bson.NewObjectId()
	} else {
		ar.ID = a.ID
	}

	if a.MakeFee != nil {
		ar.MakeFee = a.MakeFee.String()
	}

	if a.TakeFee != nil {
		ar.TakeFee = a.TakeFee.String()
	}

	if a.LendingFee != nil {
		ar.LendingFee = a.LendingFee.String()
	}

	return ar, nil
}

// SetBSON implemenets bson.Setter
func (a *Relayer) SetBSON(raw bson.Raw) error {
	decoded := &RelayerRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	a.Address = common.HexToAddress(decoded.Address)
	a.Owner = common.HexToAddress(decoded.Owner)
	a.Deposit = math.ToBigInt(decoded.Deposit)
	a.RID = decoded.RID
	a.ID = decoded.ID
	a.Domain = decoded.Domain
	a.CreatedAt = decoded.CreatedAt
	a.UpdatedAt = decoded.UpdatedAt
	if decoded.MakeFee != "" {
		a.MakeFee = math.ToBigInt(decoded.MakeFee)
	}

	if decoded.TakeFee != "" {
		a.TakeFee = math.ToBigInt(decoded.TakeFee)
	}

	if decoded.LendingFee != "" {
		a.LendingFee = math.ToBigInt(decoded.LendingFee)
	}

	return nil
}

// MarshalJSON implements the json.Marshal interface
func (a *Relayer) MarshalJSON() ([]byte, error) {
	relayer := map[string]interface{}{
		"id":        a.ID,
		"address":   a.Address.Hex(),
		"domain":    a.Domain,
		"rid":       a.RID,
		"owner":     a.Owner.Hex(),
		"deposit":   a.Deposit.String(),
		"createdAt": a.CreatedAt.String(),
		"updatedAt": a.UpdatedAt.String(),
	}

	if a.MakeFee != nil {
		relayer["makeFee"] = a.MakeFee.String()
	}

	if a.TakeFee != nil {
		relayer["takeFee"] = a.TakeFee.String()
	}

	if a.LendingFee != nil {
		relayer["lendingFee"] = a.LendingFee.String()
	}

	return json.Marshal(relayer)
}

func (a *Relayer) UnmarshalJSON(b []byte) error {
	relayer := map[string]interface{}{}
	err := json.Unmarshal(b, &relayer)
	if err != nil {
		return err
	}

	if relayer["id"] != nil && bson.IsObjectIdHex(relayer["id"].(string)) {
		a.ID = bson.ObjectIdHex(relayer["id"].(string))
	}

	if relayer["address"] != nil {
		a.Address = common.HexToAddress(relayer["address"].(string))
	}

	if relayer["owner"] != nil {
		a.Owner = common.HexToAddress(relayer["owner"].(string))
	}

	if relayer["deposit"] != nil {
		a.Deposit = math.ToBigInt(relayer["deposit"].(string))
	}

	if relayer["rid"] != nil {
		a.RID = relayer["rid"].(int)
	}

	if relayer["domain"] != nil {
		a.Domain = relayer["domain"].(string)
	}

	if relayer["makeFee"] != nil {
		a.MakeFee = math.ToBigInt(relayer["makeFee"].(string))
	}

	if relayer["takeFee"] != nil {
		a.TakeFee = math.ToBigInt(relayer["takeFee"].(string))
	}

	if relayer["lendingFee"] != nil {
		a.LendingFee = math.ToBigInt(relayer["lendingFee"].(string))
	}

	return nil
}

// Validate enforces the relayer model
func (a Relayer) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Address, validation.Required),
	)
}

// RelayerRecord corresponds to what is stored in the DB. big.Ints are encoded as strings
type RelayerRecord struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	RID        int           `json:"rid" bson:"rid"`
	Owner      string        `json:"owner" bson:"owner"`
	Deposit    string        `json:"deposit" bson:"deposit"`
	Address    string        `json:"address" bson:"address"`
	Domain     string        `json:"domain" bson:"domain"`
	MakeFee    string        `json:"makeFee,omitempty" bson:"makeFee,omitempty"`
	TakeFee    string        `json:"takeFee,omitempty" bson:"takeFee,omitempty"`
	LendingFee string        `json:"lendingFee,omitempty" bson:"lendingFee,omitempty"`
	CreatedAt  time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type RelayerBSONUpdate struct {
	*Relayer
}

func (a *RelayerBSONUpdate) GetBSON() (interface{}, error) {
	now := time.Now()

	set := bson.M{
		"updatedAt": now,
		"address":   a.Address,
	}

	setOnInsert := bson.M{
		"_id":       bson.NewObjectId(),
		"createdAt": now,
	}

	update := bson.M{
		"$set":         set,
		"$setOnInsert": setOnInsert,
	}

	return update, nil
}
