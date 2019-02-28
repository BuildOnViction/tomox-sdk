package daos

import (
	"encoding/json"
	"fmt"
	"github.com/tomochain/dex-server/app"
	"io/ioutil"
	"log"
	"net/http"
)

type PriceBoardDao struct {
}

// NewTokenDao returns a new instance of TokenDao.
func NewPriceBoardDao() *PriceBoardDao {
	return &PriceBoardDao{}
}

func (dao *PriceBoardDao) GetLatestQuotes() ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", app.Config.CoinmarketcapAPIUrl, "/cryptocurrency/quotes/latest?symbol=ETH,TOMO&convert=USD"), nil)
	req.Header.Add("X-CMC_PRO_API_KEY", app.Config.CoinmarketcapAPIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return json.Marshal(body)
}
