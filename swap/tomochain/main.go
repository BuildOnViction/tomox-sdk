package tomochain

import (
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/backend-matching-engine/swap/config"
	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/utils"
)

// Status describes status of account processing
type Status string

var logger = utils.EngineLogger

const (
	StatusCreatingAccount    Status = "creating_account"
	StatusWaitingForSigner   Status = "waiting_for_signer"
	StatusConfiguringAccount Status = "configuring_account"
	StatusRemovingSigner     Status = "removing_signer"
)

type LoadAccountHandler func(chain types.Chain, publicKey string) (*types.AddressAssociation, error)

// AccountConfigurator is responsible for configuring new Tomochain accounts that
// participate in ICO.
// Infact, AccountConfigurator will be replaced by smart contract
type AccountConfigurator struct {
	IssuerPublicKey       string
	DistributionPublicKey string

	LockUnixTimestamp uint64
	TokenAssetCode    string
	TokenPriceBTC     string
	TokenPriceETH     string
	StartingBalance   string

	LoadAccountHandler    LoadAccountHandler
	OnSubmitTransaction   func(chain types.Chain, destination string, transaction string) error
	OnAccountCreated      func(chain types.Chain, destination string)
	OnExchanged           func(chain types.Chain, destination string)
	OnExchangedTimelocked func(chain types.Chain, destination, transaction string)

	signerPublicKey  common.Address
	signerPrivateKey *ecdsa.PrivateKey

	accountStatus      map[string]Status
	accountStatusMutex sync.Mutex
}

func NewAccountConfigurator(cfg *config.Config) *AccountConfigurator {
	return &AccountConfigurator{
		IssuerPublicKey:       cfg.Tomochain.IssuerPublicKey,
		DistributionPublicKey: cfg.Tomochain.DistributionPublicKey,
		signerPublicKey:       cfg.SignerPublicKey(),
		signerPrivateKey:      cfg.SignerPrivateKey(),
		TokenAssetCode:        cfg.Tomochain.TokenAssetCode,
		StartingBalance:       cfg.Tomochain.StartingBalance,
		LockUnixTimestamp:     cfg.Tomochain.LockUnixTimestamp,
	}
}
