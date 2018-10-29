package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tomochain/backend-matching-engine/app"
	"github.com/tomochain/backend-matching-engine/contracts"
	"github.com/tomochain/backend-matching-engine/crons"
	"github.com/tomochain/backend-matching-engine/daos"
	"github.com/tomochain/backend-matching-engine/endpoints"
	"github.com/tomochain/backend-matching-engine/engine"
	"github.com/tomochain/backend-matching-engine/ethereum"
	"github.com/tomochain/backend-matching-engine/operator"
	"github.com/tomochain/backend-matching-engine/rabbitmq"
	"github.com/tomochain/backend-matching-engine/redis"
	"github.com/tomochain/backend-matching-engine/services"
	"github.com/tomochain/backend-matching-engine/ws"

	// "github.com/Proofsuite/go-ethereum/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Get application up and running",
	Long:  `Get application up and running`,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func run(cmd *cobra.Command, args []string) {

	// connect to the database
	_, err := daos.InitSession(nil)
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	log.Info("server %v is started at %v\n", app.Version, address)

	rabbitConn := rabbitmq.InitConnection(app.Config.Rabbitmq)
	redisConn := redis.NewRedisConnection(app.Config.Redis)
	provider := ethereum.NewWebsocketProvider()

	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r.Use(func(c *gin.Context) {
	// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// 	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(204)
	// 		return
	// 	}

	// 	c.Next()
	// })
	err = RegisterRouter(r, provider, redisConn, rabbitConn)
	if err != nil {
		panic(err)
	}
	r.GET("/socket", ws.ConnectionEndpoint)
	// // start the server
	panic(r.Run(address)) // listen and serve on 0.0.0.0:8080
}

func RegisterRouter(
	r *gin.Engine,
	provider *ethereum.EthereumProvider,
	redisConn *redis.RedisConnection,
	rabbitConn *rabbitmq.Connection,
) error {

	// get daos for dependency injection
	orderDao := daos.NewOrderDao()
	tokenDao := daos.NewTokenDao()
	pairDao := daos.NewPairDao()
	tradeDao := daos.NewTradeDao()
	accountDao := daos.NewAccountDao()
	walletDao := daos.NewWalletDao()

	// instantiate engine
	eng := engine.NewEngine(redisConn, rabbitConn, pairDao)

	// get services for injection
	accountService := services.NewAccountService(accountDao, tokenDao)
	ohlcvService := services.NewOHLCVService(tradeDao)
	tokenService := services.NewTokenService(tokenDao)
	tradeService := services.NewTradeService(tradeDao)
	pairService := services.NewPairService(pairDao, tokenDao, eng, tradeService)
	orderService := services.NewOrderService(orderDao, pairDao, accountDao, tradeDao, eng, provider, rabbitConn)
	orderBookService := services.NewOrderBookService(pairDao, tokenDao, orderDao, eng)
	walletService := services.NewWalletService(walletDao)
	cronService := crons.NewCronService(ohlcvService)

	var err error
	// fmt.Printf("config %v", app.Config)
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

	// deploy http and ws endpoints
	endpoints.ServeAccountResource(r, accountService)
	endpoints.ServeTokenResource(r, tokenService)
	endpoints.ServePairResource(r, pairService)
	endpoints.ServeOrderBookResource(r, orderBookService)
	endpoints.ServeOHLCVResource(r, ohlcvService)
	endpoints.ServeTradeResource(r, tradeService)
	endpoints.ServeOrderResource(r, orderService, eng)

	// deploy operator
	op, err := operator.NewOperator(
		walletService,
		tradeService,
		orderService,
		provider,
		exchange,
		rabbitConn,
	)

	if err == nil {

		// instead of using rabbit mq for working with channel from smart contract event logs
		// we use pss to subscribe directly to decentralized message queue

		//initialize rabbitmq subscriptions
		rabbitConn.SubscribeOrders(eng.HandleOrders)
		rabbitConn.SubscribeTrades(op.HandleTrades)
		rabbitConn.SubscribeOperator(orderService.HandleOperatorMessages)
		rabbitConn.SubscribeEngineResponses(orderService.HandleEngineResponse)

		// this service is for crawling data if the backend is offline for a while and need a catch-up
		cronService.InitCrons()
	}

	return err
}
