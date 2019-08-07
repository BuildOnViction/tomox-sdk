package relayer

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

// Configuration struct
type Configuration struct {
	Signer     *Signer
	SignerPath string `json:"signer_path"`
}

// NewConfiguration get config
func NewConfiguration() *Configuration {
	_, fileLocation, _, _ := runtime.Caller(1)
	file := filepath.Join(fileLocation, "./config.json")
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	configuration := &Configuration{}
	err = json.Unmarshal(raw, configuration)
	if err != nil {
		panic(err)
	}
	path := filepath.Join(fileLocation, configuration.SignerPath)
	signer := NewSigner(path, fileLocation)

	configuration.Signer = signer

	return configuration
}
