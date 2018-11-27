package tomochain

import (
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
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

// AccountConfigurator is responsible for configuring new Tomochain accounts that
// participate in ICO.
// Infact, AccountConfigurator will be replaced by smart contract
type AccountConfigurator struct {
	IssuerPublicKey       string
	DistributionPublicKey string
	SignerPrivateKey      string
	LockUnixTimestamp     uint64
	TokenAssetCode        string
	TokenPriceBTC         string
	TokenPriceETH         string
	StartingBalance       string
	OnAccountCreated      func(destination string)
	OnExchanged           func(destination string)
	OnExchangedTimelocked func(destination, transaction string)

	signerPublicKey     common.Address
	signerSequence      uint64
	signerSequenceMutex sync.Mutex
	accountStatus       map[string]Status
	accountStatusMutex  sync.Mutex
}

type Account struct {
	Typ      string         `json:"type"`
	URL      accounts.URL   `json:"url"`
	Address  common.Address `json:"address"`
	Sequence string
}
