package tomochain

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/backend-matching-engine/errors"
	"github.com/tomochain/backend-matching-engine/types"
)

func (ac *AccountConfigurator) Start() error {

	logger.Info("TomochainAccountConfigurator starting")

	if !common.IsHexAddress(ac.IssuerPublicKey) {
		return errors.New("Invalid IssuerPublicKey")
	}

	if !common.IsHexAddress(ac.DistributionPublicKey) {
		return errors.New("Invalid DistributionPublicKey")
	}

	ac.accountStatus = make(map[string]Status)

	go ac.logStats()
	return nil
}

func (ac *AccountConfigurator) logStats() {
	for {
		logger.Infof("statuses: %v", ac.accountStatus)
		time.Sleep(15 * time.Second)
	}
}

// ConfigureAccount configures a new account that participated in ICO.
// * First it creates a new account.
// * Once a signer is replaced on the account, it creates trust lines and exchanges assets.
// from coinmarket place trust lines we get exchange rate and call OnExchange event to collect the exchange rate
func (ac *AccountConfigurator) ConfigureAccount(chain types.Chain, destination, assetCode, amount string) {

	logger.Info("Configuring Tomochain account")

	ac.setAccountStatus(destination, StatusCreatingAccount)
	defer func() {
		ac.removeAccountStatus(destination)
	}()

	// Check if account exists. If it is, skip creating it.
	for {
		// get from feed
		_, exists, err := ac.getAccount(chain, destination)
		if err != nil {
			logger.Error("Error loading account from Tomochain")
			time.Sleep(2 * time.Second)
			continue
		}

		if exists {
			break
		}

		logger.Info("Creating Tomochain account")
		err = ac.createAccountTransaction(chain, destination)
		if err != nil {
			logger.Error("Error creating Tomochain account")
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	if ac.OnAccountCreated != nil {
		ac.OnAccountCreated(chain, destination)
	}

	ac.setAccountStatus(destination, StatusWaitingForSigner)

	// Wait for signer changes...
	// check if association feed is correct
	for {
		account, err := ac.LoadAccount(chain, destination)
		if err != nil {
			logger.Error("Error loading account to check trustline")
			time.Sleep(2 * time.Second)
			continue
		}

		if ac.signerExistsOnly(account) {
			break
		}

		time.Sleep(2 * time.Second)
	}

	logger.Info("Signer found")

	ac.setAccountStatus(destination, StatusConfiguringAccount)

	// When signer was created we can configure account in Bifrost without requiring
	// the user to share the account's secret key.
	logger.Info("Sending token")
	err := ac.configureAccountTransaction(chain, destination, assetCode, amount)
	if err != nil {
		logger.Error("Error configuring an account")
		return
	}

	ac.setAccountStatus(destination, StatusRemovingSigner)

	if ac.LockUnixTimestamp == 0 {
		logger.Info("Removing temporary signer")
		err = ac.removeTemporarySigner(chain, destination)
		if err != nil {
			logger.Error("Error removing temporary signer")
			return
		}

		if ac.OnExchanged != nil {
			ac.OnExchanged(chain, destination)
		}
	} else {
		logger.Info("Creating unlock transaction to remove temporary signer")
		transaction, err := ac.buildUnlockAccountTransaction(destination)
		if err != nil {
			logger.Error("Error creating unlock transaction")
			return
		}

		if ac.OnExchangedTimelocked != nil {
			ac.OnExchangedTimelocked(chain, destination, transaction)
		}
	}

	logger.Info("Account successully configured")
}

func (ac *AccountConfigurator) LoadAccount(chain types.Chain, publicKey string) (*types.AddressAssociation, error) {
	if ac.LoadAccountHandler != nil {
		return ac.LoadAccountHandler(chain, publicKey)
	}

	return &types.AddressAssociation{}, nil
}

func (ac *AccountConfigurator) setAccountStatus(account string, status Status) {
	ac.accountStatusMutex.Lock()
	defer ac.accountStatusMutex.Unlock()
	ac.accountStatus[account] = status
}

func (ac *AccountConfigurator) removeAccountStatus(account string) {
	ac.accountStatusMutex.Lock()
	defer ac.accountStatusMutex.Unlock()
	delete(ac.accountStatus, account)
}

func (ac *AccountConfigurator) getAccount(chain types.Chain, account string) (*types.AddressAssociation, bool, error) {
	hAccount, err := ac.LoadAccount(chain, account)
	return hAccount, true, err
}

// signerExistsOnly returns true if account has exactly one signer and it's
// equal to `signerPublicKey`.
func (ac *AccountConfigurator) signerExistsOnly(account *types.AddressAssociation) bool {
	tempSignerFound := false

	if account.TomochainPublicKey == ac.signerPublicKey {
		tempSignerFound = true
	}

	return tempSignerFound
}
