package daos

import (
	"errors"
	"fmt"
	m "math"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/math"
	"github.com/tomochain/tomox-sdk/ws"
)

// LendingOrderDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type LendingOrderDao struct {
	collectionName string
	dbName         string
}

// LendingOrderDaoOption opts for database option
type LendingOrderDaoOption = func(*LendingOrderDao) error

// NewLendingOrderDao returns a new instance of LendingOrderDao
func NewLendingOrderDao(opts ...LendingOrderDaoOption) *LendingOrderDao {
	dao := &LendingOrderDao{}
	dao.collectionName = "lending_items"
	dao.dbName = app.Config.DBName

	for _, op := range opts {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}

	index := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	}

	i1 := mgo.Index{
		Key: []string{"userAddress"},
	}

	i2 := mgo.Index{
		Key: []string{"status"},
	}

	i3 := mgo.Index{
		Key: []string{"collateralToken"},
	}

	i4 := mgo.Index{
		Key: []string{"lendingToken"},
	}
	i5 := mgo.Index{
		Key: []string{"createdAt"},
	}
	indexes := []mgo.Index{}
	indexes, err := db.Session.DB(dao.dbName).C(dao.collectionName).Indexes()
	if err == nil {
		if !existedIndex("index_lending_item_hash", indexes) {
			err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(index)
			if err != nil {
				panic(err)
			}
		}
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i1)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i2)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i3)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i4)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i5)
	if err != nil {
		panic(err)
	}

	return dao
}

// NewTopupDao topup dao
func NewTopupDao(opts ...LendingOrderDaoOption) *LendingOrderDao {
	dao := &LendingOrderDao{}
	dao.collectionName = "lending_topups"
	dao.dbName = app.Config.DBName

	for _, op := range opts {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}
	return dao
}

// NewRepayDao repay dao
func NewRepayDao(opts ...LendingOrderDaoOption) *LendingOrderDao {
	dao := &LendingOrderDao{}
	dao.collectionName = "lending_repays"
	dao.dbName = app.Config.DBName

	for _, op := range opts {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}
	return dao
}

// NewRecallDao recall dao
func NewRecallDao(opts ...LendingOrderDaoOption) *LendingOrderDao {
	dao := &LendingOrderDao{}
	dao.collectionName = "lending_recalls"
	dao.dbName = app.Config.DBName

	for _, op := range opts {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}
	return dao
}

// GetCollection get collection name
func (dao *LendingOrderDao) GetCollection() *mgo.Collection {
	return db.GetCollection(dao.dbName, dao.collectionName)
}

// Watch watch chaging database
func (dao *LendingOrderDao) Watch() (*mgo.ChangeStream, *mgo.Session, error) {
	return db.Watch(dao.dbName, dao.collectionName, mgo.ChangeStreamOptions{
		FullDocument:   mgo.UpdateLookup,
		MaxAwaitTimeMS: 500,
		BatchSize:      1000,
	})
}

// Create function performs the DB insertion task for LendingOrder collection
func (dao *LendingOrderDao) Create(o *types.LendingOrder) error {
	o.ID = bson.NewObjectId()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	if o.Status == "" {
		o.Status = "OPEN"
	}

	err := db.Create(dao.dbName, dao.collectionName, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetOrderCountByUserAddress get the total number of orders created by a user
// Return an integer and error
func (dao *LendingOrderDao) GetOrderCountByUserAddress(addr common.Address) (int, error) {
	q := bson.M{"userAddress": addr.Hex()}

	total, err := db.Count(dao.dbName, dao.collectionName, q)

	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return total, nil
}

// GetByID function fetches a single document from order collection based on mongoDB ID.
// Returns LendingOrder type struct
func (dao *LendingOrderDao) GetByID(id bson.ObjectId) (*types.LendingOrder, error) {
	var response *types.LendingOrder
	err := db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return response, err
}

// GetByHash function fetches a single document from order collection based on mongoDB ID.
// Returns LendingOrder type struct
func (dao *LendingOrderDao) GetByHash(hash common.Hash) (*types.LendingOrder, error) {
	q := bson.M{"hash": hash.Hex()}
	res := []types.LendingOrder{}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

// Drop drops all the order documents in the current database
func (dao *LendingOrderDao) Drop() error {
	err := db.DropCollection(dao.dbName, dao.collectionName)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *LendingOrderDao) getSideLendingOrderBook(term uint64, lendingToken common.Address, side string, srt int, limit ...int) ([]map[string]string, error) {

	sides := []map[string]string{}

	var lendingOrders []types.LendingOrder
	c := dao.GetCollection()

	// TODO: need to have limit
	err := c.Find(bson.M{
		"status":       bson.M{"$in": []string{types.OrderStatusOpen, types.OrderStatusPartialFilled}},
		"term":         strconv.FormatUint(term, 10),
		"lendingToken": lendingToken.Hex(),
		"side":         side,
	}).Sort("-createdAt").All(&lendingOrders)

	pa := make(map[string]*big.Int)
	for _, lendingOrder := range lendingOrders {
		interest := strconv.FormatUint(lendingOrder.Interest, 10)
		if val, ok := pa[interest]; ok {
			pa[interest] = math.Sub(math.Add(val, lendingOrder.Quantity), lendingOrder.FilledAmount)
		} else {
			pa[interest] = math.Sub(lendingOrder.Quantity, lendingOrder.FilledAmount)
		}
	}

	for p, a := range pa {
		s := map[string]string{
			"interest": p,
			"amount":   a.String(),
		}

		sides = append(sides, s)
	}

	sort.SliceStable(sides, func(i, j int) bool {
		return math.ToBigInt(sides[i]["interest"]).Cmp(math.ToBigInt(sides[j]["interest"])) == (0 - srt)
	})

	return sides, err
}

// GetLendingOrderBook get lending order token
func (dao *LendingOrderDao) GetLendingOrderBookInDb(term uint64, lendingToken common.Address) ([]map[string]string, []map[string]string, error) {

	borrow, err := dao.getSideLendingOrderBook(term, lendingToken, types.BORROW, -1)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	lend, err := dao.getSideLendingOrderBook(term, lendingToken, types.LEND, -1)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return borrow, lend, nil
}

func (dao *LendingOrderDao) GetLendingOrderBook(term uint64, lendingToken common.Address) ([]map[string]string, []map[string]string, error) {
	rpcClient, err := rpc.DialHTTP(app.Config.Tomochain["http_url"])
	defer rpcClient.Close()

	var result interface{}

	err = rpcClient.Call(&result, "tomox_getInvests", lendingToken.Hex(), term)
	asks := []map[string]string{}
	if result != nil && err == nil {
		for k, v := range result.(map[string]interface{}) {
			s := map[string]string{
				"interest": k,
				"amount":   fmt.Sprintf("%.0f", v.(float64)),
			}
			asks = append(asks, s)
		}

		sort.SliceStable(asks, func(i, j int) bool {
			return math.ToBigInt(asks[i]["interest"]).Cmp(math.ToBigInt(asks[j]["interest"])) == -1
		})
	}

	err = rpcClient.Call(&result, "tomox_getBorrows", lendingToken.Hex(), term)
	bids := []map[string]string{}
	if result != nil && err == nil {
		for k, v := range result.(map[string]interface{}) {
			s := map[string]string{
				"interest": k,
				"amount":   fmt.Sprintf("%.0f", v.(float64)),
			}
			bids = append(bids, s)
		}
		sort.SliceStable(bids, func(i, j int) bool {
			return math.ToBigInt(bids[i]["interest"]).Cmp(math.ToBigInt(bids[j]["interest"])) == 1
		})
	}

	if err != nil {
		logger.Error(err)
		return bids, asks, nil
	}

	return bids, asks, nil
}

// LendingOrderMsg for tomox rpc
type LendingOrderMsg struct {
	AccountNonce    hexutil.Uint64 `json:"nonce"    gencodec:"required"`
	Quantity        hexutil.Big    `json:"quantity,omitempty"`
	RelayerAddress  common.Address `json:"relayerAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	CollateralToken common.Address `json:"collateralToken,omitempty"`
	LendingToken    common.Address `json:"lendingToken,omitempty"`
	Interest        hexutil.Uint64 `json:"interest,omitempty"`
	Term            hexutil.Uint64 `json:"term,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	LendingID       hexutil.Uint64 `json:"lendingID,omitempty"`
	LendingTradeID  hexutil.Uint64 `json:"tradeId,omitempty"`
	AutoTopUp       bool           `json:"autoTopUp,omitempty"`
	// Signature values
	V hexutil.Big `json:"v" gencodec:"required"`
	R hexutil.Big `json:"r" gencodec:"required"`
	S hexutil.Big `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash common.Hash `json:"hash" rlp:"-"`
}

// AddNewLendingOrder add order
func (dao *LendingOrderDao) AddNewLendingOrder(o *types.LendingOrder) error {
	rpcClient, err := rpc.DialHTTP(app.Config.Tomochain["http_url"])
	defer rpcClient.Close()
	bigstr := o.Nonce.String()
	n, err := strconv.ParseInt(bigstr, 10, 64)
	if err != nil {
		return err
	}
	V := big.NewInt(int64(o.Signature.V))
	R := o.Signature.R.Big()
	S := o.Signature.S.Big()

	autoTopUp := (uint64(o.AutoTopUp) == uint64(1))

	msg := LendingOrderMsg{
		AccountNonce:    hexutil.Uint64(uint64(n)),
		Quantity:        hexutil.Big(*o.Quantity),
		RelayerAddress:  o.RelayerAddress,
		UserAddress:     o.UserAddress,
		CollateralToken: o.CollateralToken,
		AutoTopUp:       autoTopUp,
		Term:            hexutil.Uint64(o.Term),
		Interest:        hexutil.Uint64(o.Interest),
		LendingToken:    o.LendingToken,
		Status:          "NEW",
		Side:            o.Side,
		Type:            o.Type,
		Hash:            o.Hash,
		V:               hexutil.Big(*V),
		R:               hexutil.Big(*R),
		S:               hexutil.Big(*S),
	}
	var result interface{}
	logger.Info("tomox_sendLending", o.Status, o.Hash.Hex())
	err = rpcClient.Call(&result, "tomox_sendLending", msg)

	if err != nil {
		logger.Error(err)
		ws.SendLendingOrderMessage("ERROR", o.UserAddress, OrderErrorMsg{
			Message: err.Error(),
		})
		return err
	}

	o.Status = "ADDED"
	ws.SendLendingOrderMessage(types.LENDING_ORDER_ADDED, o.UserAddress, o)
	return nil
}

// CancelLendingOrder cancel order
func (dao *LendingOrderDao) CancelLendingOrder(o *types.LendingOrder) error {

	rpcClient, err := rpc.DialHTTP(app.Config.Tomochain["http_url"])
	defer rpcClient.Close()
	bigstr := o.Nonce.String()
	n, err := strconv.ParseInt(bigstr, 10, 64)
	if err != nil {
		return err
	}
	V := big.NewInt(int64(o.Signature.V))
	R := o.Signature.R.Big()
	S := o.Signature.S.Big()

	msg := LendingOrderMsg{
		AccountNonce:    hexutil.Uint64(uint64(n)),
		Status:          o.Status,
		Hash:            o.Hash,
		LendingID:       hexutil.Uint64(o.LendingID),
		UserAddress:     o.UserAddress,
		CollateralToken: o.CollateralToken,
		LendingToken:    o.LendingToken,
		Term:            hexutil.Uint64(o.Term),
		Interest:        hexutil.Uint64(o.Interest),
		RelayerAddress:  o.RelayerAddress,
		V:               hexutil.Big(*V),
		R:               hexutil.Big(*R),
		S:               hexutil.Big(*S),
	}
	var result interface{}
	logger.Info("tomox_sendLending", o.Status, o.Hash.Hex(), o.LendingID, o.UserAddress.Hex(), n)
	err = rpcClient.Call(&result, "tomox_sendLending", msg)

	if err != nil {
		logger.Error(err)
		ws.SendLendingOrderMessage("ERROR", o.UserAddress, OrderErrorMsg{
			Message: err.Error(),
		})
		return err
	}

	ws.SendLendingOrderMessage(types.LENDING_ORDER_CANCELLED, o.UserAddress, o)
	return nil
}

// RepayLendingOrder send repay transaction
func (dao *LendingOrderDao) RepayLendingOrder(o *types.LendingOrder) error {

	rpcClient, err := rpc.DialHTTP(app.Config.Tomochain["http_url"])
	defer rpcClient.Close()
	bigstr := o.Nonce.String()
	n, err := strconv.ParseInt(bigstr, 10, 64)
	if err != nil {
		return err
	}
	V := big.NewInt(int64(o.Signature.V))
	R := o.Signature.R.Big()
	S := o.Signature.S.Big()

	msg := LendingOrderMsg{
		AccountNonce:   hexutil.Uint64(uint64(n)),
		Status:         o.Status,
		UserAddress:    o.UserAddress,
		LendingToken:   o.LendingToken,
		Term:           hexutil.Uint64(o.Term),
		LendingTradeID: hexutil.Uint64(o.LendingTradeID),
		RelayerAddress: o.RelayerAddress,
		Type:           o.Type,
		V:              hexutil.Big(*V),
		R:              hexutil.Big(*R),
		S:              hexutil.Big(*S),
	}
	var result interface{}
	logger.Info("tomox_sendLending", o.Status, o.Hash.Hex(), o.LendingTradeID, o.UserAddress.Hex(), n)
	err = rpcClient.Call(&result, "tomox_sendLending", msg)

	if err != nil {
		logger.Error(err)
		ws.SendLendingOrderMessage("ERROR", o.UserAddress, OrderErrorMsg{
			Message: err.Error(),
		})
		return err
	}

	// ws.SendLendingOrderMessage(types.LENDING_ORDER_REPAYED, o.UserAddress, o)
	return nil
}

// TopupLendingOrder send top up lending transaction
func (dao *LendingOrderDao) TopupLendingOrder(o *types.LendingOrder) error {

	rpcClient, err := rpc.DialHTTP(app.Config.Tomochain["http_url"])
	defer rpcClient.Close()
	bigstr := o.Nonce.String()
	n, err := strconv.ParseInt(bigstr, 10, 64)
	if err != nil {
		return err
	}
	V := big.NewInt(int64(o.Signature.V))
	R := o.Signature.R.Big()
	S := o.Signature.S.Big()

	msg := LendingOrderMsg{
		AccountNonce:   hexutil.Uint64(uint64(n)),
		Status:         o.Status,
		UserAddress:    o.UserAddress,
		LendingToken:   o.LendingToken,
		Term:           hexutil.Uint64(o.Term),
		LendingTradeID: hexutil.Uint64(o.LendingTradeID),
		RelayerAddress: o.RelayerAddress,
		Quantity:       hexutil.Big(*o.Quantity),
		Type:           o.Type,
		V:              hexutil.Big(*V),
		R:              hexutil.Big(*R),
		S:              hexutil.Big(*S),
	}
	var result interface{}
	logger.Info("tomox_sendLending", o.Status, o.Hash.Hex(), o.LendingTradeID, o.UserAddress.Hex(), n)
	err = rpcClient.Call(&result, "tomox_sendLending", msg)

	if err != nil {
		logger.Error(err)
		ws.SendLendingOrderMessage("ERROR", o.UserAddress, OrderErrorMsg{
			Message: err.Error(),
		})
		return err
	}

	// ws.SendLendingOrderMessage(types.LENDING_ORDER_TOPUPED, o.UserAddress, o)
	return nil
}

// GetLendingNonce get nonce of lending order
func (dao *LendingOrderDao) GetLendingNonce(userAddress common.Address) (uint64, error) {
	rpcClient, err := rpc.DialHTTP(app.Config.Tomochain["http_url"])

	defer rpcClient.Close()

	if err != nil {
		logger.Error(err)
		return 0, err
	}

	var result interface{}
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	err = rpcClient.Call(&result, "tomox_getLendingOrderCount", userAddress)

	if err != nil {
		logger.Error(err)
		return 0, err
	}
	logger.Info("OrderNonce:", result)
	s := result.(string)
	s = strings.TrimPrefix(s, "0x")
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0, nil
	}

	return n, nil
}

// GetLendingOrderBookInterest get amount from interest
func (dao *LendingOrderDao) GetLendingOrderBookInterest(term uint64, lendingToken common.Address, interest uint64, side string) (*big.Int, error) {
	var orders []types.LendingOrder
	c := dao.GetCollection()

	//TODO: need to have limit
	err := c.Find(bson.M{
		"status":       bson.M{"$in": []string{types.OrderStatusOpen, types.OrderStatusPartialFilled}},
		"term":         strconv.FormatUint(term, 10),
		"lendingToken": lendingToken.Hex(),
		"side":         side,
		"interest":     strconv.FormatUint(interest, 10),
	}).Sort("-createdAt").All(&orders)

	amount := big.NewInt(0)

	for _, order := range orders {
		amount = math.Sub(math.Add(amount, order.Quantity), order.FilledAmount)
	}

	return amount, err
}

// UpdateStatus update lending status
func (dao *LendingOrderDao) UpdateStatus(h common.Hash, status string) error {
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{
		"status": status,
	}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// UpdateFilledAmount update lending amount
func (dao *LendingOrderDao) UpdateFilledAmount(hash common.Hash, value *big.Int) error {
	q := bson.M{"hash": hash.Hex()}
	res := []types.LendingOrder{}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		logger.Error(err)
		return err
	}

	o := res[0]
	status := ""
	filledAmount := math.Add(o.FilledAmount, value)

	if math.IsEqualOrSmallerThan(filledAmount, big.NewInt(0)) {
		filledAmount = big.NewInt(0)
		status = "OPEN"
	} else if math.IsEqualOrGreaterThan(filledAmount, o.Quantity) {
		filledAmount = o.Quantity
		status = "FILLED"
	} else {
		status = "PARTIAL_FILLED"
	}

	update := bson.M{"$set": bson.M{
		"status":       status,
		"filledAmount": filledAmount.String(),
	}}

	err = db.Update(dao.dbName, dao.collectionName, q, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetLendingOrders filter lending order
func (dao *LendingOrderDao) GetLendingOrders(lendingSpec types.LendingSpec, sort []string, offset int, size int) (*types.LendingRes, error) {

	q := bson.M{}
	q["relayer"] = lendingSpec.RelayerAddress.Hex()
	if lendingSpec.UserAddress != "" {
		q["userAddress"] = lendingSpec.UserAddress
	}
	if lendingSpec.DateFrom != 0 || lendingSpec.DateTo != 0 {
		dateFilter := bson.M{}
		if lendingSpec.DateFrom != 0 {

			dateFilter["$gte"] = time.Unix(lendingSpec.DateFrom, 0)
		}
		if lendingSpec.DateTo != 0 {
			dateFilter["$lt"] = time.Unix(lendingSpec.DateTo, 0)
		}
		q["createdAt"] = dateFilter
	}
	if lendingSpec.LendingToken != "" {
		q["lendingToken"] = lendingSpec.LendingToken
	}
	if lendingSpec.CollateralToken != "" {
		q["collateralToken"] = lendingSpec.CollateralToken
	}

	if lendingSpec.Side != "" {
		q["side"] = strings.ToUpper(lendingSpec.Side)
	}
	if lendingSpec.Status != "" {
		q["status"] = strings.ToUpper(lendingSpec.Status)
	}
	if lendingSpec.Type != "" {
		q["type"] = strings.ToUpper(lendingSpec.Type)
	}
	if lendingSpec.Hash != "" {
		q["hash"] = lendingSpec.Hash
	}
	if lendingSpec.Term != "" {
		q["term"] = lendingSpec.Term
	}
	var res types.LendingRes
	lendings := []*types.LendingOrder{}
	c, err := db.GetEx(dao.dbName, dao.collectionName, q, sort, offset, size, &lendings)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	res.Total = c
	for i := range lendings {
		lendings[i].Signature = nil
	}
	res.LendingItems = lendings
	return &res, nil
}

// GetLastTokenPrice get last token price
func (dao *LendingOrderDao) GetLastTokenPrice(bToken common.Address, qToken common.Address) (*big.Int, error) {
	rpcClient, err := rpc.DialHTTP(app.Config.Tomochain["http_url"])

	defer rpcClient.Close()

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var result interface{}
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = rpcClient.Call(&result, "tomox_getLastEpochPrice", bToken, qToken)

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	s := result.(string)
	s = strings.TrimPrefix(s, "0x")
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		return nil, errors.New("error parse price")
	}

	return n, nil
}

//GetUserLockedBalance return balance using selling
func (dao *LendingOrderDao) GetUserLockedBalance(account common.Address, token common.Address, decimals int) (*big.Int, error) {
	var orders []*types.LendingOrder
	q := bson.M{
		"$or": []bson.M{
			{
				"userAddress":  account.Hex(),
				"status":       bson.M{"$in": []string{types.OrderStatusOpen, types.OrderStatusPartialFilled}},
				"lendingToken": token.Hex(),
				"side":         "INVEST",
			},
			{
				"userAddress":     account.Hex(),
				"status":          bson.M{"$in": []string{types.OrderStatusOpen, types.OrderStatusPartialFilled}},
				"collateralToken": token.Hex(),
				"side":            "BORROW",
			},
		},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	totalLockedBalance := big.NewInt(0)
	totalInvest := big.NewInt(0)
	totalBorrow := big.NewInt(0)
	lendingTokenList := make(map[common.Address]*big.Int)
	for _, o := range orders {
		remainingAmount := math.Sub(o.Quantity, o.FilledAmount)
		if o.Side == types.LEND {
			totalInvest = math.Add(totalInvest, remainingAmount)
		} else {
			if v, ok := lendingTokenList[o.LendingToken]; ok {
				v = v.Add(v, remainingAmount)
			} else {
				lendingTokenList[o.LendingToken] = new(big.Int).Add(big.NewInt(0), remainingAmount)
			}
		}
	}
	collateralDecimals := big.NewInt(int64(m.Pow10(decimals)))
	for lt, q := range lendingTokenList {
		collateralPrice, err := dao.GetLastTokenPrice(token, lt)
		if err != nil {
			return nil, err
		}
		collateralAmount := new(big.Int).Mul(q, collateralDecimals)
		collateralAmount = math.Mul(collateralAmount, big.NewInt(int64(types.LendingRate)))
		collateralAmount = new(big.Int).Div(collateralAmount, collateralPrice)
		collateralAmount = math.Div(collateralAmount, big.NewInt(100))
		totalBorrow = totalBorrow.Add(totalBorrow, collateralAmount)
	}
	totalLockedBalance = new(big.Int).Add(totalInvest, totalBorrow)
	return totalLockedBalance, nil
}
