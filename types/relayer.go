package types

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/go-ozzo/ozzo-validation"
)

// Relayer corresponds to a single Ethereum address. It contains a list of token balances for that address
type Relayer struct {
	ID        bson.ObjectId  `json:"-" bson:"_id"`
	Address   common.Address `json:"address" bson:"address"`
	Domain    string         `json:"domain" bson:"domain"`
	CreatedAt time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" bson:"updatedAt"`
}

// GetBSON implements bson.Getter
func (a *Relayer) GetBSON() (interface{}, error) {
	ar := RelayerRecord{
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
	a.ID = decoded.ID
	a.Domain = decoded.Domain
	a.CreatedAt = decoded.CreatedAt
	a.UpdatedAt = decoded.UpdatedAt

	return nil
}

// MarshalJSON implements the json.Marshal interface
func (a *Relayer) MarshalJSON() ([]byte, error) {
	account := map[string]interface{}{
		"id":        a.ID,
		"address":   a.Address,
		"domain":    a.Domain,
		"createdAt": a.CreatedAt.String(),
		"updatedAt": a.UpdatedAt.String(),
	}

	return json.Marshal(account)
}

func (a *Relayer) UnmarshalJSON(b []byte) error {
	account := map[string]interface{}{}
	err := json.Unmarshal(b, &account)
	if err != nil {
		return err
	}

	if account["id"] != nil && bson.IsObjectIdHex(account["id"].(string)) {
		a.ID = bson.ObjectIdHex(account["id"].(string))
	}

	if account["address"] != nil {
		a.Address = common.HexToAddress(account["address"].(string))
	}

	if account["domain"] != nil {
		a.Domain = account["domain"].(string)
	}

	return nil
}

// Validate enforces the account model
func (a Relayer) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Address, validation.Required),
	)
}

// RelayerRecord corresponds to what is stored in the DB. big.Ints are encoded as strings
type RelayerRecord struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Address   string        `json:"address" bson:"address"`
	Domain    string        `json:"domain" bson:"domain"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
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
