package daos

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomodex/app"
	"github.com/tomochain/tomodex/types"
	"time"
)

type NotificationDao struct {
	collectionName string
	dbName         string
}

func NewNotificationDao() *NotificationDao {
	dbName := app.Config.DBName
	collection := "notifications"

	return &NotificationDao{
		collectionName: collection,
		dbName:         dbName,
	}
}

// Create function performs the DB insertion task for notification collection
// It accepts 1 or more notifications as input.
// All the notifications are inserted in one query itself.
func (dao *NotificationDao) Create(notifications ...*types.Notification) error {
	y := make([]interface{}, len(notifications))

	for _, notification := range notifications {
		notification.ID = bson.NewObjectId()
		notification.CreatedAt = time.Now()
		notification.UpdatedAt = time.Now()
		y = append(y, notification)
	}

	err := db.Create(dao.dbName, dao.collectionName, y...)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
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