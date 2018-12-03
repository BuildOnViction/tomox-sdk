package services

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/backend-matching-engine/ethereum"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/rabbitmq"
	"github.com/tomochain/backend-matching-engine/swap"
	swapEthereum "github.com/tomochain/backend-matching-engine/swap/ethereum"
	"github.com/tomochain/backend-matching-engine/types"
	"gopkg.in/mgo.v2/bson"
)

// need to refractor using interface.SwappEngine and only expose neccessary methods
type DepositService struct {
	configDao      interfaces.ConfigDao
	associationDao interfaces.AssociationDao
	swapEngine     *swap.Engine
	engine         interfaces.Engine
	broker         *rabbitmq.Connection
}

// NewAddressService returns a new instance of accountService
func NewDepositService(
	configDao interfaces.ConfigDao,
	associationDao interfaces.AssociationDao,
	swapEngine *swap.Engine,
	engine interfaces.Engine,
	broker *rabbitmq.Connection,
) *DepositService {

	depositService := &DepositService{configDao, associationDao, swapEngine, engine, broker}

	// set storage engine to this service
	swapEngine.SetStorage(depositService)

	swapEngine.SetQueue(depositService)

	// run watching
	swapEngine.Start()

	return depositService
}

func (s *DepositService) EthereumClient() interfaces.EthereumClient {
	provider := s.engine.Provider().(*ethereum.EthereumProvider)
	return provider.Client
}

func (s *DepositService) WethAddress() common.Address {
	provider := s.engine.Provider().(*ethereum.EthereumProvider)
	return provider.Config.WethAddress()
}

func (s *DepositService) SetDelegate(handler interfaces.SwapEngineHandler) {
	// set event handler delegate to this service
	s.swapEngine.SetDelegate(handler)
}

func (s *DepositService) GenerateAddress(chain types.Chain) (common.Address, error) {
	err := s.configDao.IncrementAddressIndex(chain)
	if err != nil {
		return swapEthereum.EmptyAddress, err
	}
	index, err := s.configDao.GetAddressIndex(chain)
	if err != nil {
		return swapEthereum.EmptyAddress, err
	}
	logger.Infof("Current index: %d", index)
	return s.swapEngine.EthereumAddressGenerator().Generate(index)
}

func (s *DepositService) SignerPublicKey() common.Address {
	return s.swapEngine.SignerPublicKey()
}

func (s *DepositService) GetSchemaVersion() uint64 {
	return s.configDao.GetSchemaVersion()
}

func (s *DepositService) RecoveryTransaction(chain types.Chain, address common.Address) error {
	return nil
}

/***** implement Storage interface ***/
func (s *DepositService) GetEthereumBlockToProcess() (uint64, error) {
	return s.configDao.GetEthereumBlockToProcess()
}

func (s *DepositService) SaveLastProcessedEthereumBlock(block uint64) error {
	return s.configDao.SaveLastProcessedEthereumBlock(block)
}

func (s *DepositService) SaveDepositTransaction(chain types.Chain, sourceAccount common.Address, txEnvelope string) error {
	return s.associationDao.SaveDepositTransaction(chain, sourceAccount, txEnvelope)
}

func (s *DepositService) QueueAdd(transaction *types.DepositTransaction) error {
	err := s.broker.PublishDepositTransaction(transaction)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// QueuePool receives and removes the head of this queue. Returns nil if no elements found.
func (s *DepositService) QueuePool() (<-chan *types.DepositTransaction, error) {
	return s.broker.QueuePoolDepositTransactions()
}

func (s *DepositService) MinimumValueWei() *big.Int {
	return s.swapEngine.MinimumValueWei()
}

func (s *DepositService) GetAssociationByChainAddress(chain types.Chain, userAddress common.Address) (*types.AddressAssociationRecord, error) {
	return s.associationDao.GetAssociationByChainAddress(chain, userAddress)
}

func (s *DepositService) SaveAssociationByChainAddress(chain types.Chain, address, associatedAddress common.Address) error {
	association := &types.AddressAssociationRecord{
		ID:                bson.NewObjectId(),
		Chain:             chain.String(),
		Address:           address.Hex(),
		AssociatedAddress: associatedAddress.Hex(),
	}

	return s.associationDao.SaveAssociation(association)
}

// Create function performs the DB insertion task for Balance collection
func (s *DepositService) GetAssociationByTomochainPublicKey(chain types.Chain, userAddress common.Address) (*types.AddressAssociation, error) {
	// get from feed
	var addressAssociationFeed types.AddressAssociationFeed
	err := s.engine.GetFeed(userAddress, chain.Bytes(), &addressAssociationFeed)

	logger.Infof("feed :%v", addressAssociationFeed)

	if err == nil {
		return addressAssociationFeed.GetJSON()
	}
	return nil, err
}
