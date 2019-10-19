package daos

import (
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/math"
)

// OrderDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type OrderDao struct {
	collectionName string
	dbName         string
}

type OrderDaoOption = func(*OrderDao) error

func OrderDaoDBOption(dbName string) func(dao *OrderDao) error {
	return func(dao *OrderDao) error {
		dao.dbName = dbName
		return nil
	}
}

// NewOrderDao returns a new instance of OrderDao
func NewOrderDao(opts ...OrderDaoOption) *OrderDao {
	dao := &OrderDao{}
	dao.collectionName = "orders"
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
		Key: []string{"baseToken"},
	}

	i4 := mgo.Index{
		Key: []string{"quoteToken"},
	}

	i5 := mgo.Index{
		Key:       []string{"pricepoint"},
		Collation: &mgo.Collation{NumericOrdering: true, Locale: "en"},
	}

	i6 := mgo.Index{
		Key: []string{"baseToken", "quoteToken", "pricepoint"},
	}

	i7 := mgo.Index{
		Key: []string{"side", "status"},
	}

	i8 := mgo.Index{
		Key: []string{"baseToken", "quoteToken", "side", "status"},
	}

	i9 := mgo.Index{
		Key: []string{"createdAt"},
	}

	err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
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

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i6)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i7)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i8)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i9)
	if err != nil {
		panic(err)
	}

	return dao
}

func (dao *OrderDao) GetCollection() *mgo.Collection {
	return db.GetCollection(dao.dbName, dao.collectionName)
}

func (dao *OrderDao) Watch() (*mgo.ChangeStream, *mgo.Session, error) {
	return db.Watch(dao.dbName, dao.collectionName, mgo.ChangeStreamOptions{FullDocument: mgo.UpdateLookup})
}

// Create function performs the DB insertion task for Order collection
func (dao *OrderDao) Create(o *types.Order) error {
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

func (dao *OrderDao) DeleteByHashes(hashes ...common.Hash) error {
	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"hash": bson.M{"$in": hashes}})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) Delete(orders ...*types.Order) error {
	hashes := []common.Hash{}
	for _, o := range orders {
		hashes = append(hashes, o.Hash)
	}

	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"hash": bson.M{"$in": hashes}})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Update function performs the DB updations task for Order collection
// corresponding to a particular order ID
func (dao *OrderDao) Update(id bson.ObjectId, o *types.Order) error {
	o.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": id}, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) Upsert(id bson.ObjectId, o *types.Order) error {
	o.UpdatedAt = time.Now()

	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": id}, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) UpsertByHash(h common.Hash, o *types.Order) error {
	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"hash": h.Hex()}, types.OrderBSONUpdate{o})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) UpdateAllByHash(h common.Hash, o *types.Order) error {
	o.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"hash": h.Hex()}, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) FindAndModify(h common.Hash, o *types.Order) (*types.Order, error) {
	o.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	updated := &types.Order{}
	change := mgo.Change{
		Update:    types.OrderBSONUpdate{o},
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}

	err := db.FindAndModify(dao.dbName, dao.collectionName, query, change, &updated)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return updated, nil
}

//UpdateByHash updates fields that are considered updateable for an order.
func (dao *OrderDao) UpdateByHash(h common.Hash, o *types.Order) error {
	o.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{
		"pricepoint":   o.PricePoint.String(),
		"amount":       o.Amount.String(),
		"status":       o.Status,
		"filledAmount": o.FilledAmount.String(),
		"updatedAt":    o.UpdatedAt,
	}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) UpdateOrderStatus(h common.Hash, status string) error {
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

func (dao *OrderDao) UpdateOrderStatusesByHashes(status string, hashes ...common.Hash) ([]*types.Order, error) {
	hexes := []string{}
	for _, h := range hashes {
		hexes = append(hexes, h.Hex())
	}

	query := bson.M{"hash": bson.M{"$in": hexes}}
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"status":    status,
		},
	}

	err := db.UpdateAll(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	orders := []*types.Order{}
	err = db.Get(dao.dbName, dao.collectionName, query, 0, 0, &orders)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	return orders, nil
}

func (dao *OrderDao) UpdateOrderFilledAmount(hash common.Hash, value *big.Int) error {
	q := bson.M{"hash": hash.Hex()}
	res := []types.Order{}
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
	} else if math.IsEqualOrGreaterThan(filledAmount, o.Amount) {
		filledAmount = o.Amount
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

func (dao *OrderDao) UpdateOrderFilledAmounts(hashes []common.Hash, amount []*big.Int) ([]*types.Order, error) {
	hexes := []string{}
	orders := []*types.Order{}
	for i := range hashes {
		hexes = append(hexes, hashes[i].Hex())
	}

	query := bson.M{"hash": bson.M{"$in": hexes}}
	err := db.Get(dao.dbName, dao.collectionName, query, 0, 0, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	updatedOrders := []*types.Order{}
	for i, o := range orders {
		status := ""
		filledAmount := math.Sub(o.FilledAmount, amount[i])

		if math.IsEqualOrSmallerThan(filledAmount, big.NewInt(0)) {
			filledAmount = big.NewInt(0)
			status = "OPEN"
		} else if math.IsEqualOrGreaterThan(filledAmount, o.Amount) {
			filledAmount = o.Amount
			status = "FILLED"
		} else {
			status = "PARTIAL_FILLED"
		}

		query := bson.M{"hash": o.Hash.Hex()}
		update := bson.M{"$set": bson.M{
			"status":       status,
			"filledAmount": filledAmount.String(),
		}}
		change := mgo.Change{
			Update:    update,
			Upsert:    true,
			Remove:    false,
			ReturnNew: true,
		}

		updated := &types.Order{}
		err := db.FindAndModify(dao.dbName, dao.collectionName, query, change, updated)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		updatedOrders = append(updatedOrders, updated)
	}

	return updatedOrders, nil
}

// GetOrderCountByUserAddress get the total number of orders created by a user
// Return an integer and error
func (dao *OrderDao) GetOrderCountByUserAddress(addr common.Address) (int, error) {
	q := bson.M{"userAddress": addr.Hex()}

	total, err := db.Count(dao.dbName, dao.collectionName, q)

	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return total, nil
}

// GetByID function fetches a single document from order collection based on mongoDB ID.
// Returns Order type struct
func (dao *OrderDao) GetByID(id bson.ObjectId) (*types.Order, error) {
	var response *types.Order
	err := db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return response, err
}

// GetByHash function fetches a single document from order collection based on mongoDB ID.
// Returns Order type struct
func (dao *OrderDao) GetByHash(hash common.Hash) (*types.Order, error) {
	q := bson.M{"hash": hash.Hex()}
	res := []types.Order{}

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

// GetByHashes return Order DAO by address
func (dao *OrderDao) GetByHashes(hashes []common.Hash) ([]*types.Order, error) {
	hexes := []string{}
	for _, h := range hashes {
		hexes = append(hexes, h.Hex())
	}

	q := bson.M{"hash": bson.M{"$in": hexes}}
	res := []*types.Order{}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetByUserAddress function fetches list of orders from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetByUserAddress(addr, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error) {
	if limit == nil {
		limit = []int{types.DefaultLimit}
	}

	var fromTemp, toTemp int64
	now := time.Now()

	if to == 0 {
		toTemp = now.Unix()
		to = toTemp
	}

	if from == 0 {
		fromTemp = now.AddDate(-1, 0, 0).Unix()
		from = fromTemp
	}

	var res []*types.Order
	var q bson.M

	if (bt == common.Address{} || qt == common.Address{}) {
		q = bson.M{
			"userAddress": addr.Hex(),
			"createdAt": bson.M{
				"$gte": time.Unix(from, 0),
				"$lt":  time.Unix(to, 0),
			},
		}
	} else {
		q = bson.M{
			"userAddress": addr.Hex(),
			"baseToken":   bt.Hex(),
			"quoteToken":  qt.Hex(),
			"createdAt": bson.M{
				"$gte": time.Unix(from, 0),
				"$lt":  time.Unix(to, 0),
			},
		}
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, limit[0], &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return []*types.Order{}, nil
	}

	return res, nil
}

func (dao *OrderDao) removeSignature(order *types.Order) {
	order.Signature = nil
}

// GetOrders filter order
func (dao *OrderDao) GetOrders(orderSpec types.OrderSpec, sort []string, offset int, size int) (*types.OrderRes, error) {

	q := bson.M{}
	if orderSpec.UserAddress != "" {
		q["userAddress"] = orderSpec.UserAddress
	}
	if orderSpec.DateFrom != 0 || orderSpec.DateTo != 0 {
		dateFilter := bson.M{}
		if orderSpec.DateFrom != 0 {

			dateFilter["$gte"] = time.Unix(orderSpec.DateFrom, 0)
		}
		if orderSpec.DateTo != 0 {
			dateFilter["$lt"] = time.Unix(orderSpec.DateTo, 0)
		}
		q["createdAt"] = dateFilter
	}
	if orderSpec.BaseToken != "" {
		q["baseToken"] = orderSpec.BaseToken
	}
	if orderSpec.QuoteToken != "" {
		q["quoteToken"] = orderSpec.BaseToken
	}

	if orderSpec.Side != "" {
		q["side"] = strings.ToUpper(orderSpec.Side)
	}
	if orderSpec.Status != "" {
		q["status"] = strings.ToUpper(orderSpec.Status)
	}
	if orderSpec.OrderType != "" {
		q["type"] = strings.ToUpper(orderSpec.OrderType)
	}
	var res types.OrderRes
	orders := []*types.Order{}
	c, err := db.GetEx(dao.dbName, dao.collectionName, q, sort, offset, size, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	res.Total = c
	for i := range orders {
		dao.removeSignature(orders[i])
	}
	res.Orders = orders
	return &res, nil
}

// GetOpenOrdersByUserAddress function fetches list of open/partial filled orders from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetOpenOrdersByUserAddress(addr common.Address) ([]*types.Order, error) {
	var res []*types.Order
	var q bson.M

	q = bson.M{
		"userAddress": addr.Hex(),
		"status":      bson.M{"$in": []string{types.OrderStatusNew, types.OrderStatusOpen, types.OrderStatusPartialFilled}},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return []*types.Order{}, nil
	}

	return res, nil
}

// GetCurrentByUserAddress function fetches list of open/partial orders from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetCurrentByUserAddress(addr common.Address, limit ...int) ([]*types.Order, error) {
	if limit == nil {
		limit = []int{types.DefaultLimit}
	}

	var res []*types.Order
	q := bson.M{
		"userAddress": addr.Hex(),
		"status": bson.M{"$in": []string{
			types.OrderStatusNew,
			types.OrderStatusOpen,
			types.OrderStatusPartialFilled,
		},
		},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, limit[0], &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return []*types.Order{}, nil
	}

	return res, nil
}

// GetHistoryByUserAddress function fetches list of orders which are not in open/partial order status
// from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetHistoryByUserAddress(addr, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error) {
	if limit == nil {
		limit = []int{types.DefaultLimit}
	}

	// Set default time range
	var fromTemp, toTemp int64
	now := time.Now()

	if to == 0 {
		toTemp = now.Unix()
		to = toTemp
	}

	if from == 0 {
		fromTemp = now.AddDate(-1, 0, 0).Unix()
		from = fromTemp
	}

	var res []*types.Order
	var q bson.M

	if (bt == common.Address{} || qt == common.Address{}) {
		q = bson.M{
			"userAddress": addr.Hex(),
			"createdAt": bson.M{
				"$gte": time.Unix(from, 0),
				"$lt":  time.Unix(to, 0),
			},
			"status": bson.M{"$nin": []string{
				types.OrderStatusNew,
				types.OrderStatusOpen,
				types.OrderStatusPartialFilled,
			},
			},
		}
	} else {
		q = bson.M{
			"userAddress": addr.Hex(),
			"baseToken":   bt.Hex(),
			"quoteToken":  qt.Hex(),
			"status": bson.M{"$nin": []string{
				types.OrderStatusNew,
				types.OrderStatusOpen,
				types.OrderStatusPartialFilled,
			},
			},
			"createdAt": bson.M{
				"$gte": time.Unix(from, 0),
				"$lt":  time.Unix(to, 0),
			},
		}
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, limit[0], &res)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return []*types.Order{}, nil
	}

	return res, nil
}

//GetUserLockedBalance return balance using selling
func (dao *OrderDao) GetUserLockedBalance(account common.Address, token common.Address, pairs []*types.Pair) (*big.Int, error) {
	var orders []*types.Order

	q := bson.M{
		"$or": []bson.M{
			{
				"userAddress": account.Hex(),
				"status":      bson.M{"$in": []string{types.OrderStatusNew, types.OrderStatusOpen, types.OrderStatusPartialFilled}},
				"quoteToken":  token.Hex(),
				"side":        "BUY",
			},
			{
				"userAddress": account.Hex(),
				"status":      bson.M{"$in": []string{types.OrderStatusNew, types.OrderStatusOpen, types.OrderStatusPartialFilled}},
				"baseToken":   token.Hex(),
				"side":        "SELL",
			},
		},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	totalLockedBalance := big.NewInt(0)
	for _, o := range orders {
		for _, p := range pairs {
			if p.BaseTokenSymbol == o.BaseTokenSymbol() && p.QuoteTokenSymbol == o.QuoteTokenSymbol() {
				if o.Side == types.BUY {
					remainingAmount := math.Sub(o.Amount, o.FilledAmount)
					amount := math.Mul(remainingAmount, o.PricePoint)
					w := math.Exp(big.NewInt(10), big.NewInt(int64(p.BaseTokenDecimals)))
					amount = math.Div(amount, w)
					totalLockedBalance = math.Add(totalLockedBalance, amount)

				} else {
					remainingAmount := math.Sub(o.Amount, o.FilledAmount)
					totalLockedBalance = math.Add(totalLockedBalance, remainingAmount)
				}
				break
			}
		}

	}

	return totalLockedBalance, nil
}

func (dao *OrderDao) GetRawOrderBook(p *types.Pair) ([]*types.Order, error) {
	var orders []*types.Order
	q := []bson.M{
		{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{types.OrderStatusNew, types.OrderStatusOpen, types.OrderStatusPartialFilled}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
			},
		},
		{
			"$sort": bson.M{
				"price": 1,
			},
		},
	}

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return orders, nil
}

func (dao *OrderDao) GetSideOrderBook(p *types.Pair, side string, sort int, limit ...int) ([]map[string]string, error) {

	sides := []map[string]string{}
	if p == nil {
		return sides, nil
	}

	sideQuery := []bson.M{
		{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{types.OrderStatusNew, types.OrderStatusOpen, types.OrderStatusPartialFilled}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
				"side":       side,
			},
		},
		{
			"$group": bson.M{
				"_id":        bson.M{"$toDecimal": "$price"},
				"pricepoint": bson.M{"$first": "$price"},
				"amount": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{{"$toDecimal": "$quantity"}, {"$toDecimal": "$filledAmount"}},
					},
				},
			},
		},
		{
			"$sort": bson.M{
				"_id": sort,
			},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"pricepoint": bson.M{"$toString": "$pricepoint"},
				"amount":     bson.M{"$toString": "$amount"},
			},
		},
	}

	if limit != nil {
		sideQuery = append(sideQuery, bson.M{
			"$limit": limit[0],
		})
	}

	err := db.Aggregate(dao.dbName, dao.collectionName, sideQuery, &sides)

	return sides, err
}

// GetOrderBook get best bids descending and best asks ascending
func (dao *OrderDao) GetOrderBook(p *types.Pair) ([]map[string]string, []map[string]string, error) {

	bids, err := dao.GetSideOrderBook(p, types.BUY, -1)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	asks, err := dao.GetSideOrderBook(p, types.SELL, 1)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return bids, asks, nil
}

func (dao *OrderDao) GetOrderBookPricePoint(p *types.Pair, pp *big.Int, side string) (*big.Int, error) {
	q := []bson.M{
		{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{types.OrderStatusNew, types.OrderStatusOpen, types.OrderStatusPartialFilled}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
				"price":      pp.String(),
				"side":       side,
			},
		},
		{
			"$group": bson.M{
				"_id":        bson.M{"$toDecimal": "$price"},
				"pricepoint": bson.M{"$first": "$price"},
				"amount": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{{"$toDecimal": "$quantity"}, {"$toDecimal": "$filledAmount"}},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"pricepoint": 1,
				"amount":     bson.M{"$toString": "$amount"},
			},
		},
	}

	res := []map[string]string{}
	err := db.Aggregate(dao.dbName, dao.collectionName, q, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return math.ToBigInt(res[0]["amount"]), nil
}

// Drop drops all the order documents in the current database
func (dao *OrderDao) Drop() error {
	err := db.DropCollection(dao.dbName, dao.collectionName)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Aggregate function calls the aggregate pipeline of mongodb
func (dao *OrderDao) Aggregate(q []bson.M) ([]*types.OrderData, error) {
	logger.Info("Query aggregate", q)
	orderData := []*types.OrderData{}
	err := db.Aggregate(dao.dbName, dao.collectionName, q, &orderData)
	if err != nil {
		logger.Error(err)
		return []*types.OrderData{}, err
	}

	return orderData, nil
}

type OrderMsg struct {
	AccountNonce    uint64         `json:"nonce"    gencodec:"required"`
	Quantity        *big.Int       `json:"quantity,omitempty"`
	Price           *big.Int       `json:"price,omitempty"`
	ExchangeAddress common.Address `json:"exchangeAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	BaseToken       common.Address `json:"baseToken,omitempty"`
	QuoteToken      common.Address `json:"quoteToken,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	PairName        string         `json:"pairName,omitempty"`
	OrderID         uint64         `json:"orderid,omitempty"`
	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash common.Hash `json:"hash" rlp:"-"`
}

// AddNewOrder add order
func (dao *OrderDao) AddNewOrder(o *types.Order, topic string) error {
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

	msg := OrderMsg{
		AccountNonce:    uint64(n),
		Quantity:        o.Amount,
		Price:           o.PricePoint,
		ExchangeAddress: o.ExchangeAddress,
		UserAddress:     o.UserAddress,
		BaseToken:       o.BaseToken,
		QuoteToken:      o.QuoteToken,
		Status:          "NEW",
		Side:            o.Side,
		Type:            o.Type,
		Hash:            o.Hash,
		PairName:        o.PairName,
		V:               V,
		R:               R,
		S:               S,
	}
	var result interface{}
	logger.Info("tomox_sendOrder", o.Status, o.Hash.Hex())
	err = rpcClient.Call(&result, "tomox_sendOrder", msg)

	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (dao *OrderDao) CancelOrder(o *types.Order, topic string) error {

	if o.Status != "CANCELLED" {
		o.Status = "CANCELLED"
	}

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

	msg := OrderMsg{
		AccountNonce:    uint64(n),
		Quantity:        o.Amount,
		Price:           o.PricePoint,
		ExchangeAddress: o.ExchangeAddress,
		UserAddress:     o.UserAddress,
		BaseToken:       o.BaseToken,
		QuoteToken:      o.QuoteToken,
		Status:          o.Status,
		Side:            o.Side,
		Type:            o.Type,
		Hash:            o.Hash,
		PairName:        o.PairName,
		OrderID:         o.OrderID,
		V:               V,
		R:               R,
		S:               S,
	}
	var result interface{}
	logger.Info("tomox_sendOrder", o.Status, o.Hash.Hex())
	err = rpcClient.Call(&result, "tomox_sendOrder", msg)

	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (dao *OrderDao) AddTopic(t []string) (string, error) {
	rpcClient, err := rpc.DialHTTP(app.Config.Tomochain["http_url"])

	defer rpcClient.Close()

	if err != nil {
		logger.Error(err)
		return "", err
	}

	var result string
	params := make(map[string]interface{})
	params["topics"] = t

	if err != nil {
		logger.Error(err)
		return "", err
	}

	err = rpcClient.Call(&result, "tomox_newTopic", params)

	if err != nil {
		logger.Error(err)
		return "", err
	}

	return result, nil
}

func (dao *OrderDao) DeleteTopic(t string) error {
	rpcClient, err := rpc.DialHTTP(app.Config.Tomochain["http_url"])

	defer rpcClient.Close()

	if err != nil {
		logger.Error(err)
		return err
	}

	var result interface{}
	params := t

	if err != nil {
		logger.Error(err)
		return err
	}

	err = rpcClient.Call(&result, "tomox_deleteTopic", params)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetOrderNonce get nonce of order
func (dao *OrderDao) GetOrderNonce(userAddress common.Address) (interface{}, error) {
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

	err = rpcClient.Call(&result, "tomox_getOrderCount", userAddress)

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info("OrderNonce:", result)
	return result, nil
}
