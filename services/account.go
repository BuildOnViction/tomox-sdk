package services

import (
	math2 "math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/math"
)

type AccountService struct {
	AccountDao   interfaces.AccountDao
	TokenDao     interfaces.TokenDao
	PairDao      interfaces.PairDao
	OrderDao     interfaces.OrderDao
	LendingDao   interfaces.LendingOrderDao
	Provider     interfaces.EthereumProvider
	OHLCVService interfaces.OHLCVService
}

// NewAccountService returns a new instance of accountService
func NewAccountService(
	accountDao interfaces.AccountDao,
	tokenDao interfaces.TokenDao,
	pairDao interfaces.PairDao,
	orderDao interfaces.OrderDao,
	lendingDao interfaces.LendingOrderDao,
	provider interfaces.EthereumProvider,
	ohlcvService interfaces.OHLCVService,
) *AccountService {
	return &AccountService{
		AccountDao:   accountDao,
		TokenDao:     tokenDao,
		PairDao:      pairDao,
		OrderDao:     orderDao,
		LendingDao:   lendingDao,
		Provider:     provider,
		OHLCVService: ohlcvService,
	}
}

func (s *AccountService) Create(a *types.Account) error {
	addr := a.Address

	acc, err := s.AccountDao.GetByAddress(addr)
	if err != nil {
		logger.Error(err)
		return err
	}

	if acc != nil {
		return ErrAccountExists
	}

	tokens, err := s.TokenDao.GetAll()
	if err != nil {
		logger.Error(err)
		return err
	}

	a.IsBlocked = false
	a.TokenBalances = make(map[common.Address]*types.TokenBalance)

	ten := big.NewInt(10)

	// currently by default, the tokens balances are set to 0
	for _, token := range tokens {
		decimals := big.NewInt(int64(token.Decimals))
		a.TokenBalances[token.ContractAddress] = &types.TokenBalance{
			Address:          token.ContractAddress,
			Symbol:           token.Symbol,
			Decimals:         token.Decimals,
			Balance:          math.Mul(big.NewInt(types.DefaultTestBalance()), math.Exp(ten, decimals)),
			InOrderBalance:   big.NewInt(types.DefaultTestInOrderBalance()),
			AvailableBalance: big.NewInt(types.DefaultTestAvailableBalance()),
		}
	}

	nativeCurrency := types.GetNativeCurrency()

	a.TokenBalances[nativeCurrency.Address] = &types.TokenBalance{
		Address:          nativeCurrency.Address,
		Symbol:           nativeCurrency.Symbol,
		Balance:          math.Mul(big.NewInt(types.DefaultTestBalance()), math.Exp(ten, big.NewInt(int64(nativeCurrency.Decimals)))),
		InOrderBalance:   big.NewInt(types.DefaultTestInOrderBalance()),
		AvailableBalance: big.NewInt(types.DefaultTestAvailableBalance()),
	}

	if a != nil {
		err = s.AccountDao.Create(a)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

//FindOrCreate find or create if not found [Deprecated]
func (s *AccountService) FindOrCreate(addr common.Address) (*types.Account, error) {
	a, err := s.AccountDao.GetByAddress(addr)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if a != nil {
		return a, nil
	}

	tokens, err := s.TokenDao.GetAll()
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
			Address:          t.ContractAddress,
			Symbol:           t.Symbol,
			Balance:          math.Mul(big.NewInt(types.DefaultTestBalance()), math.Exp(ten, decimals)),
			InOrderBalance:   big.NewInt(types.DefaultTestInOrderBalance()),
			AvailableBalance: big.NewInt(types.DefaultTestAvailableBalance()),
		}
	}

	nativeCurrency := types.GetNativeCurrency()

	a.TokenBalances[nativeCurrency.Address] = &types.TokenBalance{
		Address:          nativeCurrency.Address,
		Symbol:           nativeCurrency.Symbol,
		Balance:          math.Mul(big.NewInt(types.DefaultTestBalance()), math.Exp(ten, big.NewInt(int64(nativeCurrency.Decimals)))),
		InOrderBalance:   big.NewInt(types.DefaultTestInOrderBalance()),
		AvailableBalance: big.NewInt(types.DefaultTestAvailableBalance()),
	}

	err = s.AccountDao.Create(a)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return a, nil
}

// GetByID get account by id [Deprecated]
func (s *AccountService) GetByID(id bson.ObjectId) (*types.Account, error) {
	return s.AccountDao.GetByID(id)
}

// GetAll get all account [Deprecated]
func (s *AccountService) GetAll() ([]types.Account, error) {
	return s.AccountDao.GetAll()
}

// GetByAddress get account from address
func (s *AccountService) GetByAddress(a common.Address) (*types.Account, error) {
	account, err := s.AccountDao.GetByAddress(a)
	if err != nil {
		return nil, err
	}
	if account == nil {
		account = &types.Account{
			Address:       a,
			TokenBalances: make(map[common.Address]*types.TokenBalance),
			IsBlocked:     false,
		}
	}
	tokens, err := s.TokenDao.GetAll()
	if err != nil || tokens == nil {
		return nil, err
	}
	for _, token := range tokens {
		balance, err := s.GetTokenBalanceProvidor(a, token.ContractAddress)
		if err != nil {
			return nil, err
		}

		price, _ := s.OHLCVService.GetLastPriceCurrentByTime(balance.Symbol, time.Now())

		if balance != nil && price != nil {
			inUsdBalance := new(big.Float).Mul(price, new(big.Float).SetInt(balance.Balance))
			inUsdBalance = new(big.Float).Quo(inUsdBalance, new(big.Float).SetInt(big.NewInt(int64(math2.Pow10(balance.Decimals)))))
			balance.InUsdBalance = inUsdBalance
		}

		account.TokenBalances[token.ContractAddress] = balance
	}

	return account, nil
}

// GetTokenBalance database [Deprecated]
func (s *AccountService) GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error) {
	return s.AccountDao.GetTokenBalance(owner, token)
}

// GetTokenBalanceProvidor get balance from chain
func (s *AccountService) GetTokenBalanceProvidor(owner common.Address, tokenAddress common.Address) (*types.TokenBalance, error) {
	token, err := s.TokenDao.GetByAddress(tokenAddress)
	if err != nil || token == nil {
		return nil, err
	}
	tokenBalance := &types.TokenBalance{
		Address:          tokenAddress,
		Symbol:           token.Symbol,
		Decimals:         token.Decimals,
		AvailableBalance: big.NewInt(0),
		InOrderBalance:   big.NewInt(0),
		Balance:          big.NewInt(0),
		InUsdBalance:     big.NewFloat(0),
	}
	b, err := s.Provider.Balance(owner, tokenAddress)
	if err != nil {
		return nil, err
	}
	tokenBalance.Balance = b

	listPairs, err := s.PairDao.GetActivePairs()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	tokens, err := s.TokenDao.GetAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	sellTokenExchangeLockedBalance, err := s.OrderDao.GetUserLockedBalance(owner, tokenAddress, listPairs)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	sellTokenLendingLockedBalance, err := s.LendingDao.GetUserLockedBalance(owner, tokenAddress, tokens)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	sellTokenLockedBalance := new(big.Int).Add(sellTokenExchangeLockedBalance, sellTokenLendingLockedBalance)
	tokenBalance.InOrderBalance = sellTokenLockedBalance
	tokenBalance.AvailableBalance = math.Sub(b, sellTokenLockedBalance)
	return tokenBalance, nil
}

// GetTokenBalancesProvidor get balances from chain
func (s *AccountService) GetTokenBalancesProvidor(owner common.Address) (map[common.Address]*types.TokenBalance, error) {
	return s.AccountDao.GetTokenBalances(owner)
}

// GetTokenBalances Deprecated
func (s *AccountService) GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error) {
	return s.AccountDao.GetTokenBalances(owner)
}

// Transfer transfer favorite token
func (s *AccountService) Transfer(token common.Address, fromAddress common.Address, toAddress common.Address, amount *big.Int) error {
	return s.AccountDao.Transfer(token, fromAddress, toAddress, amount)
}

// GetFavoriteTokens get favorite token
func (s *AccountService) GetFavoriteTokens(owner common.Address) (map[common.Address]bool, error) {
	return s.AccountDao.GetFavoriteTokens(owner)
}

// AddFavoriteToken add favorite token
func (s *AccountService) AddFavoriteToken(owner, token common.Address) error {
	return s.AccountDao.AddFavoriteToken(owner, token)
}

// DeleteFavoriteToken delete favorite token
func (s *AccountService) DeleteFavoriteToken(owner, token common.Address) error {
	return s.AccountDao.DeleteFavoriteToken(owner, token)
}
