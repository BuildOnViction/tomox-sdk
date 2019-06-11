package contracts

import (
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomoxsdk/app"
	"github.com/tomochain/tomoxsdk/contracts/contractsinterfaces"
	"github.com/tomochain/tomoxsdk/daos"
	"github.com/tomochain/tomoxsdk/ethereum"
	"github.com/tomochain/tomoxsdk/services"
	"github.com/tomochain/tomoxsdk/types"
	"github.com/tomochain/tomoxsdk/utils/math"
	"github.com/tomochain/tomoxsdk/utils/testutils"
	"github.com/tomochain/tomoxsdk/utils/testutils/mocks"
)

func SetupTokenTest() (*testutils.Deployer, *types.Wallet) {
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
	walletDao := new(mocks.WalletDao)
	walletDao.On("GetDefaultAdminWallet").Return(wallet, nil)

	walletService := services.NewWalletService(walletDao)
	txService := services.NewTxService(walletDao, wallet)

	client := ethereum.NewSimulatedClient([]common.Address{wallet.Address})
	deployer := testutils.NewDeployer(walletService, txService, client)

	return deployer, wallet
}

func TestBalanceOf(t *testing.T) {
	deployer, wallet := SetupTokenTest()

	receiver := testutils.GetTestAddress1()
	amount := big.NewInt(1e18)

	token, _, _, err := deployer.DeployToken(receiver, amount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.GetSimulator()

	etherBalance, _ := simulator.BalanceAt(context.Background(), wallet.Address, nil)
	t.Logf("Ether balance is: %s", etherBalance.String())
	// commit sending tokens
	simulator.Commit()
	newEtherBalance, _ := simulator.BalanceAt(context.Background(), wallet.Address, nil)
	t.Logf("Ether balance is: %s, lost: %s", newEtherBalance.String(), math.Sub(etherBalance, newEtherBalance).String())

	balance, err := token.BalanceOf(receiver)
	if err != nil {
		t.Errorf("Error retrieving token balance: %v", err)
	}

	if balance.Cmp(amount) != 0 {
		t.Errorf("Token balance incorrect. Expected %v but instead got %v", amount, balance)
	}
}

func TestTotalSupply(t *testing.T) {
	deployer, _ := SetupTokenTest()

	receiver := testutils.GetTestAddress1()
	amount := big.NewInt(1e18)

	token, _, _, err := deployer.DeployToken(receiver, amount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.Client.(*ethereum.SimulatedClient)
	simulator.Commit()

	supply, err := token.TotalSupply()
	if err != nil {
		t.Errorf("Error retrieving total supply")
	}

	if supply.Cmp(amount) != 0 {
		t.Errorf("Token Balance Incorrect. Expected %v but instead got %v", amount, supply)
	}
}

func TestTransfer(t *testing.T) {
	deployer, wallet := SetupTokenTest()

	owner := wallet.Address
	receiver := testutils.GetTestAddress1()
	initialAmount := big.NewInt(1e18)
	transferAmount := big.NewInt(5e17)

	token, _, _, err := deployer.DeployToken(owner, initialAmount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.Client.(*ethereum.SimulatedClient)
	simulator.Commit()

	_, err = token.Transfer(receiver, transferAmount)
	if err != nil {
		t.Errorf("Could not transfer tokens: %v", err)
	}

	simulator.Commit()

	receiverBalance, err := token.BalanceOf(receiver)
	if err != nil {
		t.Errorf("Could not retrieve receiver balance %v", err)
	}
	if receiverBalance.Cmp(big.NewInt(5e17)) != 0 {
		t.Errorf("Expected receiver balance to be equal to 1/2e18 but got %v instead", receiverBalance)
	}
}

func TestApprove(t *testing.T) {
	deployer, wallet := SetupTokenTest()

	owner := wallet.Address
	spender := testutils.GetTestAddress2()
	amount := big.NewInt(1e18)

	token, _, _, err := deployer.DeployToken(owner, amount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.Client.(*ethereum.SimulatedClient)
	simulator.Commit()

	_, err = token.Approve(spender, amount)
	if err != nil {
		t.Errorf("Could not approve tokens: %v", err)
	}

	simulator.Commit()

	allowance, err := token.Allowance(owner, spender)
	if err != nil {
		t.Errorf("Could not retrieve receiver allowance %v", err)
	}
	if allowance.Cmp(amount) != 0 {
		t.Errorf("Expected receiver balance to be equal to 1/2e18 but got %v instead", allowance)
	}
}

func TestTransferEvent(t *testing.T) {
	deployer, wallet := SetupTokenTest()

	owner := wallet.Address
	receiver := testutils.GetTestAddress2()

	logs := []*contractsinterfaces.TokenTransfer{}
	amount := big.NewInt(1e18)
	done := make(chan bool)

	token, _, _, err := deployer.DeployToken(owner, amount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.Client.(*ethereum.SimulatedClient)
	simulator.Commit()

	events, err := token.ListenToTransferEvents()
	if err != nil {
		t.Errorf("Could not open transfer events channel")
	}

	go func() {
		for {
			event := <-events
			logs = append(logs, event)
			done <- true
		}
	}()

	_, err = token.Transfer(receiver, amount)
	if err != nil {
		t.Errorf("Could not transfer tokens: %v", err)
	}

	simulator.Commit()
	<-done

	if len(logs) != 1 {
		t.Errorf("Events log has not the correct length")
	}

	parsedTransfer := logs[0]
	if parsedTransfer.From != owner {
		t.Errorf("Event 'From' field is not correct")
	}
	if parsedTransfer.To != receiver {
		t.Errorf("Event 'To' field is not correct")
	}
	if parsedTransfer.Value.Cmp(amount) != 0 {
		t.Errorf("Event 'Amount' field is not correct")
	}
}
