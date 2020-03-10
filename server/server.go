package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/crons"
	"github.com/tomochain/tomox-sdk/daos"
	"github.com/tomochain/tomox-sdk/endpoints"
	"github.com/tomochain/tomox-sdk/engine"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/ethereum"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/relayer"
	"github.com/tomochain/tomox-sdk/services"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

const (
	swaggerUIDir = "/swaggerui/"
)

var logger = utils.Logger

func Start() {
	env := os.Getenv("GO_ENV")

	if err := app.LoadConfig("./config", env); err != nil {
		panic(err)
	}

	utils.InitLogger(app.Config.LogLevel)

	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(err)
	}

	logger.Infof("Server port: %v", app.Config.ServerPort)
	logger.Infof("Tomochain node HTTP url: %v", app.Config.Tomochain["http_url"])
	logger.Infof("Tomochain node WS url: %v", app.Config.Tomochain["ws_url"])
	logger.Infof("MongoDB url: %v", app.Config.MongoURL)
	logger.Infof("RabbitMQ url: %v", app.Config.RabbitMQURL)
	logger.Infof("Exchange contract address: %v", app.Config.Tomochain["exchange_address"])
	logger.Infof("Env: %v", app.Config.Env)

	_, err := daos.InitSession(nil)
	if err != nil {
		panic(err)
	}

	rabbitConn := rabbitmq.InitConnection(app.Config.RabbitMQURL)

	provider := ethereum.NewWebsocketProvider()

	router := NewRouter(provider, rabbitConn)
	// http.Handle("/", router)
	router.HandleFunc("/socket", ws.ConnectionEndpoint)

	// start the server
	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	log.Printf("server %v is started at %v\n", app.Version, address)

	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	panic(http.ListenAndServe(address, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router)))
}

func NewRouter(
	provider *ethereum.EthereumProvider,
	rabbitConn *rabbitmq.Connection,
) *mux.Router {

	r := mux.NewRouter()

	// get daos for dependency injection
	orderDao := daos.NewOrderDao()
	tokenDao := daos.NewTokenDao()
	pairDao := daos.NewPairDao()
	tradeDao := daos.NewTradeDao()
	accountDao := daos.NewAccountDao()
	walletDao := daos.NewWalletDao()
	notificationDao := daos.NewNotificationDao()

	lendingOrderDao := daos.NewLendingOrderDao()
	lendingTradeDao := daos.NewLendingTradeDao()
	lengdingPairDao := daos.NewLendingPairDao()
	// instantiate engine
	eng := engine.NewEngine(rabbitConn, orderDao, tradeDao, pairDao, provider)

	// get services for injection
	accountService := services.NewAccountService(accountDao, tokenDao, pairDao, orderDao, provider)
	ohlcvService := services.NewOHLCVService(tradeDao, pairDao)
	ohlcvService.Init()

	lendingOhlcvService := services.NewLendingOhlcvService(lendingTradeDao)
	lendingOhlcvService.Init()

	tokenService := services.NewTokenService(tokenDao)
	validatorService := services.NewValidatorService(provider, accountDao, orderDao, pairDao)
	pairService := services.NewPairService(pairDao, tokenDao, tradeDao, orderDao, ohlcvService, eng, provider)

	orderService := services.NewOrderService(orderDao, tokenDao, pairDao, accountDao, tradeDao, notificationDao, eng, validatorService, rabbitConn)
	orderService.LoadCache()
	orderBookService := services.NewOrderBookService(pairDao, tokenDao, orderDao, eng)
	tradeService := services.NewTradeService(orderDao, tradeDao, ohlcvService, notificationDao, rabbitConn)

	walletService := services.NewWalletService(walletDao)

	priceBoardService := services.NewPriceBoardService(tokenDao, tradeDao, ohlcvService)
	marketsService := services.NewMarketsService(pairDao, orderDao, tradeDao, ohlcvService, pairService)
	notificationService := services.NewNotificationService(notificationDao)

	// LEDNDING SERVICE
	lendingOrderService := services.NewLendingOrderService(lendingOrderDao, eng, rabbitConn)
	lendingTradeService := services.NewLendingTradeService(lendingOrderDao, lendingTradeDao, rabbitConn)
	lendingOrderbookService := services.NewLendingOrderBookService(lendingOrderDao)

	// deploy http and ws endpoints
	endpoints.ServeInfoResource(r, walletService, tokenService)
	endpoints.ServeAccountResource(r, accountService)
	endpoints.ServeTokenResource(r, tokenService)
	endpoints.ServePairResource(r, pairService)
	endpoints.ServeOrderBookResource(r, orderBookService)
	endpoints.ServeOHLCVResource(r, ohlcvService)

	endpoints.ServeTradeResource(r, tradeService)
	endpoints.ServeOrderResource(r, orderService, accountService)

	endpoints.ServePriceBoardResource(r, priceBoardService)
	endpoints.ServeMarketsResource(r, marketsService, ohlcvService)
	endpoints.ServeNotificationResource(r, notificationService)

	// Endpoint for lending

	endpoints.ServeLendingOrderResource(r, lendingOrderService)
	endpoints.ServeLendingTradeResource(r, lendingTradeService)
	endpoints.ServeLendingOrderBookResource(r, lendingOrderbookService)
	endpoints.ServeLendingOhlcvResource(r, lendingOhlcvService)
	endpoints.ServeLendingTradeResource(r, lendingTradeService)

	exchangeAddress := common.HexToAddress(app.Config.Tomochain["exchange_address"])
	contractAddress := common.HexToAddress(app.Config.Tomochain["exchange_contract_address"])
	lendingContractAddress := common.HexToAddress(app.Config.Tomochain["lending_contract_address"])
	relayerEngine := relayer.NewRelayer(app.Config.Tomochain["http_url"], exchangeAddress, contractAddress, lendingContractAddress)
	relayerService := services.NewRelayerService(relayerEngine, tokenDao, pairDao, lengdingPairDao)
	endpoints.ServeRelayerResource(r, relayerService)

	// Swagger UI
	sh := http.StripPrefix(swaggerUIDir, http.FileServer(http.Dir("."+swaggerUIDir)))
	r.PathPrefix(swaggerUIDir).Handler(sh)

	//initialize rabbitmq subscriptions
	rabbitConn.SubscribeOrders(eng.HandleOrders)
	rabbitConn.SubscribeEngineResponses(orderService.HandleEngineResponse)

	rabbitConn.SubscribeOrderResponses(orderService.HandleEngineResponse)
	rabbitConn.SubscribeTradeResponses(tradeService.HandleTradeResponse)

	// Subscribe lending
	// for create/cancel order
	rabbitConn.SubscribeLendingOrders(lendingOrderService.HandleLendingOrdersCreateCancel)
	// for database changing response
	rabbitConn.SubscribeLendingOrderResponses(lendingOrderService.HandleLendingOrderResponse)
	rabbitConn.SubscribeLendingTradeResponses(lendingTradeService.HandleLendingTradeResponse)
	// start cron service
	cronService := crons.NewCronService(ohlcvService, priceBoardService, pairService, relayerService, eng)
	// initialize MongoDB Change Streams
	go orderService.WatchChanges()
	go tradeService.WatchChanges()

	// lending mongo watch change
	go lendingOrderService.WatchChanges()
	go lendingTradeService.WatchChanges()
	cronService.InitCrons()
	return r
}
