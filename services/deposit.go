package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/swap"
	"github.com/tomochain/backend-matching-engine/types"
)

// need to refractor using interface.SwappEngine and only expose neccessary methods
type DepositService struct {
	DepositDao interfaces.DepositDao
	SwapEngine *swap.Engine
}

// NewAddressService returns a new instance of accountService
func NewDepositService(DepositDAO interfaces.DepositDao, SwapEngine *swap.Engine) *DepositService {
	return &DepositService{DepositDAO, SwapEngine}
}

func (s *DepositService) GenerateAddress(chain types.Chain) (*common.Address, error) {
	err := s.DepositDao.IncrementAddressIndex(chain)
	if err != nil {
		return nil, err
	}
	index, err := s.DepositDao.GetAddressIndex(chain)
	if err != nil {
		return nil, err
	}

	return s.SwapEngine.EthereumAddressGenerator.Generate(index)
}

func (s *DepositService) GetSchemaVersion() uint64 {
	return s.DepositDao.GetSchemaVersion()
}
