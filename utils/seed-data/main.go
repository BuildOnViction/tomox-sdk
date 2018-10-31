package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

var (
	app = cli.NewApp()
)

func init() {
	// Initialize the CLI app and start tomo
	app.Commands = []cli.Command{
		cli.Command{
			Name: "genesis",
			Action: func(c *cli.Context) error {
				return generateGenesis(c.String("cbf"), c.String("out"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "contract-build-folder, cbf", Value: "../../../contracts/build/contracts"},
				cli.StringFlag{Name: "output-folder, out", Value: "../../../protocol/OrderBook"},
			},
		},
		cli.Command{
			Name: "tokens",
			Action: func(c *cli.Context) error {
				return generateToken(c.String("cr"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "contract-result, cr", Value: "/contract-results.txt"},
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

type Token struct {
	Symbol          string `json:"symbol"`
	ContractAddress string `json:"contractAddress"`
}

type TokenCode struct {
	Code    string `json:"code"`
	Balance string `json:"balance"`
}

type Genesis struct {
	Alloc map[string]TokenCode `json:"alloc"`
}

func getTokenCode(buildFolder, symbol string) TokenCode {
	contractPath := path.Join(buildFolder, fmt.Sprintf("%s.json", symbol))
	byteValue, _ := ioutil.ReadFile(contractPath)
	var contract map[string]string
	json.Unmarshal(byteValue, &contract)
	tokenCode := TokenCode{
		Code: contract["deployedBytecode"],
		// Code:    contract["bytecode"],
		Balance: "0x0",
	}
	return tokenCode
}

func getAbsolutePath(basePath, folder string) string {
	if folder[0] == '/' {
		return folder
	}

	return path.Join(basePath, folder)

}

func generateToken(contractResultFile string) error {

	return nil
}

func generateGenesis(folder, outFolder string) error {
	_, fileName, _, _ := runtime.Caller(1)
	basePath := path.Dir(fileName)
	buildFolder := getAbsolutePath(basePath, folder)
	outputFolder := getAbsolutePath(basePath, outFolder)

	fmt.Printf("Contract folder :%s\n", buildFolder)

	templatePath := path.Join(basePath, "genesis.gohtml")
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Print(err)
		return err
	}

	// first step: read all tokens and deployedBytecode (bytecode of smartcontract without deploying by wallet but creation block)
	tokenPath := path.Join(basePath, "tokens.json")
	tokenFile, err := os.Open(tokenPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}
	defer tokenFile.Close()

	genesis := Genesis{
		Alloc: make(map[string]TokenCode),
	}
	scanner := bufio.NewScanner(tokenFile)
	var re = regexp.MustCompile(`^0x`)
	for scanner.Scan() {
		var token Token
		json.Unmarshal(scanner.Bytes(), &token)
		fmt.Printf("Token content :%v\n", token)
		// get deployedBytecode of the token
		tokenCode := getTokenCode(buildFolder, token.Symbol)
		contractAddress := strings.ToLower(re.ReplaceAllString(token.ContractAddress, ""))
		genesis.Alloc[contractAddress] = tokenCode
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	genesisPath := path.Join(outputFolder, "genesis.json")
	f, err := os.Create(genesisPath)
	tpl.Execute(f, genesis)
	if err != nil {
		log.Print("execute: ", err)
		return err
	}
	return nil

}
