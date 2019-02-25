package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/spf13/viper"
	"gopkg.in/urfave/cli.v1"
)

var (
	app = cli.NewApp()
)

var networks = map[string]string{
	"ethereum":         "1",
	"rinkeby":          "4",
	"tomochain":        "88",
	"tomochainTestnet": "89",
	"development":      "8888",
}

func batch(filePath string, networkId string, funcs ...func(string, string) error) error {
	var err error
	for _, funcObj := range funcs {
		err = funcObj(filePath, networkId)
		if err != nil {
			break
		}
	}
	return err
}

func init() {
	// Initialize the CLI app and start tomo
	app.Commands = []cli.Command{
		cli.Command{
			Name: "seeds",
			Action: func(c *cli.Context) error {
				filePath := c.String("ccf")
				networkId := getNetworkID(os.Args[2])
				return batch(
					filePath,
					networkId,
					generateConfig,
				)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "contract-config-folder, ccf", Value: "../../deployment/utils"},
			},
		},
	}
}

func main() {

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func getAbsolutePath(basePath, folder string) string {
	if folder[0] == '/' {
		return folder
	}

	return path.Join(basePath, folder)

}

func getGroupsFromContractResultFile(contractResultFile string, networkId string) map[string]interface{} {
	// // now matching data from contract-resultFile
	// resultData, _ := ioutil.ReadFile(contractResultFile)
	// // ?m: is notation tell this will match multiline
	// tokenAndAddress := regexp.MustCompile(`(?m:^\s*([\w]+)\s*:\s*(.*?)\s*$)`)
	// // TOMO: 0x4f696e8a1a3fb3aea9f72eb100ea8d97c5130b32
	// groups = make(map[string]string)
	// matches := tokenAndAddress.FindAllStringSubmatch(string(resultData), -1)
	// for _, match := range matches {
	// 	groups[match[1]] = match[2]
	// }

	// return groups
	var ret map[string]interface{}
	bytes, _ := ioutil.ReadFile(contractResultFile)
	json.Unmarshal(bytes, &ret)

	return ret[networkId].(map[string]interface{})
}

func generateConfig(filePath string, networkId string) error {
	_, fileName, _, _ := runtime.Caller(1)
	basePath := path.Dir(fileName)
	contractResultFile := getAbsolutePath(basePath, fmt.Sprintf("%s/%s", filePath, "addresses.json"))

	groups := getGroupsFromContractResultFile(contractResultFile, networkId)

	configPath := path.Join(basePath, "../../config")
	v := viper.New()

	// Choose config file based on deployment network
	switch networkId {
	case networks["tomochain"]:
		v.SetConfigName("config.prod")
	case networks["tomochainTestnet"]:
		v.SetConfigName("config.dev")
	case networks["development"]:
		v.SetConfigName("config.local")
	default:
		v.SetConfigName("config.local")
	}

	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read the configuration file: %s", err)
	}

	ethereumConfig := v.GetStringMap("ethereum")

	ethereumConfig["exchange_address"] = groups["Exchange"]

	v.SetDefault("ethereum", ethereumConfig)

	err := v.WriteConfigAs(path.Join(configPath, "config.yaml"))

	return err
}

func getNetworkID(networkName string) string {
	return networks[networkName]
}
