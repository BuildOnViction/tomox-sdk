package daos

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tomochain/backend-matching-engine/app"
	"github.com/tomochain/backend-matching-engine/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	schemaVersionKey        = "swap_schema_version"
	ethereumAddressIndexKey = "ethereum_address_index"
	ethereumLastBlockKey    = "ethereum_last_block"
)

// DepositDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type DepositDao struct {
	collectionName string
	dbName         string
}

// NewBalanceDao returns a new instance of AddressDao
func NewDepositDao() *DepositDao {
	dbName := app.Config.DBName
	// we save deposit information in swarm feed and use config for retrieving params.
	collection := "config"
	index := mgo.Index{
		Key:    []string{"key"},
		Unique: true,
	}

	err := db.Session.DB(dbName).C(collection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return &DepositDao{collection, dbName}
}

// Create function performs the DB insertion task for Balance collection
func (dao *DepositDao) GetAssociationByChainAddress(chain types.Chain, address *common.Address) (*types.AddressAssociation, error) {
	// get from feed
	return nil, nil
}

func (dao *DepositDao) AddProcessedTransaction(chain types.Chain, transactionID string, receivingAddress *common.Address) (bool, error) {
	return false, nil
}

func (dao *DepositDao) GetAssociationByTomochainPublicKey(tomochainPublicKey *common.Address) (*types.AddressAssociation, error) {
	return nil, nil
}

func (dao *DepositDao) AddRecoveryTransaction(sourceAccount *common.Address, txEnvelope string) error {
	return nil
}

func (dao *DepositDao) GetSchemaVersion() uint64 {
	// get version
	var response []types.KeyValue
	err := db.Get(dao.dbName, dao.collectionName, bson.M{"key": schemaVersionKey}, 0, 1, &response)
	if err != nil {
		logger.Error(err)
		return types.SwapSchemaVersion
	}

	version, _ := strconv.ParseUint(response[0].Value, 10, 32)
	return version
}

func (dao *DepositDao) getAddressIndexKey(chain types.Chain) (string, error) {
	switch chain {
	case types.ChainEthereum:
		return ethereumAddressIndexKey, nil
	default:
		return "", errors.New("Invalid chain")
	}
}

func (dao *DepositDao) GetAddressIndex(chain types.Chain) (uint64, error) {
	key, err := dao.getAddressIndexKey(chain)
	if err != nil {
		return 0, err
	}

	var response []types.KeyValue
	err = db.Get(dao.dbName, dao.collectionName, bson.M{"key": key}, 0, 1, &response)
	if err != nil {
		return 0, err
	}

	index, _ := strconv.ParseUint(response[0].Value, 10, 32)
	return index, nil
}

func (dao *DepositDao) IncrementAddressIndex(chain types.Chain) error {
	// update database
	key, err := dao.getAddressIndexKey(chain)
	if err != nil {
		return err
	}

	err = db.Update(dao.dbName, dao.collectionName, bson.M{"key": key}, bson.M{
		"$inc": bson.M{
			key: 1,
		},
	})

	return err
}

// func (dao *DepositDao) UpdateTokenBalance(owner, token common.Address, tokenBalance *types.TokenBalance) error {
// 	q := bson.M{
// 		"address": owner.Hex(),
// 	}

// 	updateQuery := bson.M{
// 		"$set": bson.M{
// 			"tokenBalances." + token.Hex() + ".balance":        tokenBalance.Balance.String(),
// 			"tokenBalances." + token.Hex() + ".allowance":      tokenBalance.Allowance.String(),
// 			"tokenBalances." + token.Hex() + ".lockedBalance":  tokenBalance.LockedBalance.String(),
// 			"tokenBalances." + token.Hex() + ".pendingBalance": tokenBalance.PendingBalance.String(),
// 		},
// 	}

// 	err := db.Update(dao.dbName, dao.collectionName, q, updateQuery)
// 	return err
// }

// Drop drops all the order documents in the current database
func (dao *DepositDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}

// ResetBlockCounters changes last processed bitcoin and ethereum block to default value.
// Used in stress tests.
func (dao *DepositDao) ResetBlockCounters() error {
	// _, err = keyValueStore.Update(nil, map[string]interface{}{"key": ethereumLastBlockKey}).Set("value", 0).Exec()
	// if err != nil {
	// 	return errors.Wrap(err, "Error reseting `ethereumLastBlockKey`")
	// }

	return nil
}

// AddRecoveryTransaction inserts recovery account ID and transaction envelope
// func (dao *DepositDao) AddRecoveryTransaction(sourceAccount string, txEnvelope string) error
