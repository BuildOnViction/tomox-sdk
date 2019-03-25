package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tomochain/dex-server/app"
	"github.com/tomochain/dex-server/contracts"
	"github.com/tomochain/dex-server/crons"
	"github.com/tomochain/dex-server/daos"
	"github.com/tomochain/dex-server/endpoints"
	"github.com/tomochain/dex-server/engine"
	"github.com/tomochain/dex-server/errors"
	"github.com/tomochain/dex-server/ethereum"
	"github.com/tomochain/dex-server/operator"
	"github.com/tomochain/dex-server/rabbitmq"
	"github.com/tomochain/dex-server/services"
	"github.com/tomochain/dex-server/swap"
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/ws"
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
	priceBoardDao := daos.NewPriceBoardDao()

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
	orderService := services.NewOrderService(orderDao, pairDao, accountDao, tradeDao, eng, validatorService, rabbitConn)
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
	priceBoardService := services.NewPriceBoardService(tokenDao, tradeDao, priceBoardDao)

	// start cron service
	cronService := crons.NewCronService(ohlcvService, priceBoardService, pairService, eng)

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

	//initialize rabbitmq subscriptions
	rabbitConn.SubscribeOrders(eng.HandleOrders)
	rabbitConn.SubscribeEngineResponses(orderService.HandleEngineResponse)
	rabbitConn.SubscribeTrades(op.HandleTrades)
	rabbitConn.SubscribeOperator(orderService.HandleOperatorMessages)

	// initialize MongoDB Change Streams
	orderService.WatchChanges()

	cronService.InitCrons()
	return r
}
