package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/swap"
	"github.com/tomochain/backend-matching-engine/swap/ethereum"
	"github.com/tomochain/backend-matching-engine/swap/queue"
	"github.com/tomochain/backend-matching-engine/types"

	"github.com/tomochain/backend-matching-engine/errors"
)

// need to refractor using interface.SwappEngine and only expose neccessary methods
type DepositService struct {
	depositDao interfaces.DepositDao
	swapEngine *swap.Engine
	engine     interfaces.Engine
}

// NewAddressService returns a new instance of accountService
func NewDepositService(
	depositDao interfaces.DepositDao,
	swapEngine *swap.Engine,
	engine interfaces.Engine,
) *DepositService {

	depositService := &DepositService{depositDao, swapEngine, engine}

	// set event handler delegate to this service
	swapEngine.SetDelegate(depositService)
	// set storage engine to this service
	swapEngine.SetStorage(depositService)
	// run watching
	swapEngine.Start()

	return depositService
}

func (s *DepositService) GenerateAddress(chain types.Chain) (common.Address, error) {
	err := s.depositDao.IncrementAddressIndex(chain)
	if err != nil {
		return ethereum.EmptyAddress, err
	}
	index, err := s.depositDao.GetAddressIndex(chain)
	if err != nil {
		return ethereum.EmptyAddress, err
	}
	logger.Infof("Current index: %d", index)
	return s.swapEngine.EthereumAddressGenerator().Generate(index)
}

func (s *DepositService) SignerPublicKey() string {
	return s.swapEngine.SignerPublicKey()
}

func (s *DepositService) GetSchemaVersion() uint64 {
	return s.depositDao.GetSchemaVersion()
}

func (s *DepositService) RecoveryTransaction(chain types.Chain, address common.Address) error {
	return nil
}

/***** implement Storage interface ***/
func (s *DepositService) GetEthereumBlockToProcess() (uint64, error) {
	return s.depositDao.GetEthereumBlockToProcess()
}

func (s *DepositService) SaveLastProcessedEthereumBlock(block uint64) error {
	return s.depositDao.SaveLastProcessedEthereumBlock(block)
}

/***** events from engine ****/
// onNewEthereumTransaction checks if transaction is valid and adds it to
// the transactions queue for TomochainAccountConfigurator to consume.
//
// Transaction added to transactions queue should be in a format described in
// queue.Transaction (especialy amounts). Pooling service should not have to deal with any
// conversions.
func (s *DepositService) OnNewEthereumTransaction(transaction ethereum.Transaction) error {
	logger.Info("Processing transaction")

	// Let's check if tx is valid first.

	// Check if value is above minimum required
	if transaction.ValueWei.Cmp(s.swapEngine.MinimumValueWei()) < 0 {
		logger.Debug("Value is below minimum required amount, skipping")
		return nil
	}

	addressTo := common.HexToAddress(transaction.To)

	addressAssociation, err := s.GetAssociationByChainAddress(types.ChainEthereum, addressTo)
	if err != nil {
		return errors.Wrap(err, "Error getting association")
	}

	if addressAssociation == nil {
		logger.Debug("Associated address not found, skipping")
		return nil
	}

	// Add transaction as processing.
	processed, err := s.depositDao.AddProcessedTransaction(types.ChainEthereum, transaction.Hash, addressTo)
	if err != nil {
		return err
	}

	if processed {
		logger.Debug("Transaction already processed, skipping")
		return nil
	}

	// Add tx to the processing queue
	queueTx := queue.Transaction{
		TransactionID: transaction.Hash,
		AssetCode:     queue.AssetCodeETH,
		// Amount in the base unit of currency.
		Amount:             transaction.ValueWei.String(),
		TomochainPublicKey: addressAssociation.TomochainPublicKey.String(),
	}

	err = s.swapEngine.TransactionsQueue().QueueAdd(queueTx)
	if err != nil {
		return errors.Wrap(err, "Error adding transaction to the processing queue")
	}
	logger.Info("Transaction added to transaction queue")

	// Broadcast event to address stream
	logger.Infof("Broadcasting event: %v", transaction)
	logger.Info("Transaction processed successfully")
	return nil
}

func (s *DepositService) OnTomochainAccountCreated(destination string) {
	publicKey := common.HexToAddress(destination)
	association, err := s.depositDao.GetAssociationByTomochainPublicKey(publicKey)
	if err != nil {
		logger.Error("Error getting association")
		return
	}

	if association == nil {
		logger.Error("Association not found")
		return
	}
	// broast cast event association
	logger.Infof("Broasting event: %v", association)
}

func (s *DepositService) OnExchanged(destination string) {
	publicKey := common.HexToAddress(destination)
	association, err := s.depositDao.GetAssociationByTomochainPublicKey(publicKey)
	if err != nil {
		logger.Error("Error getting association")
		return
	}

	if association == nil {
		logger.Error("Association not found")
		return
	}

	logger.Infof("Broasting event: %v", association)
}

func (s *DepositService) OnExchangedTimelocked(destination, transaction string) {
	publicKey := common.HexToAddress(destination)
	association, err := s.depositDao.GetAssociationByTomochainPublicKey(publicKey)
	if err != nil {
		logger.Error("Error getting association")
		return
	}

	if association == nil {
		logger.Error("Association not found")
		return
	}

	// Save tx to database
	err = s.depositDao.AddRecoveryTransaction(publicKey, transaction)
	if err != nil {
		logger.Error("Error saving unlock transaction to DB")
		return
	}

	logger.Infof("Broasting event: %v", association)
}

// Create function performs the DB insertion task for Balance collection
func (s *DepositService) GetAssociationByChainAddress(chain types.Chain, userAddress common.Address) (*types.AddressAssociation, error) {
	// get from feed
	var addressAssociationFeed types.AddressAssociationFeed
	err := s.engine.GetFeed(userAddress, chain.Bytes(), &addressAssociationFeed)

	logger.Infof("feed :%v", addressAssociationFeed)

	if err == nil {
		return addressAssociationFeed.GetJSON()
	}
	return nil, err
}
