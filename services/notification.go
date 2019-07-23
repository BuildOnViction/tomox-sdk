package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomoxsdk/interfaces"
	"github.com/tomochain/tomoxsdk/types"
)

// NotificationService struct with daos required, responsible for communicating with dao
// NotificationService functions are responsible for interacting with dao and implements business logic.
type NotificationService struct {
	NotificationDao interfaces.NotificationDao
}

// NewNotificationService returns a new instance of NewNotificationService
func NewNotificationService(
	notificationDao interfaces.NotificationDao,
) *NotificationService {
	return &NotificationService{
		NotificationDao: notificationDao,
	}
}

// Create inserts a new notification into the database
func (s *NotificationService) Create(n *types.Notification) ([]*types.Notification, error) {
	notifications, err := s.NotificationDao.Create(n)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return notifications, nil
}

// GetAll fetches all the notifications from db
func (s *NotificationService) GetAll() ([]types.Notification, error) {
	return s.NotificationDao.GetAll()
}

// GetByUserAddress fetches all the notifications related to user address
func (s *NotificationService) GetByUserAddress(addr common.Address, limit int, offset int) ([]*types.Notification, error) {
	return s.NotificationDao.GetByUserAddress(addr, limit, offset)
}

// GetSortDecByUserAddress function fetches list of orders from order collection based on user address, result sorted by created date.
// Returns array of notification type struct
func (s *NotificationService) GetSortDecByUserAddress(addr common.Address, limit int, offset int) ([]*types.Notification, error) {
	return s.NotificationDao.GetSortDecByUserAddress(addr, limit, offset)
}

// GetByID fetches the detailed document of a notification using its mongo ID
func (s *NotificationService) GetByID(id bson.ObjectId) (*types.Notification, error) {
	return s.NotificationDao.GetByID(id)
}

// Update updates the detailed document of a notification using its mongo ID
func (s *NotificationService) Update(n *types.Notification) (*types.Notification, error) {
	updated, err := s.NotificationDao.FindAndModify(n.ID, n)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return updated, nil
}

// MarkRead update UNREAD status to READ status
func (s *NotificationService) MarkRead(id bson.ObjectId) error {
	return s.NotificationDao.MarkRead(id)
}

// MarkUnRead update READ status to UNREAD status
func (s *NotificationService) MarkUnRead(id bson.ObjectId) error {
	return s.NotificationDao.MarkUnRead(id)
}

// MarkAllRead update all UNREAD status to READ status
func (s *NotificationService) MarkAllRead(addr common.Address) error {
	return s.NotificationDao.MarkAllRead(addr)
}
