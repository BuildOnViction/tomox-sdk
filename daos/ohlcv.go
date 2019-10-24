package daos

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/types"
)

// OHLCVDao for ohlcv database struct
type OHLCVDao struct {
	collectionName string
	dbName         string
}

// NewOHLCVDao for init ohlcv database struct
func NewOHLCVDao() *OHLCVDao {
	dbName := app.Config.DBName
	collection := "ohlcv"
	index := mgo.Index{
		Key:    []string{"timestamp", "uint", "duration"},
		Unique: true,
	}
	db.Session.DB(dbName).C(collection).EnsureIndex(index)
	return &OHLCVDao{collection, dbName}
}

// Create function performs the DB insertion task for token collection
func (dao *OHLCVDao) Create(tick *types.Tick) error {
	err := db.Create(dao.dbName, dao.collectionName, tick)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// GetAll function fetches all the tokens in the token collection of mongodb.
func (dao *OHLCVDao) GetAll() ([]types.Tick, error) {
	var response []types.Tick
	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

// GetOhlcv get OHLCV data
func (dao *OHLCVDao) GetOhlcv(pairs []types.PairAddresses, duration int64, unit string, start, end int64) ([]*types.Tick, error) {
	res := make([]*types.Tick, 0)

	match := bson.M{
		"timestamp": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	if len(pairs) >= 1 {
		or := make([]bson.M, 0)

		for _, pair := range pairs {
			or = append(or, bson.M{
				"$and": []bson.M{
					{
						"baseToken":  pair.BaseToken.Hex(),
						"quoteToken": pair.QuoteToken.Hex(),
					},
				},
			},
			)
		}

		match["$or"] = or
	}
	sort := []string{"-timestamp"}
	err := db.GetAndSort(dao.dbName, dao.collectionName, match, sort, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return res, nil
}
