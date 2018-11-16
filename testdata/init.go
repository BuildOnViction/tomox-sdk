package testdata

import (
	"github.com/tomochain/backend-matching-engine/app"
	"github.com/tomochain/backend-matching-engine/daos"
)

func init() {
	// the test may be started from the home directory or a subdirectory
	err := app.LoadConfig("./config", "../config")
	if err != nil {
		panic(err)
	}
	// connect to the database
	if err := daos.InitSession(nil); err != nil {
		panic(err)
	}
}
