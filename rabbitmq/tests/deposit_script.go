package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/tomochain/tomoxsdk/app"
	"github.com/tomochain/tomoxsdk/rabbitmq"
	"github.com/tomochain/tomoxsdk/types"
)

func SetupTest() *rabbitmq.Connection {
	err := app.LoadConfig("./config", "test")
	if err != nil {
		panic(err)
	}

	return rabbitmq.InitConnection(app.Config.RabbitMQURL)
}

func main() {

	if len(os.Args) > 1 {
		cmd := os.Args[1]

		if cmd == "send" {
			sendMessage()
			return
		}
	}

	signalInterrupt := make(chan os.Signal, 1)
	signal.Notify(signalInterrupt, os.Interrupt)

	go processTransaction()

	<-signalInterrupt
}

func sendMessage() {
	connection := SetupTest()

	transaction := &types.DepositTransaction{
		Chain:         types.ChainEthereum,
		TransactionID: "30",
		AssetCode:     "ETH",
		Amount:        "100",
	}

	connection.PublishDepositTransaction(transaction)
}

func processTransaction() {

	connection := SetupTest()

	msgs, err := connection.QueuePoolDepositTransactions()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for transaction := range msgs {
		fmt.Printf("Got transaction :%v\n", transaction)
	}
}
