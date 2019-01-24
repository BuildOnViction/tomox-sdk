package contracts

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	etherTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/tomochain/dex-server/app"
	"github.com/tomochain/dex-server/daos"
	"github.com/tomochain/dex-server/ethereum"
	"github.com/tomochain/dex-server/services"
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/utils"
	"github.com/tomochain/dex-server/utils/math"
	"github.com/tomochain/dex-server/utils/testutils"
	"github.com/tomochain/dex-server/utils/testutils/mocks"
)

func SetupTest() (*testutils.Deployer, *types.Wallet, common.Address, common.Address, *types.Wallet, *types.Wallet) {
	err := app.LoadConfig("../config", "")
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")

	_, err = daos.InitSession(nil)
	if err != nil {
		panic(err)
	}

	wallet := testutils.GetTestWallet()
	maker := testutils.GetTestWallet1()
	taker := testutils.GetTestWallet2()

	walletDao := new(mocks.WalletDao)
	walletDao.On("GetDefaultAdminWallet").Return(wallet, nil)

	walletService := services.NewWalletService(walletDao)
	txService := services.NewTxService(walletDao, wallet)
	accs := []common.Address{wallet.Address, maker.Address, taker.Address}
	gasLimit, err := strconv.ParseUint("47b7600", 16, 64)
	if err != nil {
		panic(err)
	}
	client := ethereum.NewSimulatedClientWithGasLimit(accs, gasLimit)
	deployer := testutils.NewDeployer(walletService, txService, client)
	if err != nil {
		panic(err)
	}

	feeAccount := common.HexToAddress(app.Config.Ethereum["fee_account"])
	wethToken := common.HexToAddress(app.Config.Ethereum["weth_address"])

	return deployer, wallet, feeAccount, wethToken, maker, taker
}

func TestSetFeeAccount(t *testing.T) {
	deployer, _, feeAccount, wethToken, _, _ := SetupTest()
	exchange, _, _, err := deployer.DeployExchange(feeAccount, wethToken)
	if err != nil {
		t.Errorf("Could not deploy exchange: %v", err)
	}

	simulator := deployer.Client.(*ethereum.SimulatedClient)
	simulator.Commit()

	txOpts, _ := exchange.DefaultTxOptions()
	newFeeAccount := testutils.GetTestAddress1()

	_, err = exchange.SetFeeAccount(newFeeAccount, txOpts)
	if err != nil {
		t.Errorf("Could not see new fee account: %v", err)
	}

	simulator.Commit()

	feeAccount, err = exchange.FeeAccount()
	if err != nil {
		t.Errorf("Error retrieving fee account address: %v", err)
	}

	if newFeeAccount != feeAccount {
		t.Errorf("Fee account not set correctly")
	}
}

func TestSetOperator(t *testing.T) {
	deployer, _, feeAccount, wethToken, _, _ := SetupTest()

	exchange, _, _, err := deployer.DeployExchange(feeAccount, wethToken)
	if err != nil {
		t.Errorf("Could not deploy exchange")
	}

	simulator := deployer.Client.(*ethereum.SimulatedClient)
	simulator.Commit()

	txOpts, _ := exchange.DefaultTxOptions()
	operator := testutils.GetTestAddress1()

	_, err = exchange.SetOperator(operator, true, txOpts)
	if err != nil {
		t.Errorf("Could not set operator: %v", err)
	}

	simulator.Commit()

	isOperator, err := exchange.Operator(operator)
	if err != nil {
		t.Errorf("Error calling the operator variable: %v", err)
	}

	if isOperator != true {
		t.Errorf("Operator variable should be equal to true but got false")
	}
}

func TestTrade(t *testing.T) {
	deployer, admin, feeAccount, wethToken, maker, taker := SetupTest()
	simulator := deployer.GetSimulator()

	pricepoint := big.NewInt(1e8)
	amount := big.NewInt(10)
	buyAmount := big.NewInt(1000)
	sellAmount := big.NewInt(100)
	// expires := big.NewInt(1e7)

	exchange, exchangeAddr, _, err := deployer.DeployExchange(feeAccount, wethToken)
	if err != nil {
		t.Errorf("Could not deploy exchange")
	}

	txOpts, _ := exchange.DefaultTxOptions()
	t.Logf("Admin public key: %s", crypto.PubkeyToAddress(admin.PrivateKey.PublicKey).Hex())
	// t.Logf("Admin public key: %s", txOpts.From.Hex())

	_, err = exchange.SetOperator(admin.Address, true, txOpts)
	if err != nil {
		t.Errorf("Could not set operator: %v", err)
	}

	relayerWallet := testutils.GetTestWallet4()
	_, err = exchange.SetFeeAccount(relayerWallet.Address, txOpts)
	if err != nil {
		t.Errorf("Could not see new fee account: %v", err)
	}

	//Initially Maker owns 1e18 units of sellToken and Taker owns 1e18 units buyToken
	sellToken, sellTokenAddr, _, err := deployer.DeployToken(maker.Address, sellAmount)
	if err != nil {
		t.Errorf("Error deploying token 1: %v", err)
	}

	// etherBalance, _ := simulator.BalanceAt(context.Background(), maker.Address, nil)
	// t.Logf("Ether balance is: %s", etherBalance.String())
	buyToken, buyTokenAddr, _, err := deployer.DeployToken(taker.Address, buyAmount)
	if err != nil {
		t.Errorf("Error deploying token 2: %v", err)
	}

	simulator.Commit()

	t.Logf("Maker address :%s, Taker address :%s", maker.Address.Hex(), taker.Address.Hex())
	exchange.PrintErrors()

	sellToken.SetTxSender(maker)
	_, err = sellToken.Approve(exchangeAddr, sellAmount)
	if err != nil {
		t.Errorf("Could not approve sellToken: %v", err)
	}

	buyToken.SetTxSender(taker)
	_, err = buyToken.Approve(exchangeAddr, buyAmount)
	if err != nil {
		t.Errorf("Could not approve buyToken: %v", err)
	}

	exchange.Interface.RegisterPair(txOpts, buyTokenAddr, sellTokenAddr, pricepoint)
	exchange.Interface.RegisterPair(txOpts, sellTokenAddr, buyTokenAddr, pricepoint)

	simulator.Commit()

	sellAllowed, err := sellToken.Allowance(maker.Address, exchangeAddr)
	buyAllowed, err := buyToken.Allowance(taker.Address, exchangeAddr)
	ok, _ := exchange.Operator(admin.Address)
	callTx := exchange.GetTxCallOptions()
	pairRegistered, _ := exchange.Interface.PairIsRegistered(callTx, sellTokenAddr, buyTokenAddr)
	t.Logf("Allowed :sell(%s) - buy(%s), has Operator: %t, is pair registered: %t",
		sellAllowed.String(), buyAllowed.String(), ok, pairRegistered)

	//Maker creates an order that exchanges 'sellAmount' of sellToken for 'buyAmount' of buyToken
	makerOrder := &types.Order{
		ExchangeAddress: exchangeAddr,
		Side:            types.SELL,
		Amount:          amount,
		PricePoint:      pricepoint,
		Nonce:           big.NewInt(0),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		BaseToken:       buyTokenAddr,
		QuoteToken:      sellTokenAddr,
		UserAddress:     maker.Address,
		Status:          "OPEN",
	}

	makerOrder.Sign(maker)

	//Taker creates an order that exchanges 'buyAmount' of buyToken for 'sellAmount' of sellToken
	takerOrder := &types.Order{
		ExchangeAddress: exchangeAddr,
		Side:            types.BUY,
		Amount:          amount,
		PricePoint:      pricepoint,
		Nonce:           big.NewInt(0),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		BaseToken:       sellTokenAddr,
		QuoteToken:      buyTokenAddr,
		UserAddress:     taker.Address,
		Status:          "OPEN",
	}

	takerOrder.Sign(taker)

	trade := types.NewTrade(makerOrder, takerOrder, amount, pricepoint)
	trade.Sign(admin)
	// err = trade.Validate()

	matches := types.NewMatches(
		[]*types.Order{makerOrder},
		takerOrder,
		[]*types.Trade{trade},
	)

	utils.TerminalLogger.Info("Orderbook matches: ")
	utils.PrintJSON(matches)

	// Now try to update balance directly on merkle trie, using matches result, then check balance

	// fake trade
	sellToken.SetTxSender(maker)
	_, err = sellToken.Transfer(taker.Address, amount)
	if err != nil {
		t.Error(err)
	}

	buyToken.SetTxSender(taker)
	_, err = buyToken.Transfer(maker.Address, amount)
	if err != nil {
		t.Error(err)
	}

	// real trade, txOpts is from admin
	txOpts.GasLimit = 3000000

	// priceMultiplier, _ := exchange.Interface.GetPairPricepointMultiplier(callTx, sellTokenAddr, buyTokenAddr)
	// t.Logf("Price multiplier :%s", priceMultiplier.String())
	// trans, err := exchange.Trade(matches, txOpts)
	// if err != nil {
	// 	t.Errorf("Could not execute trade: %v", err)
	// 	return
	// }
	// utils.TerminalLogger.Info("Trade transactions: ")
	// utils.PrintJSON(trans)

	simulator.Commit()

	// TokenSell: InitialSellTokenAmount + amount * (amountSell/amountBuy)
	sellTokenTakerBalance, _ := sellToken.BalanceOf(taker.Address)
	sellTokenMakerBalance, _ := sellToken.BalanceOf(maker.Address)
	buyTokenTakerBalance, _ := buyToken.BalanceOf(taker.Address)
	buyTokenMakerBalance, _ := buyToken.BalanceOf(maker.Address)

	t.Logf("Sell token balance is: maker(%s) - taker(%s), buy token balance is: maker(%s) - taker(%s)",
		sellTokenMakerBalance.String(), sellTokenTakerBalance.String(),
		buyTokenMakerBalance.String(), buyTokenTakerBalance.String())

	expectedSellAmount := math.Sub(sellAmount, amount)
	expectedBuyAmount := math.Sub(buyAmount, amount)

	if sellTokenTakerBalance.Cmp(amount) != 0 {
		t.Errorf("Expected Taker balance of sellToken to be equal to %v but got %v instead", amount, sellTokenTakerBalance)
	}

	if sellTokenMakerBalance.Cmp(expectedSellAmount) != 0 {
		t.Errorf("Expected Maker balance of sellToken to be equal to %v but got %v instead", expectedSellAmount, sellTokenMakerBalance)
	}

	if buyTokenTakerBalance.Cmp(expectedBuyAmount) != 0 {
		t.Errorf("Expected Taker balance of buyToken to be equal to %v but got %v instead", expectedBuyAmount, buyTokenTakerBalance)
	}

	if buyTokenMakerBalance.Cmp(amount) != 0 {
		t.Errorf("Expected Maker balance of buyToken to be equal to %v but got %v instead", amount, buyTokenMakerBalance)
	}
}

func TestReadTrie(t *testing.T) {
	var (
		key0, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key1, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		key2, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		addr0   = crypto.PubkeyToAddress(key0.PublicKey)
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
		ctx     = context.Background()
	)

	simulatedBackend := backends.NewSimulatedBackend(core.GenesisAlloc{
		addr0: {Balance: big.NewInt(1000000000)},
		addr1: {Balance: big.NewInt(1000000000)},
		addr2: {Balance: big.NewInt(1000000000)},
	}, 10000000)

	database := simulatedBackend.Database()

	triedb := trie.NewDatabase(database)

	tx := etherTypes.NewTransaction(
		0,
		addr1,
		big.NewInt(100),
		21000,
		nil,
		nil,
	)

	tx, _ = etherTypes.SignTx(tx, etherTypes.HomesteadSigner{}, key0)
	simulatedBackend.SendTransaction(ctx, tx)
	simulatedBackend.Commit()

	balance, _ := simulatedBackend.BalanceAt(ctx, addr0, nil)

	t.Logf("Balance of account0: %s", balance.String())

	utils.PrintJSON(triedb.Nodes())
}
