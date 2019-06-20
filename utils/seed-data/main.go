package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/spf13/viper"
	"gopkg.in/urfave/cli.v1"
)

var (
	app = cli.NewApp()
)

func batch(networkId string, funcs ...func(string) error) error {
	var err error
	for _, funcObj := range funcs {
		err = funcObj(networkId)
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
				networkId := os.Args[2]
				return batch(
					networkId,
					generateConfig,
				)
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

func generateConfig(networkId string) error {
	_, fileName, _, _ := runtime.Caller(1)
	basePath := path.Dir(fileName)

	configPath := path.Join(basePath, "../../config")
	v := viper.New()

	// Choose config file based on deployment network
	switch networkId {
	case "88":
		v.SetConfigName("config.prod")
	case "89":
		v.SetConfigName("config.dev")
	case "8888":
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

	v.SetDefault("ethereum", ethereumConfig)

	err := v.WriteConfigAs(path.Join(configPath, "config.yaml"))

	return err
}
