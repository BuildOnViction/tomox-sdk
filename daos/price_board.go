package daos

import (
	"fmt"
	"github.com/tomochain/tomoxsdk/errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
	"github.com/tomochain/tomoxsdk/app"
)

type PriceBoardDao struct {
}

// NewTokenDao returns a new instance of TokenDao.
func NewPriceBoardDao() *PriceBoardDao {
	return &PriceBoardDao{}
}

func (dao *PriceBoardDao) GetLatestQuotes() (map[string]float64, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/cryptocurrency/quotes/latest?symbol=%s&convert=USD", app.Config.CoinmarketcapAPIUrl, app.Config.SupportedCurrencies)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("X-CMC_PRO_API_KEY", app.Config.CoinmarketcapAPIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	status := gjson.Get(string(body), "status")
	statusErrorCode := status.Get("error_code")
	statusErrorMessage := status.Get("error_message")

	if statusErrorCode.Int() != 0 {
		logger.Error(statusErrorMessage.String())
		return nil, errors.New(statusErrorMessage.String())
	}

	data := gjson.Get(string(body), "data")
	result := make(map[string]float64)
	data.ForEach(func(key, value gjson.Result) bool {
		result[key.String()] = value.Get("quote.USD.price").Float()
		return true // keep iterating
	})

	return result, nil
}
