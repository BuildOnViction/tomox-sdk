package daos

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomodex/app"
	"github.com/tomochain/tomodex/types"
)

// StopOrderDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type StopOrderDao struct {
	collectionName string
	dbName         string
}

type StopOrderDaoOption = func(*StopOrderDao) error

// NewOrderDao returns a new instance of OrderDao
func NewStopOrderDao(opts ...StopOrderDaoOption) *StopOrderDao {
	dao := &StopOrderDao{}
	dao.collectionName = "stop_orders"
	dao.dbName = app.Config.DBName

	for _, op := range opts {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}

	index := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	}

	i1 := mgo.Index{
		Key: []string{"userAddress"},
	}

	i2 := mgo.Index{
		Key: []string{"status"},
	}

	i3 := mgo.Index{
		Key: []string{"baseToken"},
	}

	i4 := mgo.Index{
		Key: []string{"quoteToken"},
	}

	i5 := mgo.Index{
		Key:       []string{"stopPrice"},
		Collation: &mgo.Collation{NumericOrdering: true, Locale: "en"},
	}

	i6 := mgo.Index{
		Key: []string{"baseToken", "quoteToken", "stopPrice"},
	}

	i7 := mgo.Index{
		Key: []string{"side", "status"},
	}

	i8 := mgo.Index{
		Key: []string{"baseToken", "quoteToken", "side", "status"},
	}

	err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i1)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i2)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i3)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i4)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i5)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i6)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i7)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i8)
	if err != nil {
		panic(err)
	}

	return dao
}

// Create function performs the DB insertion task for Order collection
func (dao *StopOrderDao) Create(o *types.StopOrder) error {
	o.ID = bson.NewObjectId()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	if o.Status == "" {
		o.Status = types.OrderStatusOpen
	}

	err := db.Create(dao.dbName, dao.collectionName, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Update function performs the DB updations task for Order collection
// corresponding to a particular order ID
func (dao *StopOrderDao) Update(id bson.ObjectId, so *types.StopOrder) error {
	so.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": id}, so)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//UpdateByHash updates fields that are considered updateable for an order.
func (dao *StopOrderDao) UpdateByHash(h common.Hash, so *types.StopOrder) error {
	so.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{
		"stopPrice":    so.StopPrice.String(),
		"limitPrice":   so.LimitPrice.String(),
		"amount":       so.Amount.String(),
		"status":       so.Status,
		"filledAmount": so.FilledAmount.String(),
		"makeFee":      so.MakeFee.String(),
		"takeFee":      so.TakeFee.String(),
		"updatedAt":    so.UpdatedAt,
	}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *StopOrderDao) Upsert(id bson.ObjectId, o *types.StopOrder) error {
	o.UpdatedAt = time.Now()

	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": id}, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *StopOrderDao) UpsertByHash(h common.Hash, so *types.StopOrder) error {
	_, err := db.Upsert(
		dao.dbName,
		dao.collectionName,
		bson.M{"hash": h.Hex()},
		types.StopOrderBSONUpdate{
			StopOrder: so,
		},
	)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *StopOrderDao) UpdateAllByHash(h common.Hash, so *types.StopOrder) error {
	so.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"hash": h.Hex()}, so)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *StopOrderDao) FindAndModify(h common.Hash, so *types.StopOrder) (*types.StopOrder, error) {
	so.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	updated := &types.StopOrder{}
	change := mgo.Change{
		Update: types.StopOrderBSONUpdate{
			StopOrder: so,
		},
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}

	err := db.FindAndModify(dao.dbName, dao.collectionName, query, change, &updated)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return updated, nil
}

func (dao *StopOrderDao) UpdateOrderStatus(h common.Hash, status string) error {
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{
		"status": status,
	}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *StopOrderDao) UpdateOrderStatusesByHashes(status string, hashes ...common.Hash) ([]*types.StopOrder, error) {
	hexes := make([]string, 0)
	for _, h := range hashes {
		hexes = append(hexes, h.Hex())
	}

	query := bson.M{"hash": bson.M{"$in": hexes}}
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"status":    status,
		},
	}

	err := db.UpdateAll(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	orders := make([]*types.StopOrder, 0)
	err = db.Get(dao.dbName, dao.collectionName, query, 0, 0, &orders)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	return orders, nil
}

// Drop drops all the order documents in the current database
func (dao *StopOrderDao) Drop() error {
	err := db.DropCollection(dao.dbName, dao.collectionName)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
