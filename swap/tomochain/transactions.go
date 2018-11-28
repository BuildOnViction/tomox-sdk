package tomochain

import (
	"strconv"

	"github.com/tomochain/backend-matching-engine/errors"
)

func (ac *AccountConfigurator) createAccountTransaction(destination string) error {
	transaction, err := ac.buildTransaction(
		ac.signerPublicKey.String(),
		ac.SignerPrivateKey,
		destination,
		ac.StartingBalance,
	)
	if err != nil {
		return errors.Wrap(err, "Error building transaction")
	}

	err = ac.submitTransaction(transaction)
	if err != nil {
		return errors.Wrap(err, "Error submitting a transaction")
	}

	return nil
}

// configureAccountTransaction is using a signer on an user accounts to configure the account.
func (ac *AccountConfigurator) configureAccountTransaction(destination, intermediateAssetCode, amount string) error {

	var tokenPrice string
	switch intermediateAssetCode {
	case "ETH":
		tokenPrice = ac.TokenPriceETH
	default:
		return errors.Errorf("Invalid intermediateAssetCode: $%s", intermediateAssetCode)
	}

	// // Send WETH token using smart contract
	// build.Payment(
	// 	build.SourceAccount{ac.DistributionPublicKey},
	// 	build.Destination{destination},
	// 	build.CreditAmount{
	// 		Code:   intermediateAssetCode,
	// 		Issuer: ac.IssuerPublicKey,
	// 		Amount: amount,
	// 	},
	// )

	transaction, err := ac.buildTransaction(destination, ac.SignerPrivateKey, tokenPrice)
	if err != nil {
		return errors.Wrap(err, "Error building a transaction")
	}

	err = ac.submitTransaction(transaction)
	if err != nil {
		return errors.Wrap(err, "Error submitting a transaction")
	}

	return nil
}

// removeTemporarySigner is removing temporary signer from an account.
func (ac *AccountConfigurator) removeTemporarySigner(destination string) error {
	// Remove signer ? need to remove this account wallet? ac.signerPublicKey

	transaction, err := ac.buildTransaction(destination, ac.SignerPrivateKey)
	if err != nil {
		return errors.Wrap(err, "Error building a transaction")
	}

	err = ac.submitTransaction(transaction)
	if err != nil {
		return errors.Wrap(err, "Error submitting a transaction")
	}

	return nil
}

// buildUnlockAccountTransaction creates and returns unlock account transaction.
func (ac *AccountConfigurator) buildUnlockAccountTransaction(source string) (string, error) {
	// Remove signer, ac.LockUnixTimestamp

	return ac.buildTransaction(source, ac.SignerPrivateKey)
}

func (ac *AccountConfigurator) buildTransaction(source string, signer string, params ...string) (string, error) {
	// muts := []build.TransactionMutator{
	// 	build.SourceAccount{source},
	// 	build.Network{ac.NetworkPassphrase},
	// }

	// if source == ac.signerPublicKey {
	// 	muts = append(muts, build.Sequence{ac.getSignerSequence()})
	// } else {
	// 	muts = append(muts, build.AutoSequence{ac.Horizon})
	// }

	// muts = append(muts, mutators...)
	// tx, err := build.Transaction(muts...)
	// if err != nil {
	// 	return "", err
	// }
	// txe, err := tx.Sign(signer)
	// if err != nil {
	// 	return "", err
	// }
	// return txe.Base64()
	return "hash", nil
}

func (ac *AccountConfigurator) submitTransaction(transaction string) error {
	logger.Info("Submitting transaction")

	err := ac.submitTransaction(transaction)
	if err != nil {
		ac.updateSignerSequence()
		logger.Error("Error submitting transaction")
		return errors.Wrap(err, "Error submitting transaction")
	}

	logger.Info("Transaction successfully submitted")
	return nil
}

func (ac *AccountConfigurator) updateSignerSequence() error {
	ac.signerSequenceMutex.Lock()
	defer ac.signerSequenceMutex.Unlock()

	account, err := ac.LoadAccount(ac.signerPublicKey.String())
	if err != nil {
		err = errors.Wrap(err, "Error loading issuing account")
		logger.Error(err)
		return err
	}

	ac.signerSequence, err = strconv.ParseUint(account.Sequence, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Invalid DistributionPublicKey sequence")
		logger.Error(err)
		return err
	}

	return nil
}

func (ac *AccountConfigurator) LoadAccount(publicKey string) (Account, error) {
	return Account{}, nil
}

func (ac *AccountConfigurator) getSignerSequence() uint64 {
	ac.signerSequenceMutex.Lock()
	defer ac.signerSequenceMutex.Unlock()
	ac.signerSequence++
	sequence := ac.signerSequence
	return sequence
}
