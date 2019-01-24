package services

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/dex-server/app"
	"github.com/tomochain/dex-server/interfaces"
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/utils"
	"github.com/tomochain/dex-server/utils/math"
)

type ValidatorService struct {
	ethereumProvider interfaces.EthereumProvider
	accountDao       interfaces.AccountDao
	orderDao         interfaces.OrderDao
	pairDao          interfaces.PairDao
}

func NewValidatorService(
	ethereumProvider interfaces.EthereumProvider,
	accountDao interfaces.AccountDao,
	orderDao interfaces.OrderDao,
	pairDao interfaces.PairDao,
) *ValidatorService {

	return &ValidatorService{
		ethereumProvider,
		accountDao,
		orderDao,
		pairDao,
	}
}

func (s *ValidatorService) ValidateAvailableBalance(o *types.Order) error {
	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	pair, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	totalRequiredAmount := o.TotalRequiredSellAmount(pair)

	// balanceRecord, err := s.accountDao.GetTokenBalances(o.UserAddress)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }

	var sellTokenBalance *big.Int
	var sellTokenAllowance *big.Int

	// we implement retries in the case the provider connection fell asleep
	err = utils.Retry(3, func() error {
		sellTokenBalance, err = s.ethereumProvider.BalanceOf(o.UserAddress, o.SellToken())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	err = utils.Retry(3, func() error {
		sellTokenAllowance, err = s.ethereumProvider.Allowance(o.UserAddress, exchangeAddress, o.SellToken())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenLockedBalance, err := s.orderDao.GetUserLockedBalance(o.UserAddress, o.SellToken(), pair)
	if err != nil {
		logger.Error(err)
		return err
	}

	availableSellTokenBalance := math.Sub(sellTokenBalance, sellTokenLockedBalance)
	availableSellTokenAllowance := math.Sub(sellTokenAllowance, sellTokenLockedBalance)

	//Sell Token Balance
	fmt.Println(sellTokenBalance)
	fmt.Println(totalRequiredAmount)
	if sellTokenBalance.Cmp(totalRequiredAmount) == -1 {
		return fmt.Errorf("Insufficient %v Balance", o.SellTokenSymbol())
	}

	if availableSellTokenBalance.Cmp(totalRequiredAmount) == -1 {
		return fmt.Errorf("Insufficient % available", o.SellTokenSymbol())
	}

	if sellTokenAllowance.Cmp(totalRequiredAmount) == -1 {
		return fmt.Errorf("Insufficient %v Allowance", o.SellTokenSymbol())
	}

	if availableSellTokenAllowance.Cmp(totalRequiredAmount) == -1 {
		return fmt.Errorf("Insufficient %v allowance available", o.SellTokenSymbol())
	}

	// sellTokenBalanceRecord := balanceRecord[o.SellToken()]
	// if sellTokenBalanceRecord == nil {
	// 	return errors.New("Account error: Balance record not found")
	// }

	// sellTokenBalanceRecord.Balance.Set(sellTokenBalance)
	// sellTokenBalanceRecord.Allowance.Set(sellTokenAllowance)
	// err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.SellToken(), sellTokenBalanceRecord)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }

	return nil
}

func (s *ValidatorService) ValidateBalance(o *types.Order) error {
	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	pair, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	totalRequiredAmount := o.TotalRequiredSellAmount(pair)

	// balanceRecord, err := s.accountDao.GetTokenBalances(o.UserAddress)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }

	var sellTokenBalance *big.Int
	var sellTokenAllowance *big.Int

	// we implement retries in the case the provider connection fell asleep
	err = utils.Retry(3, func() error {
		sellTokenBalance, err = s.ethereumProvider.BalanceOf(o.UserAddress, o.SellToken())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	err = utils.Retry(3, func() error {
		sellTokenAllowance, err = s.ethereumProvider.Allowance(o.UserAddress, exchangeAddress, o.SellToken())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	//Sell Token Balance
	if sellTokenBalance.Cmp(totalRequiredAmount) == -1 {
		return fmt.Errorf("Insufficient %v Balance", o.SellTokenSymbol())
	}

	if sellTokenAllowance.Cmp(totalRequiredAmount) == -1 {
		return fmt.Errorf("Insufficient %v Allowance", o.SellTokenSymbol())
	}

	// sellTokenBalanceRecord := balanceRecord[o.SellToken()]
	// if sellTokenBalanceRecord == nil {
	// 	return errors.New("Account error: Balance record not found")
	// }

	// sellTokenBalanceRecord.Balance.Set(sellTokenBalance)
	// sellTokenBalanceRecord.Allowance.Set(sellTokenAllowance)
	// err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.SellToken(), sellTokenBalanceRecord)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }

	return nil
}
