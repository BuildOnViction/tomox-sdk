package services

import (
    "fmt"
    "math/big"

    "github.com/tomochain/tomox-sdk/interfaces"
    "github.com/tomochain/tomox-sdk/types"
    "github.com/tomochain/tomox-sdk/utils"
    "github.com/tomochain/tomox-sdk/utils/math"
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

// ValidateAvailableBalance get balance
func (s *ValidatorService) ValidateAvailableBalance(o *types.Order) error {
    logger.Info("ValidateAvailableBalance start...")
    pair, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
    if err != nil {
        logger.Error(err)
        return err
    }

    totalRequiredAmount := o.TotalRequiredSellAmount(pair)

    var sellTokenBalance *big.Int

    // we implement retries in the case the provider connection fell asleep
    err = utils.Retry(3, func() error {
        sellTokenBalance, err = s.ethereumProvider.Balance(o.UserAddress, o.SellToken())
        return nil
    })

    if err != nil {
        logger.Error(err)
        return err
    }
    pairs, err := s.pairDao.GetActivePairs()
    sellTokenLockedBalance, err := s.orderDao.GetUserLockedBalance(o.UserAddress, o.SellToken(), pairs)
    if err != nil {
        logger.Error(err)
        return err
    }

    availableSellTokenBalance := math.Sub(sellTokenBalance, sellTokenLockedBalance)

    //Sell Token Balance
    if sellTokenBalance.Cmp(totalRequiredAmount) == -1 {
        return fmt.Errorf("insufficient %v Balance", o.SellTokenSymbol())
    }

    if availableSellTokenBalance.Cmp(totalRequiredAmount) == -1 {
        return fmt.Errorf("insufficient %v available", o.SellTokenSymbol())
    }

    return nil
}

func (s *ValidatorService) ValidateBalance(o *types.Order) error {
    //exchangeAddress := common.HexToAddress(app.Config.Tomochain["exchange_address"])

    pair, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
    if err != nil {
        logger.Error(err)
        return err
    }

    totalRequiredAmount := o.TotalRequiredSellAmount(pair)

    balanceRecord, err := s.accountDao.GetTokenBalance(o.UserAddress, o.SellToken())
    if err != nil {
        logger.Error(err)
        return err
    }

    var sellTokenBalance *big.Int
    sellTokenBalance = balanceRecord.Balance

    //Sell Token Balance
    if sellTokenBalance.Cmp(totalRequiredAmount) == -1 {
        return fmt.Errorf("Insufficient %v Balance", o.SellTokenSymbol())
    }

    return nil
}
