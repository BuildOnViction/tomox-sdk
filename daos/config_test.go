package daos

import (
	"testing"

	"github.com/globalsign/mgo"
	"github.com/tomochain/tomoxsdk/app"
	"github.com/tomochain/tomoxsdk/types"
)

func TestConfigIncrementIndex(t *testing.T) {
	err := app.LoadConfig("../config", "test")
	if err != nil {
		panic(err)
	}
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
