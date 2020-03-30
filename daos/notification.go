package daos

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

type NotificationDao struct {
	collectionName string
	dbName         string
}

func NewNotificationDao() *NotificationDao {
	dao := &NotificationDao{}
	dao.collectionName = "notifications"
	dao.dbName = app.Config.DBName

	i1 := mgo.Index{
		Key: []string{"recipient"},
	}

	err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i1)

	i2 := mgo.Index{
		Key:         []string{"createdAt"},
		Background:  true,
		ExpireAfter: time.Duration(30*24*60*60) * time.Second, // 30 days
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i2)

	if err != nil {
		logger.Warning("Index failed", err)
	}

	return dao
}

// Create function performs the DB insertion task for notification collection
// It accepts 1 or more notifications as input.
// All the notifications are inserted in one query itself.
func (dao *NotificationDao) Create(notifications ...*types.Notification) ([]*types.Notification, error) {
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
func (dao *NotificationDao) GetAll() ([]types.Notification, error) {
	var response []types.Notification

	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

// GetByUserAddress function fetches list of orders from order collection based on user address.
// Returns array of Order type struct
func (dao *NotificationDao) GetByUserAddress(addr common.Address, limit int, offset int) ([]*types.Notification, error) {
	if limit == 0 {
		limit = 10 // Get last 10 records
	}

	var res []*types.Notification
	q := bson.M{"recipient": addr.Hex()}

	err := db.Get(dao.dbName, dao.collectionName, q, offset, limit, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return []*types.Notification{}, nil
	}

	return res, nil
}

// GetSortDecByUserAddress function fetches list of orders from order collection based on user address, result sorted by created date.
// Returns array of notification type struct
func (dao *NotificationDao) GetSortDecByUserAddress(addr common.Address, limit int, offset int) ([]*types.Notification, error) {
	if limit == 0 {
		limit = 10 // Get last 10 records
	}

	var res []*types.Notification
	q := bson.M{"recipient": addr.Hex()}
	sort := []string{"-createdAt"}
	err := db.GetAndSort(dao.dbName, dao.collectionName, q, sort, offset, limit, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return []*types.Notification{}, nil
	}

	return res, nil
}

// GetByID function fetches details of a notification based on its mongo id
func (dao *NotificationDao) GetByID(id bson.ObjectId) (*types.Notification, error) {
	var response *types.Notification

	err := db.GetByID(dao.dbName, dao.collectionName, id, &response)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

func (dao *NotificationDao) FindAndModify(id bson.ObjectId, n *types.Notification) (*types.Notification, error) {
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

func (dao *NotificationDao) Update(n *types.Notification) error {
	n.UpdatedAt = time.Now()
	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": n.ID}, n)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *NotificationDao) Upsert(id bson.ObjectId, n *types.Notification) error {
	n.UpdatedAt = time.Now()

	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": id}, n)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *NotificationDao) Delete(notifications ...*types.Notification) error {
	ids := make([]bson.ObjectId, 0)
	for _, n := range notifications {
		ids = append(ids, n.ID)
	}

	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"_id": bson.M{"$in": ids}})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *NotificationDao) DeleteByIds(ids ...bson.ObjectId) error {
	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"_id": bson.M{"$in": ids}})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Aggregate function calls the aggregate pipeline of mongodb
func (dao *NotificationDao) Aggregate(q []bson.M) ([]*types.Notification, error) {
	var res []*types.Notification

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// Drop drops all the order documents in the current database
func (dao *NotificationDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}

// MarkStatus update UNREAD status to READ status
func (dao *NotificationDao) MarkStatus(id bson.ObjectId, status string) error {
	query := bson.M{"_id": id}
	updated := &types.Notification{}
	changeData := bson.M{}
	changeData["status"] = status
	changeData["updatedAt"] = time.Now()
	changeDataSet := bson.M{"$set": changeData}

	change := mgo.Change{
		Update:    changeDataSet,
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}

	err := db.FindAndModify(dao.dbName, dao.collectionName, query, change, &updated)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// MarkRead update UNREAD status to READ status
func (dao *NotificationDao) MarkRead(id bson.ObjectId) error {
	return dao.MarkStatus(id, types.StatusRead)
}

// MarkUnRead update READ status to UNREAD status
func (dao *NotificationDao) MarkUnRead(id bson.ObjectId) error {
	return dao.MarkStatus(id, types.StatusUnread)
}

// MarkAllRead update all UNREAD status to READ status
func (dao *NotificationDao) MarkAllRead(addr common.Address) error {
	query := bson.M{"recipient": addr.Hex(), "status": types.StatusUnread}
	changeData := bson.M{}
	changeData["status"] = types.StatusRead
	changeData["updatedAt"] = time.Now()
	changeDataSet := bson.M{"$set": changeData}
	changeInfo, err := db.ChangeAll(dao.dbName, dao.collectionName, query, changeDataSet)
	if err != nil {
		logger.Error(err)
		return err
	}
	if changeInfo.Matched == 0 {
		return errors.New("User address not found or all user notificaions have been read")
	}
	if changeInfo.Updated < changeInfo.Matched {
		return errors.New("Update process is not completed")

	}
	return nil
}
