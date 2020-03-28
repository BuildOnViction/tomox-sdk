package daos

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/types"
)

type RelayerDao struct {
	collectionName string
	dbName         string
}

func NewRelayerDao() *RelayerDao {
	dbName := app.Config.DBName
	collection := "relayers"
	index := mgo.Index{
		Key:    []string{"address", "domain"},
		Unique: true,
	}

	err := db.Session.DB(dbName).C(collection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return &RelayerDao{collection, dbName}
}

func (dao *RelayerDao) Create(a *types.Relayer) error {
	a.ID = bson.NewObjectId()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	err := db.Create(dao.dbName, dao.collectionName, a)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *RelayerDao) FindOrCreate(addr common.Address) (*types.Relayer, error) {
	a := &types.Relayer{Address: addr}
	query := bson.M{"address": addr.Hex()}
	updated := &types.Relayer{}

	change := mgo.Change{
		Update:    types.RelayerBSONUpdate{Relayer: a},
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

func (dao *RelayerDao) GetAll() (res []types.Relayer, err error) {
	err = db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &res)
	return
}

func (dao *RelayerDao) GetByID(id bson.ObjectId) (*types.Relayer, error) {
	res := []types.Relayer{}
	q := bson.M{"_id": id}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &res[0], nil
}

func (dao *RelayerDao) GetByAddress(owner common.Address) (*types.Relayer, error) {
	res := []types.Relayer{}
	q := bson.M{"address": owner.Hex()}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

func (dao *RelayerDao) GetByHost(host string) (*types.Relayer, error) {
	res := []types.Relayer{}
	q := bson.M{"domain": host}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

func (dao *RelayerDao) DeleteByAddress(addr common.Address) error {
	query := bson.M{"address": addr.Hex()}
	return db.RemoveItem(dao.dbName, dao.collectionName, query)
}

func (dao *RelayerDao) UpdateByAddress(addr common.Address, relayer *types.Relayer) error {
	q := bson.M{"address": addr.Hex()}

	update := bson.M{
		"$set": bson.M{
			"makeFee":    relayer.MakeFee.String(),
			"takeFee":    relayer.TakeFee.String(),
			"lendingFee": relayer.LendingFee.String(),
			"domain":     relayer.Domain,
		},
	}
	err := db.Update(dao.dbName, dao.collectionName, q, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
