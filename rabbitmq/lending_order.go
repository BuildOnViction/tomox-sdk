package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

// SubscribeLendingOrders sub a order consumer
func (c *Connection) SubscribeLendingOrders(fn func(*Message) error) error {
	ch := c.GetChannel("lendingOrderSubscribe")
	if ch == nil {
		return errors.New("Fail to open lendingorderSubscribe chanel")
	}
	q := c.GetQueue(ch, "lending_order")
	if q == nil {
		return errors.New("Fail to open lending order queue")
	}
	go func() {
		msgs, err := c.Consume(ch, q)
		if err != nil {
			logger.Error(err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				msg := &Message{}
				err := json.Unmarshal(d.Body, msg)
				if err != nil {
					logger.Error(err)
					continue
				}

				go fn(msg)
			}
		}()

		<-forever
	}()
	return nil
}

// PublishLendingOrderMessage publish message to queue
func (c *Connection) PublishLendingOrderMessage(o *types.LendingOrder) error {
	b, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.PublishLendingOrder(&Message{
		Type: "NEW_LENDING_ORDER",
		Data: b,
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// PublishCancelLendingOrderMessage publish cancel message to queue
func (c *Connection) PublishCancelLendingOrderMessage(o *types.LendingOrder) error {
	b, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.PublishOrder(&Message{
		Type: "CANCEL_LENDING_ORDER",
		Data: b,
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// PublishLendingOrder publish a lending order
func (c *Connection) PublishLendingOrder(order *Message) error {
	ch := c.GetChannel("lendingOrderPublish")
	q := c.GetQueue(ch, "lending_order")

	bytes, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Failed to marshal order: ", err)
		return errors.New("Failed to marshal order: " + err.Error())
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
