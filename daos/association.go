package daos

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/backend-matching-engine/app"
	"github.com/tomochain/backend-matching-engine/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AssociationDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type AssociationDao struct {
	collectionName string
	dbName         string
}

// NewBalanceDao returns a new instance of AddressDao
func NewAssociationDao() *AssociationDao {
	dbName := app.Config.DBName
	// we save deposit information in swarm feed and use config for retrieving params.
	collection := "associations"
	index := mgo.Index{
		Key:    []string{"chain", "address"},
		Unique: true,
	}

	err := db.Session.DB(dbName).C(collection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return &AssociationDao{collection, dbName}
}

// return the lowercase of the key
func (dao *AssociationDao) getAddressKey(address common.Address) string {
	return "0x" + common.Bytes2Hex(address.Bytes())
}

// Drop drops all the order documents in the current database
func (dao *AssociationDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}

// SaveDepositTransaction update the transaction envelope for association item
func (dao *AssociationDao) SaveDepositTransaction(chain types.Chain, sourceAccount common.Address, txEnvelope string) error {
	// txEnvolope is rlp of result
	err := db.Update(dao.dbName, dao.collectionName, bson.M{
		"chain":   chain.String(),
		"address": sourceAccount,
	}, bson.M{
		"$set": bson.M{
			"txEnvelope": txEnvelope,
		},
	})
	return err
}

func (dao *AssociationDao) GetAssociationByChainAddress(chain types.Chain, userAddress common.Address) (*types.AddressAssociationRecord, error) {
	var response types.AddressAssociationRecord
	err := db.GetOne(dao.dbName, dao.collectionName, bson.M{
		"chain":   chain.String(),
		"address": dao.getAddressKey(userAddress),
	}, &response)

	return &response, err
}

// SaveAssociation using upsert to update for existing users
func (dao *AssociationDao) SaveAssociation(record *types.AddressAssociationRecord) error {
	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{
		"associatedAddress": record.AssociatedAddress,
	}, bson.M{
		"$set": bson.M{
			"associatedAddress": record.AssociatedAddress,
			"chain":             record.Chain,
			"address":           strings.ToLower(record.Address),
		},
	})
	return err
}
