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

func batch(filePath string, funcs ...func(string) error) error {
	var err error
	for _, funcObj := range funcs {
		err = funcObj(filePath)
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
				return batch(
					filePath,
					generateConfig,
				)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "contract-config-folder, ccf", Value: "../db/utils"},
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

func getGroupsFromContractResultFile(contractResultFile string) map[string]interface{} {
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
	// TODO: Fix this hard-coded part 8888
	return ret["8888"].(map[string]interface{})
}

func generateConfig(filePath string) error {
	_, fileName, _, _ := runtime.Caller(1)
	basePath := path.Dir(fileName)
	contractResultFile := getAbsolutePath(basePath, fmt.Sprintf("%s/%s", filePath, "addresses.json"))

	groups := getGroupsFromContractResultFile(contractResultFile)

	configPath := path.Join(basePath, "../../config")
	v := viper.New()
	v.SetConfigName("config.sample")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read the configuration file: %s", err)
	}

	ethereumConfig := v.GetStringMap("ethereum")

	ethereumConfig["exchange_address"] = groups["Exchange"]
	ethereumConfig["weth_address"] = groups["WETH"]

	v.SetDefault("ethereum", ethereumConfig)

	err := v.WriteConfigAs(path.Join(configPath, "config.yaml"))

	return err
}
