package daos

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tidwall/gjson"
	"github.com/tomochain/tomoxsdk/app"
	"github.com/tomochain/tomoxsdk/errors"
	"github.com/tomochain/tomoxsdk/types"
)

type FiatPriceDao struct {
	collectionName string
	dbName         string
}

// NewFiatPriceDao returns a new instance of FiatPriceDao.
func NewFiatPriceDao() *FiatPriceDao {
	dbName := app.Config.DBName
	collection := "fiat_price"

	return &FiatPriceDao{
		collectionName: collection,
		dbName:         dbName,
	}
}

func (dao *FiatPriceDao) GetLatestQuotes() (map[string]float64, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/cryptocurrency/quotes/latest?symbol=%s&convert=USD", app.Config.CoinmarketcapAPIUrl, app.Config.SupportedCurrencies)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("X-CMC_PRO_API_KEY", app.Config.CoinmarketcapAPIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	status := gjson.Get(string(body), "status")
	statusErrorCode := status.Get("error_code")
	statusErrorMessage := status.Get("error_message")

	if statusErrorCode.Int() != 0 {
		logger.Error(statusErrorMessage.String())
		return nil, errors.New(statusErrorMessage.String())
	}

	data := gjson.Get(string(body), "data")
	result := make(map[string]float64)
	data.ForEach(func(key, value gjson.Result) bool {
		result[key.String()] = value.Get("quote.USD.price").Float()
		return true // keep iterating
	})

	return result, nil
}

// Create function performs the DB insertion task for notification collection
// It accepts 1 or more notifications as input.
// All the notifications are inserted in one query itself.
func (dao *FiatPriceDao) Create(notifications ...*types.Notification) ([]*types.Notification, error) {
	y := make([]interface{}, len(notifications))
	result := make([]*types.Notification, len(notifications))

	for _, notification := range notifications {
		notification.ID = bson.NewObjectId()
		notification.CreatedAt = time.Now()
		notification.UpdatedAt = time.Now()
		y = append(y, notification)
		result = append(result, notification)
	}

	err := db.Create(dao.dbName, dao.collectionName, y...)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return result, nil
}

// GetAll function fetches all the notifications in the notification collection of mongodb.
func (dao *FiatPriceDao) GetAll() ([]types.Notification, error) {
	var response []types.Notification

	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

func (dao *FiatPriceDao) FindAndModify(id bson.ObjectId, n *types.Notification) (*types.Notification, error) {
	n.UpdatedAt = time.Now()
	query := bson.M{"_id": id}
	updated := &types.Notification{}
	change := mgo.Change{
		Update:    types.NotificationBSONUpdate{Notification: n},
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

func (dao *FiatPriceDao) Update(n *types.Notification) error {
	n.UpdatedAt = time.Now()
	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": n.ID}, n)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *FiatPriceDao) Upsert(id bson.ObjectId, n *types.Notification) error {
	n.UpdatedAt = time.Now()

	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": id}, n)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Aggregate function calls the aggregate pipeline of mongodb
func (dao *FiatPriceDao) Aggregate(q []bson.M) ([]*types.Notification, error) {
	var res []*types.Notification

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// Drop drops all the order documents in the current database
func (dao *FiatPriceDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}
