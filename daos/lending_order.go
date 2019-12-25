package daos

import (
	"math/big"
	"sort"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	dao.collectionName = "lending_orders"
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
		if !existedIndex("index_order_hash", indexes) {
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
func (dao *LendingOrderDao) GetLendingOrderBook(term uint64, lendingToken common.Address) ([]map[string]string, []map[string]string, error) {

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

// LendingOrderMsg for tomox rpc
type LendingOrderMsg struct {
	AccountNonce    uint64         `json:"nonce"    gencodec:"required"`
	Quantity        *big.Int       `json:"quantity,omitempty"`
	RelayerAddress  common.Address `json:"relayerAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	CollateralToken common.Address `json:"collateralToken,omitempty"`
	LendingToken    common.Address `json:"lendingToken,omitempty"`
	Interest        uint64         `json:"interest,omitempty"`
	Term            uint64         `json:"term,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	LendingID       uint64         `json:"lendingID,omitempty"`
	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

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

	msg := LendingOrderMsg{
		AccountNonce:    uint64(n),
		Quantity:        o.Quantity,
		RelayerAddress:  o.RelayerAddress,
		UserAddress:     o.UserAddress,
		CollateralToken: o.CollateralToken,
		Term:            o.Term,
		Interest:        o.Interest,
		LendingToken:    o.LendingToken,
		Status:          "NEW",
		Side:            o.Side,
		Type:            o.Type,
		Hash:            o.Hash,
		V:               V,
		R:               R,
		S:               S,
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
	o.Status = "OPEN"
	dao.Create(o)
	o.Status = "ADDED"
	ws.SendLendingOrderMessage("ORDER_ADDED", o.UserAddress, o)
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
		AccountNonce:    uint64(n),
		Status:          o.Status,
		Hash:            o.Hash,
		LendingID:       o.LendingID,
		UserAddress:     o.UserAddress,
		CollateralToken: o.CollateralToken,
		LendingToken:    o.LendingToken,
		Term:            o.Term,
		Interest:        o.Interest,
		RelayerAddress:  o.RelayerAddress,
		V:               V,
		R:               R,
		S:               S,
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

	ws.SendLendingOrderMessage("ORDER_CANCELLED", o.UserAddress, o)
	dao.Create(o)
	return nil
}

// GetLendingNonce get nonce of order
func (dao *LendingOrderDao) GetLendingNonce(userAddress common.Address) (uint64, error) {
	return 0, nil
}
