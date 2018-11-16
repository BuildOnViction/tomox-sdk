package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

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
			Name: "seeds",
			Action: func(c *cli.Context) error {
				filePath := c.String("ccf")
				return batch(
					filePath,
					generateConfig,
					generateTokens,
					generatePairs,
					generateAccounts,
				)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "client-config-folder, ccf", Value: "../../../client/src/config"},
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

type ImageInsert struct {
	ImageURL string `json:"url"`
	// meta can be null
	ImageMeta string `json:"meta"`
}

type TokenInsert struct {
	*Token
	Name string
	ImageInsert
	ID string
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
		return fmt.Errorf("Failed to read the configuration file: %s", err)
	}

	ethereumConfig := v.GetStringMap("ethereum")

	ethereumConfig["exchange_address"] = groups["Exchange"]
	ethereumConfig["weth_address"] = groups["WETH"]

	v.SetDefault("ethereum", ethereumConfig)

	err := v.WriteConfigAs(path.Join(configPath, "config.yaml"))

	return err
}

func generatePairs(filePath string) error {
	_, fileName, _, _ := runtime.Caller(1)
	basePath := path.Dir(fileName)
	// first create a list from pairs.json, then update it using matches
	pairsFile := path.Join(basePath, "pairs.json")
	contractResultFile := getAbsolutePath(basePath, fmt.Sprintf("%s/%s", filePath, "addresses.json"))
	groups := getGroupsFromContractResultFile(contractResultFile)
	buffer := &bytes.Buffer{}
	file, err := os.Open(pairsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var objList []map[string]interface{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var obj map[string]interface{}
		json.Unmarshal(scanner.Bytes(), &obj)
		objList = append(objList, obj)
	}

	for _, obj := range objList {
		if baseTokenAddress, ok := groups[obj["baseTokenSymbol"].(string)]; ok {
			obj["baseTokenAddress"] = baseTokenAddress
		}
		if quoteTokenAddress, ok := groups[obj["quoteTokenSymbol"].(string)]; ok {
			obj["quoteTokenAddress"] = quoteTokenAddress
		}
		bytes, _ := json.Marshal(obj)
		buffer.Write(bytes)
		buffer.WriteString("\n")
	}

	fmt.Println(buffer.String())
	ioutil.WriteFile(pairsFile, buffer.Bytes(), os.ModePerm)

	return nil
}

func generateAccounts(filePath string) error {
	_, fileName, _, _ := runtime.Caller(1)
	basePath := path.Dir(fileName)
	// first create a list from pairs.json, then update it using matches
	accountFile := path.Join(basePath, "accounts.json")
	contractResultFile := getAbsolutePath(basePath, fmt.Sprintf("%s/%s", filePath, "addresses.json"))
	groups := getGroupsFromContractResultFile(contractResultFile)

	bytes, _ := ioutil.ReadFile(accountFile)

	var obj map[string]interface{}
	json.Unmarshal(bytes, &obj)

	tokenBalances := obj["tokenBalances"].(map[string]interface{})
	updateTokenBalances := make(map[string]interface{})

	for oldAddress, tokenBalance := range tokenBalances {
		tokenBalanceMap := tokenBalance.(map[string]interface{})
		if address, ok := groups[tokenBalanceMap["symbol"].(string)]; ok {
			tokenBalanceMap["address"] = address
			updateTokenBalances[address.(string)] = tokenBalance
		} else {
			updateTokenBalances[oldAddress] = tokenBalance
		}
	}
	obj["tokenBalances"] = updateTokenBalances
	bytes, _ = json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(bytes))
	ioutil.WriteFile(accountFile, bytes, os.ModePerm)

	return nil
}

func generateTokens(filePath string) error {
	_, fileName, _, _ := runtime.Caller(1)
	basePath := path.Dir(fileName)
	contractResultFile := getAbsolutePath(basePath, fmt.Sprintf("%s/%s", filePath, "addresses.json"))
	imagesConfigFile := getAbsolutePath(basePath, fmt.Sprintf("%s/%s", filePath, "images.json"))
	imagesConfigBytes, _ := ioutil.ReadFile(imagesConfigFile)
	// with RawMessage we can deserialize whatever type
	var imagesConfigMap map[string]map[string]*ImageInsert
	json.Unmarshal(imagesConfigBytes, &imagesConfigMap)
	imagesConfig := imagesConfigMap["8888"]
	// fmt.Println(imagesConfig)

	tplStr := `{"_id":{"$oid":"{{.ID}}"},"name":"{{.Name}}","symbol":"{{.Symbol}}","contractAddress":"{{.ContractAddress}}","image":{"url":"{{.ImageURL}}","meta":"{{.ImageMeta}}"},"decimals":18,"quote":false,"createdAt":"Sun Sep 02 2018 17:34:37 GMT+0900 (Korean Standard Time)","updatedAt":"Sun Sep 02 2018 17:34:37 GMT+0900 (Korean Standard Time)"}`
	tpl, _ := template.New("token").Parse(tplStr)
	startIndex, _ := new(big.Int).SetString("5b8ba09da75a9b1320ca4974", 16)
	oneBig := big.NewInt(1)
	groups := getGroupsFromContractResultFile(contractResultFile)
	buffer := &bytes.Buffer{}
	for symbol, address := range groups {
		if symbol == "Exchange" {
			continue
		}

		startIndex = startIndex.Add(startIndex, oneBig)
		tokenInsert := &TokenInsert{
			Token: &Token{
				Symbol:          symbol,
				ContractAddress: address.(string),
			},
			Name: symbol,
			ID:   startIndex.Text(16),
		}

		imageData, ok := imagesConfig[symbol]
		if ok {
			tokenInsert.ImageURL = imageData.ImageURL
			tokenInsert.ImageMeta = imageData.ImageMeta
		}

		tpl.Execute(buffer, tokenInsert)
		buffer.WriteString("\n")
	}
	tokenFile := path.Join(basePath, "tokens.json")
	ioutil.WriteFile(tokenFile, buffer.Bytes(), os.ModePerm)
	fmt.Printf("Token json data: %s\n", buffer.String())
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
