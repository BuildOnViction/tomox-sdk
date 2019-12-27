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

// LendingTradeDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type LendingTradeDao struct {
	collectionName string
	dbName         string
}

// NewLendingTradeDao returns a new instance of LendingTradeDao.
func NewLendingTradeDao() *LendingTradeDao {
	dbName := app.Config.DBName
	collection := "lending_trades"

	i1 := mgo.Index{
		Key: []string{"collateralToken"},
	}

	i2 := mgo.Index{
		Key: []string{"lendingToken"},
	}

	i3 := mgo.Index{
		Key: []string{"createdAt"},
	}

	i4 := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
		Sparse: true,
	}

	i5 := mgo.Index{
		Key:    []string{"borrowingHash"},
		Sparse: true,
	}

	i6 := mgo.Index{
		Key:    []string{"investingHash"},
		Sparse: true,
	}

	i7 := mgo.Index{
		Key: []string{"createdAt", "status", "collateralToken", "lendingToken"},
	}

	indexes := []mgo.Index{}
	indexes, err := db.Session.DB(dbName).C(collection).Indexes()
	if err == nil {
		if !existedIndex("index_lending_trade_hash", indexes) {
			db.Session.DB(dbName).C(collection).EnsureIndex(i4)
		}
	}

	db.Session.DB(dbName).C(collection).EnsureIndex(i1)
	db.Session.DB(dbName).C(collection).EnsureIndex(i2)
	db.Session.DB(dbName).C(collection).EnsureIndex(i3)
	db.Session.DB(dbName).C(collection).EnsureIndex(i5)
	db.Session.DB(dbName).C(collection).EnsureIndex(i6)
	db.Session.DB(dbName).C(collection).EnsureIndex(i7)

	return &LendingTradeDao{collection, dbName}
}

// GetCollection get trade collection name
func (dao *LendingTradeDao) GetCollection() *mgo.Collection {
	return db.GetCollection(dao.dbName, dao.collectionName)
}

// Watch changing database
func (dao *LendingTradeDao) Watch() (*mgo.ChangeStream, *mgo.Session, error) {
	return db.Watch(dao.dbName, dao.collectionName, mgo.ChangeStreamOptions{
		FullDocument:   mgo.UpdateLookup,
		MaxAwaitTimeMS: 500,
		BatchSize:      1000,
	})
}

// Create function performs the DB insertion task for trade collection
// It accepts 1 or more trades as input.
// All the trades are inserted in one query itself.
func (dao *LendingTradeDao) Create(trades ...*types.LendingTrade) error {
	y := make([]interface{}, len(trades))

	for _, trade := range trades {
		trade.ID = bson.NewObjectId()
		trade.CreatedAt = time.Now()
		trade.UpdatedAt = time.Now()
		y = append(y, trade)
	}

	err := db.Create(dao.dbName, dao.collectionName, y...)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Update update lending trade record
func (dao *LendingTradeDao) Update(trade *types.LendingTrade) error {
	trade.UpdatedAt = time.Now()
	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": trade.ID}, trade)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Upsert update lending trade record by id
func (dao *LendingTradeDao) Upsert(id bson.ObjectId, t *types.LendingTrade) error {
	t.UpdatedAt = time.Now()

	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": id}, t)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// UpsertByHash update lending trade record by hash
func (dao *LendingTradeDao) UpsertByHash(h common.Hash, t *types.LendingTrade) error {
	t.UpdatedAt = time.Now()

	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"hash": h}, t)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetAll function fetches all the trades in mongodb
func (dao *LendingTradeDao) GetAll() ([]types.LendingTrade, error) {
	var response []types.LendingTrade
	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

// Aggregate function calls the aggregate pipeline of mongodb
func (dao *LendingTradeDao) Aggregate(q []bson.M) ([]*types.Tick, error) {
	var res []*types.Tick

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetByHash fetches the first record that matches a certain hash
func (dao *LendingTradeDao) GetByHash(h common.Hash) (*types.LendingTrade, error) {
	q := bson.M{"hash": h.Hex()}

	res := []*types.LendingTrade{}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res[0], nil
}

// GetLendingTradeByOrderBook get trade by term and lendingToken
func (dao *LendingTradeDao) GetLendingTradeByOrderBook(term uint64, lendingToken common.Address, from, to int64, n int) ([]*types.LendingTrade, error) {
	res := make([]*types.LendingTrade, 0)

	var q bson.M

	if from == 0 || to == 0 {
		q = bson.M{
			"term":         strconv.FormatUint(term, 10),
			"lendingToken": lendingToken.Hex(),
		}
	} else {
		q = bson.M{
			"term":         strconv.FormatUint(term, 10),
			"lendingToken": lendingToken.Hex(),
			"createdAt": bson.M{
				"$gte": strconv.FormatInt(from, 10),
				"$lt":  strconv.FormatInt(to, 10),
			},
		}
	}

	sort := []string{"-createdAt"}

	err := db.GetAndSort(dao.dbName, dao.collectionName, q, sort, 0, n, &res)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetSortedLendingTradesByUserAddress get trade of user address
func (dao *LendingTradeDao) GetSortedLendingTradesByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.LendingTrade, error) {
	if limit == nil {
		limit = []int{types.DefaultLimit}
	}

	var res []*types.LendingTrade
	var q bson.M

	if (bt == common.Address{} || qt == common.Address{}) {
		if from == 0 || to == 0 {
			q = bson.M{
				"$or": []bson.M{
					{"maker": a.Hex()},
					{"taker": a.Hex()},
				},
			}
		} else {
			q = bson.M{
				"createdAt": bson.M{
					"$gte": strconv.FormatInt(from, 10),
					"$lt":  strconv.FormatInt(to, 10),
				},
				"$or": []bson.M{
					{"maker": a.Hex()},
					{"taker": a.Hex()},
				},
			}
		}
	} else {
		if from == 0 || to == 0 {
			q = bson.M{
				"baseToken":  bt.Hex(),
				"quoteToken": qt.Hex(),
				"$or": []bson.M{
					{"maker": a.Hex()},
					{"taker": a.Hex()},
				},
			}
		} else {
			q = bson.M{
				"baseToken":  bt.Hex(),
				"quoteToken": qt.Hex(),
				"createdAt": bson.M{
					"$gte": strconv.FormatInt(from, 10),
					"$lt":  strconv.FormatInt(to, 10),
				},
				"$or": []bson.M{
					{"maker": a.Hex()},
					{"taker": a.Hex()},
				},
			}
		}
	}

	sort := []string{"-createdAt"}

	err := db.GetAndSort(dao.dbName, dao.collectionName, q, sort, 0, limit[0], &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetByUserAddress fetches all the trades corresponding to a particular user address.
func (dao *LendingTradeDao) GetByUserAddress(a common.Address) ([]*types.LendingTrade, error) {
	var res []*types.LendingTrade
	q := bson.M{"$or": []bson.M{{"maker": a.Hex()}, {"taker": a.Hex()}}}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// UpdateTradeStatus update trade status
func (dao *LendingTradeDao) UpdateTradeStatus(h common.Hash, status string) error {
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{"status": status}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetTradeByTime get range trade
func (dao *LendingTradeDao) GetTradeByTime(dateFrom, dateTo int64, pageOffset int, pageSize int) ([]*types.LendingTrade, error) {
	q := bson.M{}

	dateFilter := bson.M{}
	dateFilter["$gte"] = time.Unix(dateFrom, 0)
	dateFilter["$lt"] = time.Unix(dateTo, 0)
	q["createdAt"] = dateFilter

	trades := []*types.LendingTrade{}
	_, err := db.GetEx(dao.dbName, dao.collectionName, q, []string{"-createdAt"}, pageOffset, pageSize, &trades)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return trades, nil
}
