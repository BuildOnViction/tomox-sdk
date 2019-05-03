package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomodex/interfaces"
	"github.com/tomochain/tomodex/types"
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

// GetByUserAddress fetches all the orders placed by passed user address
func (s *NotificationService) GetByUserAddress(addr common.Address, limit ...int) ([]*types.Notification, error) {
	return s.NotificationDao.GetByUserAddress(addr, limit...)
}
