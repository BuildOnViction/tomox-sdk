package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomoxsdk/app"
	"github.com/tomochain/tomoxsdk/contracts"
	"github.com/tomochain/tomoxsdk/crons"
	"github.com/tomochain/tomoxsdk/daos"
	"github.com/tomochain/tomoxsdk/endpoints"
	"github.com/tomochain/tomoxsdk/engine"
	"github.com/tomochain/tomoxsdk/errors"
	"github.com/tomochain/tomoxsdk/ethereum"
	"github.com/tomochain/tomoxsdk/operator"
	"github.com/tomochain/tomoxsdk/rabbitmq"
	"github.com/tomochain/tomoxsdk/services"
	"github.com/tomochain/tomoxsdk/swap"
	"github.com/tomochain/tomoxsdk/types"
	"github.com/tomochain/tomoxsdk/ws"
)

const (
	swaggerUIDir = "/swaggerui/"
)

func Start() {
	env := os.Getenv("GO_ENV")

	if err := app.LoadConfig("./config", env); err != nil {
		panic(err)
	}

	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(err)
	}

	// connect to the database
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
	configDao := daos.NewConfigDao()
	associationDao := daos.NewAssociationDao()
	fiatPriceDao := daos.NewFiatPriceDao()
	notificationDao := daos.NewNotificationDao()

	// instantiate engine
	eng := engine.NewEngine(rabbitConn, orderDao, tradeDao, pairDao, provider)
	swapEngine := swap.NewEngine(app.Config.Deposit)

	// get services for injection
	accountService := services.NewAccountService(accountDao, tokenDao)
	ohlcvService := services.NewOHLCVService(tradeDao)
	tokenService := services.NewTokenService(tokenDao)
	tradeService := services.NewTradeService(tradeDao)
	validatorService := services.NewValidatorService(provider, accountDao, orderDao, pairDao)
	pairService := services.NewPairService(pairDao, tokenDao, tradeDao, orderDao, eng, provider)
	orderService := services.NewOrderService(orderDao, pairDao, accountDao, tradeDao, notificationDao, eng, validatorService, rabbitConn)
	orderBookService := services.NewOrderBookService(pairDao, tokenDao, orderDao, eng)

	walletService := services.NewWalletService(walletDao)

	// txservice for deposit
	// wallet := &types.NewWalletFromPrivateKey(app.Config.Deposit.Tomochain.SignerPrivateKey)
	// we already have them so no need to re-calculate
	wallet := &types.Wallet{
		Address:    app.Config.Deposit.Tomochain.GetPublicKey(),
		PrivateKey: app.Config.Deposit.Tomochain.GetPrivateKey(),
	}
	txService := services.NewTxService(walletDao, wallet)
	depositService := services.NewDepositService(configDao, associationDao, pairDao, orderDao, swapEngine, eng, rabbitConn)
	priceBoardService := services.NewPriceBoardService(tokenDao, tradeDao)
	fiatPriceService := services.NewFiatPriceService(tokenDao, fiatPriceDao)
	marketsService := services.NewMarketsService(pairDao, orderDao, tradeDao, ohlcvService, fiatPriceService)
	notificationService := services.NewNotificationService(notificationDao)

	// start cron service
	cronService := crons.NewCronService(ohlcvService, priceBoardService, pairService, fiatPriceService)

	// get exchange contract instance
	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	exchange, err := contracts.NewExchange(
		walletService,
		exchangeAddress,
		provider.Client,
	)

	if err != nil {
		panic(err)
	}

	// deploy operator
	op, err := operator.NewOperator(
		walletService,
		tradeService,
		orderService,
		provider,
		exchange,
		rabbitConn,
		accountService,
		tokenService,
	)

	if err != nil {
		panic(err)
	}

	// deploy http and ws endpoints
	endpoints.ServeInfoResource(r, walletService, tokenService)
	endpoints.ServeAccountResource(r, accountService)
	endpoints.ServeTokenResource(r, tokenService)
	endpoints.ServePairResource(r, pairService)
	endpoints.ServeOrderBookResource(r, orderBookService)
	endpoints.ServeOHLCVResource(r, ohlcvService)
	endpoints.ServeTradeResource(r, tradeService)
	endpoints.ServeOrderResource(r, orderService, accountService)

	endpoints.ServeDepositResource(r, depositService, walletService, txService)
	endpoints.ServePriceBoardResource(r, priceBoardService)
	endpoints.ServeMarketsResource(r, marketsService)
	endpoints.ServeNotificationResource(r, notificationService)

	// Swagger UI
	sh := http.StripPrefix(swaggerUIDir, http.FileServer(http.Dir("."+swaggerUIDir)))
	r.PathPrefix(swaggerUIDir).Handler(sh)

	//initialize rabbitmq subscriptions
	rabbitConn.SubscribeOrders(eng.HandleOrders)
	rabbitConn.SubscribeEngineResponses(orderService.HandleEngineResponse)
	rabbitConn.SubscribeTrades(op.HandleTrades)
	rabbitConn.SubscribeOperator(orderService.HandleOperatorMessages)

	// Initialize fiat price
	fiatPriceService.InitFiatPrice()

	cronService.InitCrons()
	return r
}
