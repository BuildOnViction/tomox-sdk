package swap

import (
	"math/big"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/stellar/go/support/log"
	"github.com/tomochain/backend-matching-engine/swap/config"
	"github.com/tomochain/backend-matching-engine/swap/ethereum"
	"github.com/tomochain/backend-matching-engine/swap/queue"
	"github.com/tomochain/backend-matching-engine/swap/tomochain"
)

type Engine struct {
	// Config                       *config.Config                 `inject:""`
	EthereumListener             *ethereum.Listener             `inject:""`
	EthereumAddressGenerator     *ethereum.AddressGenerator     `inject:""`
	TomochainAccountConfigurator *tomochain.AccountConfigurator `inject:""`
	TransactionsQueue            queue.Queue                    `inject:""`

	MinimumValueEth string
	SignerPublicKey string

	minimumValueSat int64
	minimumValueWei *big.Int
}

func NewEngine(cfg *config.Config) *Engine {
	engine := &Engine{
		SignerPublicKey: cfg.SignerPublicKey(),
	}

	ethereumClient := &ethclient.Client{}
	ethereumListener := &ethereum.Listener{}

	ethereumClient, err = ethclient.Dial("http://" + cfg.Ethereum.RpcServer)
	if err != nil {
		log.WithField("err", err).Error("Error connecting to geth")
		os.Exit(-1)
	}

	ethereumListener.Enabled = true
	ethereumListener.NetworkID = cfg.Ethereum.NetworkID
	server.MinimumValueEth = cfg.Ethereum.MinimumValueEth

	ethereumAddressGenerator, err := ethereum.NewAddressGenerator(cfg.Ethereum.MasterPublicKey)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}

	tomochainAccountConfigurator := &tomochain.AccountConfigurator{
		NetworkPassphrase:     cfg.Tomochain.NetworkPassphrase,
		IssuerPublicKey:       cfg.Tomochain.IssuerPublicKey,
		DistributionPublicKey: cfg.Tomochain.DistributionPublicKey,
		SignerSecretKey:       cfg.Tomochain.SignerSecretKey,
		NeedsAuthorize:        cfg.Tomochain.NeedsAuthorize,
		TokenAssetCode:        cfg.Tomochain.TokenAssetCode,
		StartingBalance:       cfg.Tomochain.StartingBalance,
		LockUnixTimestamp:     cfg.Tomochain.LockUnixTimestamp,
	}

	if cfg.Tomochain.StartingBalance == "" {
		tomochainAccountConfigurator.StartingBalance = "2.1"
	}

	if cfg.Bitcoin != nil {
		tomochainAccountConfigurator.TokenPriceBTC = cfg.Bitcoin.TokenPrice
	}

	if cfg.Ethereum != nil {
		tomochainAccountConfigurator.TokenPriceETH = cfg.Ethereum.TokenPrice
	}

	engine.TomochainAccountConfigurator = tomochainAccountConfigurator
	engine.EthereumAddressGenerator = ethereumAddressGenerator
	engine.EthereumListener = ethereumListener
}

func (e *Engine) Start() error {
	e.EthereumListener.TransactionHandler = e.onNewEthereumTransaction
	e.TomochainAccountConfigurator.OnAccountCreated = e.onTomochainAccountCreated
	e.TomochainAccountConfigurator.OnExchanged = e.onExchanged
	e.TomochainAccountConfigurator.OnExchangedTimelocked = e.OnExchangedTimelocked

	var err error
	e.minimumValueWei, err = ethereum.EthToWei(e.MinimumValueEth)
	if err != nil {
		return errors.Wrap(err, "Invalid minimum accepted Ethereum transaction value")
	}

	if e.minimumValueWei.Cmp(new(big.Int)) == 0 {
		return errors.New("Minimum accepted Ethereum transaction value must be larger than 0")
	}

	err = e.EthereumListener.Start(e.Config.Ethereum.RpcServer)
	if err != nil {
		return errors.Wrap(err, "Error starting EthereumListener")
	}

	err := e.TomochainAccountConfigurator.Start()
	if err != nil {
		return errors.Wrap(err, "Error starting TomochainAccountConfigurator")
	}

	// client will update swarm feed association so that we do not have to build broadcast engine

	signalInterrupt := make(chan os.Signal, 1)
	signal.Notify(signalInterrupt, os.Interrupt)

	go e.poolTransactionsQueue()

	<-signalInterrupt
	s.shutdown()

	return nil
}

// poolTransactionsQueue pools transactions queue which contains only processed and
// validated transactions and sends it to TomochainAccountConfigurator for account configuration.
func (s *Engine) poolTransactionsQueue() {
	logger.Infof("Started pooling transactions queue")

	for {
		transaction, err := s.TransactionsQueue.QueuePool()
		if err != nil {
			logger.Logf("Error pooling transactions queue")
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
