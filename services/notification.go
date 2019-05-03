package services

import (
	"github.com/tomochain/tomodex/interfaces"
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
