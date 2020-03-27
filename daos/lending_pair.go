package daos

import (
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/types"
)

// LendingPairDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type LendingPairDao struct {
	collectionName string
	dbName         string
}

// LendingPairDaoOption option
type LendingPairDaoOption = func(*LendingPairDao) error

// LendingPairDaoDBOption option
func LendingPairDaoDBOption(dbName string) func(dao *LendingPairDao) error {
	return func(dao *LendingPairDao) error {
		dao.dbName = dbName
		return nil
	}
}

// NewLendingPairDao returns a new instance of AddressDao
func NewLendingPairDao(options ...LendingPairDaoOption) *LendingPairDao {
	dao := &LendingPairDao{}
	dao.collectionName = "lending_pairs"
	dao.dbName = app.Config.DBName

	for _, op := range options {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}
	index := mgo.Index{
		Key:    []string{"lendingTokenAddress", "term"},
		Unique: true,
	}
	err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return dao
}

// Create function performs the DB insertion task for pair collection
func (dao *LendingPairDao) Create(pair *types.LendingPair) error {
	pair.ID = bson.NewObjectId()
	pair.CreatedAt = time.Now()
	pair.UpdatedAt = time.Now()

	err := db.Create(dao.dbName, dao.collectionName, pair)
	return err
}

// GetAll function fetches all the pairs in the pair collection of mongodb.
// for GetAll return continous memory
func (dao *LendingPairDao) GetAll() ([]types.LendingPair, error) {
	var res []types.LendingPair
	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &res)
	if err != nil {
		return nil, err
	}

	ret := []types.LendingPair{}
	keys := make(map[string]bool)

	for _, it := range res {
		code := it.LendingTokenAddress.Hex() + "::" + string(it.Term)
		if _, value := keys[code]; !value {
			keys[code] = true
			ret = append(ret, it)
		}
	}

	return ret, nil
}

func (dao *LendingPairDao) GetAllByCoinbase(addr common.Address) ([]types.LendingPair, error) {
	var res []types.LendingPair
	err := db.Get(dao.dbName, dao.collectionName, bson.M{"relayerAddress": addr.Hex()}, 0, 0, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetByID function fetches details of a pair using pair's mongo ID.
func (dao *LendingPairDao) GetByID(id bson.ObjectId) (*types.LendingPair, error) {
	var response *types.LendingPair
	err := db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return response, err
}

// DeleteByLendingKey delete token by term and lending token address
func (dao *LendingPairDao) DeleteByLendingKey(term uint64, lendingAddress common.Address) error {
	query := bson.M{"lendingTokenAddress": lendingAddress.Hex(), "term": strconv.FormatUint(term, 10)}
	return db.RemoveItem(dao.dbName, dao.collectionName, query)
}

// GetByLendingID get pair from lending token and term
func (dao *LendingPairDao) GetByLendingID(term uint64, lendingAddress common.Address) (*types.LendingPair, error) {
	var res types.LendingPair
	query := bson.M{"lendingTokenAddress": lendingAddress.Hex(), "term": strconv.FormatUint(term, 10)}
	err := db.GetOne(dao.dbName, dao.collectionName, query, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
