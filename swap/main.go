package swap

import (
	"math/big"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/gommon/log"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/swap/config"
	"github.com/tomochain/backend-matching-engine/swap/errors"
	"github.com/tomochain/backend-matching-engine/swap/ethereum"
	"github.com/tomochain/backend-matching-engine/swap/queue"
	"github.com/tomochain/backend-matching-engine/swap/tomochain"
	"github.com/tomochain/backend-matching-engine/utils"
)

// swap is engine
var logger = utils.EngineLogger

type Engine struct {
	Config                       *config.Config                 `inject:""`
	EthereumListener             *ethereum.Listener             `inject:""`
	EthereumAddressGenerator     *ethereum.AddressGenerator     `inject:""`
	TomochainAccountConfigurator *tomochain.AccountConfigurator `inject:""`
	TransactionsQueue            queue.Queue                    `inject:""`

	MinimumValueEth string
	SignerPublicKey string

	MinimumValueSat int64
	MinimumValueWei *big.Int
}

func NewEngine(cfg *config.Config) *Engine {
	engine := &Engine{
		SignerPublicKey: cfg.SignerPublicKey(),
	}
	if cfg.Ethereum != nil {
		ethereumListener := &ethereum.Listener{}
		ethereumClient, err := ethclient.Dial("http://" + cfg.Ethereum.RpcServer)
		if err != nil {
			logger.Error("Error connecting to geth")
			os.Exit(-1)
		}

		// config ethereum listener
		ethereumListener.Enabled = true
		ethereumListener.NetworkID = cfg.Ethereum.NetworkID
		ethereumListener.Client = ethereumClient

		engine.MinimumValueEth = cfg.Ethereum.MinimumValueEth

		ethereumAddressGenerator, err := ethereum.NewAddressGenerator(cfg.Ethereum.MasterPublicKey)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}

		engine.EthereumAddressGenerator = ethereumAddressGenerator
		engine.EthereumListener = ethereumListener
	}

	if cfg.Tomochain != nil {

		tomochainAccountConfigurator := &tomochain.AccountConfigurator{
			IssuerPublicKey:       cfg.Tomochain.IssuerPublicKey,
			DistributionPublicKey: cfg.Tomochain.DistributionPublicKey,
			SignerPrivateKey:      cfg.Tomochain.SignerPrivateKey,
			TokenAssetCode:        cfg.Tomochain.TokenAssetCode,
			StartingBalance:       cfg.Tomochain.StartingBalance,
			LockUnixTimestamp:     cfg.Tomochain.LockUnixTimestamp,
		}

		if cfg.Tomochain.StartingBalance == "" {
			tomochainAccountConfigurator.StartingBalance = "2.1"
		}

		if cfg.Ethereum != nil {
			tomochainAccountConfigurator.TokenPriceETH = cfg.Ethereum.TokenPrice
		}

		engine.TomochainAccountConfigurator = tomochainAccountConfigurator
	}

	engine.Config = cfg
	return engine
}

func (engine *Engine) SetDelegate(handler interfaces.SwapEngineHandler) {
	// delegate some handlers
	engine.EthereumListener.TransactionHandler = handler.OnNewEthereumTransaction
	engine.TomochainAccountConfigurator.OnAccountCreated = handler.OnTomochainAccountCreated
	engine.TomochainAccountConfigurator.OnExchanged = handler.OnExchanged
	engine.TomochainAccountConfigurator.OnExchangedTimelocked = handler.OnExchangedTimelocked
}

func (engine *Engine) Start() error {

	var err error
	engine.MinimumValueWei, err = ethereum.EthToWei(engine.MinimumValueEth)
	if err != nil {
		return errors.Wrap(err, "Invalid minimum accepted Ethereum transaction value")
	}

	if engine.MinimumValueWei.Cmp(new(big.Int)) == 0 {
		return errors.New("Minimum accepted Ethereum transaction value must be larger than 0")
	}

	err = engine.EthereumListener.Start(engine.Config.Ethereum.RpcServer)
	if err != nil {
		return errors.Wrap(err, "Error starting EthereumListener")
	}

	err = engine.TomochainAccountConfigurator.Start()
	if err != nil {
		return errors.Wrap(err, "Error starting TomochainAccountConfigurator")
	}

	// client will update swarm feed association so that we do not have to build broadcast engine

	signalInterrupt := make(chan os.Signal, 1)
	signal.Notify(signalInterrupt, os.Interrupt)

	go engine.poolTransactionsQueue()

	<-signalInterrupt
	engine.shutdown()

	return nil
}

// poolTransactionsQueue pools transactions queue which contains only processed and
// validated transactions and sends it to TomochainAccountConfigurator for account configuration.
func (s *Engine) poolTransactionsQueue() {
	logger.Infof("Started pooling transactions queue")

	for {
		transaction, err := s.TransactionsQueue.QueuePool()
		if err != nil {
			logger.Infof("Error pooling transactions queue")
			time.Sleep(time.Second)
			continue
		}

		if transaction == nil {
			time.Sleep(time.Second)
			continue
		}

		logger.Infof("Received transaction from transactions queue: %v", transaction)
		go s.TomochainAccountConfigurator.ConfigureAccount(
			transaction.TomochainPublicKey,
			string(transaction.AssetCode),
			transaction.Amount,
		)
	}
}

func (e *Engine) shutdown() {
	// do something
}

type GenerateAddressResponse struct {
	ProtocolVersion int    `json:"protocol_version"`
	Chain           string `json:"chain"`
	Address         string `json:"address"`
	Signer          string `json:"signer"`
}
