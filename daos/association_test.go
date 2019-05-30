package daos

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/globalsign/mgo"
	"github.com/tomochain/tomoxsdk/app"
	"github.com/tomochain/tomoxsdk/types"
	"github.com/tomochain/tomoxsdk/utils"
)

func TestDepositHistory(t *testing.T) {
	err := app.LoadConfig("../config", "test")
	if err != nil {
		panic(err)
	}
	session, err := mgo.Dial(app.Config.MongoURL)
	if err != nil {
		panic(err)
	}

	db = &Database{session}
	associationDao := NewAssociationDao()

	// test get history
	chain := types.ChainEthereum
	associatedAddress := common.HexToAddress("0x59b8515e7ff389df6926cd52a086b0f1f46c630a")
	addressAssociation, err := associationDao.GetAssociationByChainAssociatedAddress(chain, associatedAddress)

	if err != nil {
		panic(err)
	}

	associationTransactions := []types.AssociationTransaction{}

	if addressAssociation != nil {
		for _, txEnvelope := range addressAssociation.TxEnvelopes {
			bytes := common.Hex2Bytes(txEnvelope)

			// t.Logf("Got bytes: %v", bytes)

			var associationTransaction types.AssociationTransaction
			err = rlp.DecodeBytes(bytes, &associationTransaction)
			if err != nil {
				continue
			}

			associationTransactions = append(associationTransactions, associationTransaction)
		}
	}

	utils.PrintJSON(associationTransactions)

}
