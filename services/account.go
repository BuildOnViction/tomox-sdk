package services

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/dex-server/interfaces"
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/utils/math"
	"gopkg.in/mgo.v2/bson"
)

type AccountService struct {
	accountDao interfaces.AccountDao
	tokenDao   interfaces.TokenDao
}

// NewAddressService returns a new instance of accountService
func NewAccountService(
	accountDao interfaces.AccountDao,
	tokenDao interfaces.TokenDao,
) *AccountService {
	return &AccountService{accountDao, tokenDao}
}

func (s *AccountService) Create(a *types.Account) error {
	addr := a.Address

	acc, err := s.accountDao.GetByAddress(addr)
	if err != nil {
		logger.Error(err)
		return err
	}

	if acc != nil {
		return ErrAccountExists
	}

	tokens, err := s.tokenDao.GetAll()
	if err != nil {
		logger.Error(err)
		return err
	}

	a.IsBlocked = false
	a.TokenBalances = make(map[common.Address]*types.TokenBalance)

	ten := big.NewInt(10)

	// currently by default, the tokens balances are set to 0
	for _, token := range tokens {
		fmt.Println(int64(token.Decimals))
		decimals := big.NewInt(int64(token.Decimals))
		a.TokenBalances[token.ContractAddress] = &types.TokenBalance{
			Address:        token.ContractAddress,
			Symbol:         token.Symbol,
			Balance:        math.Mul(big.NewInt(types.DefaultTestBalance()), ten.Exp(ten, decimals, nil)),
			LockedBalance:  big.NewInt(types.DefaultTestLockedBalance()),
			PendingBalance: big.NewInt(types.DefaultTestPendingBalance()),
		}
	}

	nativeCurrency := types.GetNativeCurrency()

	a.TokenBalances[nativeCurrency.Address] = &types.TokenBalance{
		Address:        nativeCurrency.Address,
		Symbol:         nativeCurrency.Symbol,
		Balance:        math.Mul(big.NewInt(types.DefaultTestBalance()), ten.Exp(ten, big.NewInt(int64(nativeCurrency.Decimals)), nil)),
		LockedBalance:  big.NewInt(types.DefaultTestLockedBalance()),
		PendingBalance: big.NewInt(types.DefaultTestPendingBalance()),
	}

	if a != nil {
		err = s.accountDao.Create(a)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (s *AccountService) FindOrCreate(addr common.Address) (*types.Account, error) {
	a, err := s.accountDao.GetByAddress(addr)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if a != nil {
		return a, nil
	}

	tokens, err := s.tokenDao.GetAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	a = &types.Account{
		Address:       addr,
		IsBlocked:     false,
		TokenBalances: make(map[common.Address]*types.TokenBalance),
	}

	ten := big.NewInt(10)

	// currently by default, the tokens balances are set to 0
	for _, t := range tokens {
		decimals := big.NewInt(int64(t.Decimals))
		a.TokenBalances[t.ContractAddress] = &types.TokenBalance{
			Address:        t.ContractAddress,
			Symbol:         t.Symbol,
			Balance:        math.Mul(big.NewInt(types.DefaultTestBalance()), ten.Exp(ten, decimals, nil)),
			LockedBalance:  big.NewInt(types.DefaultTestLockedBalance()),
			PendingBalance: big.NewInt(types.DefaultTestPendingBalance()),
		}
	}

	nativeCurrency := types.GetNativeCurrency()

	a.TokenBalances[nativeCurrency.Address] = &types.TokenBalance{
		Address:        nativeCurrency.Address,
		Symbol:         nativeCurrency.Symbol,
		Balance:        math.Mul(big.NewInt(types.DefaultTestBalance()), ten.Exp(ten, big.NewInt(int64(nativeCurrency.Decimals)), nil)),
		LockedBalance:  big.NewInt(types.DefaultTestLockedBalance()),
		PendingBalance: big.NewInt(types.DefaultTestPendingBalance()),
	}

	err = s.accountDao.Create(a)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return a, nil
}

func (s *AccountService) GetByID(id bson.ObjectId) (*types.Account, error) {
	return s.accountDao.GetByID(id)
}

func (s *AccountService) GetAll() ([]types.Account, error) {
	return s.accountDao.GetAll()
}

func (s *AccountService) GetByAddress(a common.Address) (*types.Account, error) {
	return s.accountDao.GetByAddress(a)
}

func (s *AccountService) GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error) {
	return s.accountDao.GetTokenBalance(owner, token)
}

func (s *AccountService) GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error) {
	return s.accountDao.GetTokenBalances(owner)
}

func (s *AccountService) Transfer(token common.Address, fromAddress common.Address, toAddress common.Address, amount *big.Int) error {
	return s.accountDao.Transfer(token, fromAddress, toAddress, amount)
}
