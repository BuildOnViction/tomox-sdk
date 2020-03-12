package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

// MockServices is a that tolds different mock services to be passed
// around easily for testing setup

// UintToPaddedString converts an int to string of length 19 by padding with 0
func UintToPaddedString(num int64) string {
	return fmt.Sprintf("%019d", num)
}

// GetTickChannelID is used to get the channel id for OHLCV data streaming
// it takes pairname, duration and units of data streaming
func GetTickChannelID(bt, qt common.Address, unit string, duration int64) string {
	pair := GetPairKey(bt, qt)
	return fmt.Sprintf("%s::%d::%s", pair, duration, unit)
}

// GetPairKey return the pair key identifier corresponding to two
func GetPairKey(bt, qt common.Address) string {
	return strings.ToLower(fmt.Sprintf("%s::%s", bt.Hex(), qt.Hex()))
}

func GetTradeChannelID(bt, qt common.Address) string {
	return strings.ToLower(fmt.Sprintf("%s::%s", bt.Hex(), qt.Hex()))
}

// GetLendingTradeChannelID get channel from term and lending token
func GetLendingTradeChannelID(term uint64, lendingToken common.Address) string {
	return strings.ToLower(fmt.Sprintf("%s::%s", strconv.FormatUint(term, 10), lendingToken.Hex()))
}

func GetOHLCVChannelID(bt, qt common.Address, unit string, duration int64) string {
	pair := GetPairKey(bt, qt)
	return fmt.Sprintf("%s::%d::%s", pair, duration, unit)
}
func GetLendingOhlcvChannelID(term uint64, lendingToken common.Address, unit string, duration int64) string {
	return fmt.Sprintf("%d::%s::%d::%s", term, lendingToken.Hex(), duration, unit)
}

func GetOrderBookChannelID(bt, qt common.Address) string {
	return strings.ToLower(fmt.Sprintf("%s::%s", bt.Hex(), qt.Hex()))
}
func GetLendingOrderBookChannelID(term uint64, lendingToken common.Address) string {
	return strings.ToLower(fmt.Sprintf("%s::%s", strconv.FormatUint(term, 10), lendingToken.Hex()))
}
func GetPriceBoardChannelID(bt, qt common.Address) string {
	return strings.ToLower(fmt.Sprintf("%s::%s", bt.Hex(), qt.Hex()))
}

func GetMarketsChannelID(channel string) string {
	return strings.ToLower(channel)
}

func Retry(retries int, fn func() error) error {
	if err := fn(); err != nil {
		retries--
		if retries <= 0 {
			return err
		}

		// preventing thundering herd problem (https://en.wikipedia.org/wiki/Thundering_herd_problem)
		time.Sleep(time.Second)

		return Retry(retries, fn)
	}

	return nil
}

func PrintJSON(x interface{}) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Print(string(b), "\n")
}

func JSON(x interface{}) string {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	return fmt.Sprint(string(b), "\n")
}

func PrintError(msg string, err error) {
	log.Printf("\n%v: %v\n", msg, err)
}

// Util function to handle unused variables while testing
func Use(...interface{}) {

}

// GetAddressFromPublicKey get derived address from public key
func GetAddressFromPublicKey(pk []byte) common.Address {
	var buf []byte
	hash := sha3.NewKeccak256()
	hash.Write(pk[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	publicAddr := hexutil.Encode(buf[12:])

	publicAddress := common.HexToAddress(publicAddr)

	return publicAddress
}

// Union get set of two array
func Union(a, b []common.Address) []common.Address {
	m := make(map[common.Address]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			a = append(a, item)
		}
	}
	return a
}

func IsNativeTokenByAddress(address common.Address) bool {
	return (address.Hex() == "0x0000000000000000000000000000000000000001")
}

// GetModTime get round time by step
func GetModTime(ts, interval int64, unit string) (int64, int64) {
	var modTime, intervalInSeconds int64
	switch unit {
	case "sec":
		intervalInSeconds = interval
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "min":
		intervalInSeconds = interval * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "hour":
		intervalInSeconds = interval * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "day":
		intervalInSeconds = interval * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "week":
		intervalInSeconds = interval * 7 * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "month":
		d := time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.UTC).Day()
		intervalInSeconds = interval * int64(d) * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "year":
		// Number of days in current year
		d := time.Date(time.Now().Year()+1, 1, 1, 0, 0, 0, 0, time.UTC).Sub(time.Date(time.Now().Year(), 0, 0, 0, 0, 0, 0, time.UTC)).Hours() / 24
		intervalInSeconds = interval * int64(d) * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))
	}

	return modTime, intervalInSeconds
}

// UnitToSecond time uint to second
func UnitToSecond(interval int64, unit string) int64 {
	var intervalInSeconds int64
	switch unit {
	case "sec":
		intervalInSeconds = interval
	case "min":
		intervalInSeconds = interval * 60
	case "hour":
		intervalInSeconds = interval * 60 * 60
	case "day":
		intervalInSeconds = interval * 24 * 60 * 60
	case "week":
		intervalInSeconds = interval * 7 * 24 * 60 * 60
	case "month":
		d := time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.UTC).Day()
		intervalInSeconds = interval * int64(d) * 24 * 60 * 60
	case "year":
		d := time.Date(time.Now().Year()+1, 1, 1, 0, 0, 0, 0, time.UTC).Sub(time.Date(time.Now().Year(), 0, 0, 0, 0, 0, 0, time.UTC)).Hours() / 24
		intervalInSeconds = interval * int64(d) * 24 * 60 * 60

	}

	return intervalInSeconds
}

// ToBigInt string to bigint
func ToBigInt(s string) *big.Int {
	res := big.NewInt(0)
	res.SetString(s, 10)
	return res
}
