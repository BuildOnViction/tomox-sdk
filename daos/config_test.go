package daos

import (
	"testing"

	"github.com/tomochain/backend-matching-engine/app"
	"github.com/tomochain/backend-matching-engine/types"
	mgo "gopkg.in/mgo.v2"
)

func init() {
	err := app.LoadConfig("../config", "test")
	if err != nil {
		panic(err)
	}
}

func TestConfigIncrementIndex(t *testing.T) {

	session, err := mgo.Dial(app.Config.MongoURL)
	if err != nil {
		panic(err)
	}

	db = &Database{session}
	configDao := NewConfigDao()
	configDao.IncrementAddressIndex(types.ChainEthereum)
	index, err := configDao.GetAddressIndex(types.ChainEthereum)
	t.Logf("Current Address Index: %d, err  :%v", index, err)
}
