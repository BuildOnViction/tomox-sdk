package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
)

type Notification struct {
	ID        bson.ObjectId  `json:"-" bson:"_id"`
	Recipient common.Address `json:"recipient" bson:"recipient"`
	Message   string         `json:"message" bson:"message"`
	Type      string         `json:"type" bson:"type"`
	Status    string         `json:"status" bson:"status"`
	CreatedAt time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" bson:"updatedAt"`
}

type NotificationBSONUpdate struct {
	*Notification
}
