package rabbitmq

import (
	"encoding/json"

	"github.com/tomochain/tomox-sdk/types"
)

// SubscribeQueue subscribe queue
func (c *Connection) SubscribeQueue(fn func(*types.EngineResponse) error, queue string) error {
	ch := c.GetChannel("erSub")
	q := c.GetQueue(ch, queue)

	go func() {
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		if err != nil {
			logger.Fatal("Failed to register a consumer:", err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				var res *types.EngineResponse
				err := json.Unmarshal(d.Body, &res)
				if err != nil {
					logger.Error(err)
					continue
				}
				go fn(res)
			}
		}()

		<-forever
	}()
	return nil
}

// SubscribeOrderResponses subscribe order responses
func (c *Connection) SubscribeOrderResponses(fn func(*types.EngineResponse) error) error {
	return c.SubscribeQueue(fn, "orderResponse")
}

// SubscribeTradeResponses subscribe trade responses
func (c *Connection) SubscribeTradeResponses(fn func(*types.EngineResponse) error) error {
	return c.SubscribeQueue(fn, "tradeResponse")
}

// SubscribeEngineResponses subscribe engine responses
func (c *Connection) SubscribeEngineResponses(fn func(*types.EngineResponse) error) error {
	ch := c.GetChannel("erSub")
	q := c.GetQueue(ch, "engineResponse")

	go func() {
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		if err != nil {
			logger.Fatal("Failed to register a consumer:", err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				var res *types.EngineResponse
				err := json.Unmarshal(d.Body, &res)
				if err != nil {
					logger.Error(err)
					continue
				}
				go fn(res)
			}
		}()

		<-forever
	}()
	return nil
}

// PublishMessage publish message to rabbitmq
func (c *Connection) PublishMessage(res *types.EngineResponse, queue string) error {
	ch := c.GetChannel("erPub")
	q := c.GetQueue(ch, queue)

	bytes, err := json.Marshal(res)
	if err != nil {
		logger.Error("Failed to marshal engine response: ", err)
		return err
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error("Failed to publish order: ", err)
		return err
	}

	return nil
}

// PublishOrderResponse publish order response to queue
func (c *Connection) PublishOrderResponse(res *types.EngineResponse) error {
	return c.PublishMessage(res, "orderResponse")
}

// PublishTradeResponse publish trade response to queue
func (c *Connection) PublishTradeResponse(res *types.EngineResponse) error {
	return c.PublishMessage(res, "tradeResponse")
}

// PublishEngineResponse publish engine response to queue
func (c *Connection) PublishEngineResponse(res *types.EngineResponse) error {
	ch := c.GetChannel("erPub")
	q := c.GetQueue(ch, "engineResponse")

	bytes, err := json.Marshal(res)
	if err != nil {
		logger.Error("Failed to marshal engine response: ", err)
		return err
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error("Failed to publish order: ", err)
		return err
	}

	return nil
}
