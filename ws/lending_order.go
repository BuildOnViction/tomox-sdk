package ws

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/types"
)

// LendingOrderConnection array for lending connections
type LendingOrderConnection []*Client

var lockLendingOrder = &sync.Mutex{}

var lendingOrderConnections map[string]LendingOrderConnection

// GetLendingOrderConnections returns the connection associated with an order ID
func GetLendingOrderConnections(a common.Address) LendingOrderConnection {
	c := lendingOrderConnections[a.Hex()]
	if c == nil {
		logger.Warning("No connection found")
		return nil
	}

	return lendingOrderConnections[a.Hex()]
}

// LendingOrderSocketUnsubscribeHandler unsubscrible order
func LendingOrderSocketUnsubscribeHandler(a common.Address) func(client *Client) {
	return func(client *Client) {
		logger.Info("In unsubscription handler")
		orderConnection := lendingOrderConnections[a.Hex()]
		if orderConnection == nil {
			logger.Info("No subscriptions")
		}

		if orderConnection != nil {
			logger.Info("%v connections before unsubscription", len(lendingOrderConnections[a.Hex()]))
			lockLendingOrder.Lock()
			for i, c := range orderConnection {
				if client == c {
					orderConnection = append(orderConnection[:i], orderConnection[i+1:]...)
				}
			}
			lockLendingOrder.Unlock()

		}

		lendingOrderConnections[a.Hex()] = orderConnection
		logger.Info("%v connections after unsubscription", len(lendingOrderConnections[a.Hex()]))
	}
}

// RegisterLendingOrderConnection registers a connection with and orderID.
// It is called whenever a message is recieved over order channel
func RegisterLendingOrderConnection(a common.Address, c *Client) {
	logger.Info("Registering new order connection")

	if lendingOrderConnections == nil {
		lendingOrderConnections = make(map[string]LendingOrderConnection)
	}

	if lendingOrderConnections[a.Hex()] == nil {
		logger.Info("Registering a new order connection")
		lendingOrderConnections[a.Hex()] = LendingOrderConnection{c}
		RegisterConnectionUnsubscribeHandler(c, LendingOrderSocketUnsubscribeHandler(a))
		logger.Info("Number of connections for this address: %v", len(lendingOrderConnections))
	}

	if lendingOrderConnections[a.Hex()] != nil {

		if !isClientConnected(lendingOrderConnections[a.Hex()], c) {
			logger.Info("Registering a new order connection")
			lockLendingOrder.Lock()
			lendingOrderConnections[a.Hex()] = append(lendingOrderConnections[a.Hex()], c)
			lockLendingOrder.Unlock()
			RegisterConnectionUnsubscribeHandler(c, LendingOrderSocketUnsubscribeHandler(a))
			logger.Info("Number of connections for this address: %v", len(lendingOrderConnections))
		}
	}
}

// SendLendingOrderMessage send lending order message
func SendLendingOrderMessage(msgType types.SubscriptionEvent, a common.Address, payload interface{}) {
	conn := GetLendingOrderConnections(a)
	if conn == nil {
		return
	}

	for _, c := range conn {
		c.SendMessage(OrderChannel, msgType, payload)
	}
}
