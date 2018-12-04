package services

import (
	"math/big"
	"strings"

	"github.com/tomochain/backend-matching-engine/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/backend-matching-engine/ethereum"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/rabbitmq"
	"github.com/tomochain/backend-matching-engine/swap"
	swapEthereum "github.com/tomochain/backend-matching-engine/swap/ethereum"
	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/utils/math"
	"gopkg.in/mgo.v2/bson"
)

// need to refractor using interface.SwappEngine and only expose neccessary methods
type DepositService struct {
	configDao      interfaces.ConfigDao
	associationDao interfaces.AssociationDao
	pairDao        interfaces.PairDao
	orderDao       interfaces.OrderDao
	swapEngine     *swap.Engine
	engine         interfaces.Engine
	broker         *rabbitmq.Connection
}

// NewAddressService returns a new instance of accountService
func NewDepositService(
	configDao interfaces.ConfigDao,
	associationDao interfaces.AssociationDao,
	pairDao interfaces.PairDao,
	orderDao interfaces.OrderDao,
	swapEngine *swap.Engine,
	engine interfaces.Engine,
	broker *rabbitmq.Connection,
) *DepositService {

	depositService := &DepositService{configDao, associationDao, pairDao, orderDao, swapEngine, engine, broker}

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

func (s *DepositService) GenerateAddress(chain types.Chain) (common.Address, uint64, error) {

	err := s.configDao.IncrementAddressIndex(chain)
	if err != nil {
		return swapEthereum.EmptyAddress, 0, err
	}
	index, err := s.configDao.GetAddressIndex(chain)
	if err != nil {
		return swapEthereum.EmptyAddress, 0, err
	}
	logger.Infof("Current index: %d", index)
	address, err := s.swapEngine.EthereumAddressGenerator().Generate(index)
	return address, index, err
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

func (s *DepositService) GetAssociationByChainAssociatedAddress(chain types.Chain, associatedAddress common.Address) (*types.AddressAssociationRecord, error) {
	return s.associationDao.GetAssociationByChainAssociatedAddress(chain, associatedAddress)
}

func (s *DepositService) SaveAssociationByChainAddress(chain types.Chain, address common.Address, index uint64, associatedAddress common.Address, pairAddresses *types.PairAddresses) error {

	association := &types.AddressAssociationRecord{
		ID:                bson.NewObjectId(),
		Chain:             chain,
		Address:           address.Hex(),
		AddressIndex:      index,
		Status:            types.PENDING,
		AssociatedAddress: associatedAddress.Hex(),
		PairName:          pairAddresses.Name,
		BaseTokenAddress:  pairAddresses.BaseToken.Hex(),
		QuoteTokenAddress: pairAddresses.QuoteToken.Hex(),
	}

	return s.associationDao.SaveAssociation(association)
}

func (s *DepositService) SaveAssociationStatusByChainAddress(chain types.Chain, address common.Address, status string) error {
	return s.associationDao.SaveAssociationStatus(chain, address, status)
}

func (s *DepositService) getTokenAmountFromOracle(baseTokenSymbol, quoteTokenSymbol string, quoteAmount *big.Int) (*big.Int, error) {
	return quoteAmount, nil
}

func (s *DepositService) GetBaseTokenAmount(pairName string, quoteAmount *big.Int) (*big.Int, error) {

	tokenSymbols := strings.Split(pairName, "/")
	if len(tokenSymbols) != 2 {
		return nil, errors.Errorf("Pair name is wrong format: %s", pairName)
	}
	baseTokenSymbol := tokenSymbols[0]
	quoteTokenSymbol := tokenSymbols[1]

	// this is 1:1 exchange
	if baseTokenSymbol == quoteTokenSymbol {
		return quoteAmount, nil
	}

	pair, err := s.pairDao.GetByTokenSymbols(baseTokenSymbol, quoteTokenSymbol)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		// there is no exchange rate yet
		return s.getTokenAmountFromOracle(baseTokenSymbol, quoteTokenSymbol, quoteAmount)
	}

	logger.Debugf("Got pair :%v", pair)

	// get best Bid, the highest bid available
	bids, err := s.orderDao.GetSideOrderBook(pair, types.BUY, -1, 1)
	if err != nil {
		return nil, err
	}

	// if there is no exchange rate, should return one from oracle service like coin market cap
	if len(bids) < 1 {
		return s.getTokenAmountFromOracle(baseTokenSymbol, quoteTokenSymbol, quoteAmount)
	}

	pricepoint := new(big.Int)
	pricepoint.SetString(bids[0]["pricepoint"], 10)

	tokenAmount := new(big.Int)
	tokenAmount = math.Div(quoteAmount, pricepoint)
	tokenAmount = math.Mul(tokenAmount, pair.PriceMultiplier)

	return tokenAmount, nil
}

// Create function performs the DB insertion task for Balance collection
func (s *DepositService) GetAssociationByUserAddress(chain types.Chain, userAddress common.Address) (*types.AddressAssociation, error) {
	// get from feed
	var addressAssociationFeed types.AddressAssociationFeed
	err := s.engine.GetFeed(userAddress, chain.Bytes(), &addressAssociationFeed)

	logger.Infof("feed :%v", addressAssociationFeed)

	if err == nil {
		return addressAssociationFeed.GetJSON()
	}
	return nil, err
}
