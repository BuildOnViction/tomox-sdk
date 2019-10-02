package endpoints

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/ethereum/go-ethereum/common"
    "github.com/gorilla/mux"
    "github.com/tomochain/tomox-sdk/interfaces"
    "github.com/tomochain/tomox-sdk/types"
    "github.com/tomochain/tomox-sdk/utils/httputils"
    "github.com/tomochain/tomox-sdk/ws"
)

type tradeEndpoint struct {
    tradeService interfaces.TradeService
}

// ServeTradeResource sets up the routing of trade endpoints and the corresponding handlers.
// TODO trim down to one single endpoint with the 3 following params: base, quote, address
func ServeTradeResource(
    r *mux.Router,
    tradeService interfaces.TradeService,
) {
    e := &tradeEndpoint{tradeService}
    r.HandleFunc("/api/trades", e.HandleGetTrades)
    r.HandleFunc("/api/trades/history", e.HandleGetTradesHistory)
    ws.RegisterChannel(ws.TradeChannel, e.tradeWebsocket)
}

// HandleGetTrades is responsible for getting pair's trade history requests
func (e *tradeEndpoint) HandleGetTrades(w http.ResponseWriter, r *http.Request) {
    v := r.URL.Query()
    bt := v.Get("baseToken")
    qt := v.Get("quoteToken")
    fromParam := v.Get("from")
    toParam := v.Get("to")

    pageOffset := v.Get("pageOffset")
    pageSize := v.Get("pageSize")
    sortBy := v.Get("sortBy")
    sortType := v.Get("sortType")

    sortedList := make(map[string]string)
    sortedList["time"] = "createdAt"

    var tradeSpec types.TradeSpec
    if bt != "" {
        if !common.IsHexAddress(bt) {
            httputils.WriteError(w, http.StatusBadRequest, "Invalid base token address")
            return
        } else {
            tradeSpec.BaseToken = bt
        }
    }

    if qt != "" {
        if !common.IsHexAddress(qt) {
            httputils.WriteError(w, http.StatusBadRequest, "Invalid base token address")
            return
        } else {
            tradeSpec.QuoteToken = qt
        }
    }

    if toParam != "" {
        t, _ := strconv.Atoi(toParam)
        tradeSpec.DateTo = int64(t)
    }
    if fromParam != "" {
        t, _ := strconv.Atoi(fromParam)
        tradeSpec.DateFrom = int64(t)
    }

    offset := 0
    size := types.DefaultLimit
    sortDB := []string{}
    if sortType != "asc" && sortType != "dec" {
        sortType = "asc"
    }
    if sortBy != "" {
        if val, ok := sortedList[sortBy]; ok {
            if sortType == "asc" {
                sortDB = append(sortDB, "+"+val)
            } else {
                sortDB = append(sortDB, "-"+val)
            }

        }
    }
    if pageOffset != "" {
        t, err := strconv.Atoi(pageOffset)
        if err != nil {
            httputils.WriteError(w, http.StatusBadRequest, "Invalid page offset")
            return
        }
        offset = t
    }
    if pageSize != "" {
        t, err := strconv.Atoi(pageSize)
        if err != nil {
            httputils.WriteError(w, http.StatusBadRequest, "Invalid page size")
            return
        }
        size = t
    }

    res, err := e.tradeService.GetTrades(&tradeSpec, sortDB, offset*size, size)
    if err != nil {
        logger.Error(err)
        httputils.WriteError(w, http.StatusInternalServerError, "")
        return
    }

    if res == nil {
        r := types.TradeRes{
            Total:  0,
            Trades: []*types.Trade{},
        }
        httputils.WriteJSON(w, http.StatusOK, r)
        return
    }

    httputils.WriteJSON(w, http.StatusOK, res)
}

// HandleGetTradesHistory is responsible for handling user's trade history requests
func (e *tradeEndpoint) HandleGetTradesHistory(w http.ResponseWriter, r *http.Request) {
    v := r.URL.Query()
    addr := v.Get("address")
    bt := v.Get("baseToken")
    qt := v.Get("quoteToken")
    fromParam := v.Get("from")
    toParam := v.Get("to")

    pageOffset := v.Get("pageOffset")
    pageSize := v.Get("pageSize")
    sortBy := v.Get("sortBy")
    sortType := v.Get("sortType")

    sortedList := make(map[string]string)
    sortedList["time"] = "createdAt"
    var tradeSpec types.TradeSpec

    if addr == "" {
        httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
        return
    }

    if !common.IsHexAddress(addr) {
        httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
        return
    }

    if bt != "" {
        if !common.IsHexAddress(bt) {
            httputils.WriteError(w, http.StatusBadRequest, "Invalid base token address")
            return
        } else {
            tradeSpec.BaseToken = bt
        }
    }

    if qt != "" {
        if !common.IsHexAddress(qt) {
            httputils.WriteError(w, http.StatusBadRequest, "Invalid base token address")
            return
        } else {
            tradeSpec.QuoteToken = qt
        }
    }

    if toParam != "" {
        t, _ := strconv.Atoi(toParam)
        tradeSpec.DateTo = int64(t)
    }
    if fromParam != "" {
        t, _ := strconv.Atoi(fromParam)
        tradeSpec.DateFrom = int64(t)
    }

    offset := 0
    size := types.DefaultLimit
    sortDB := []string{}
    if sortType != "asc" && sortType != "dec" {
        sortType = "asc"
    }
    if sortBy != "" {
        if val, ok := sortedList[sortBy]; ok {
            if sortType == "asc" {
                sortDB = append(sortDB, "+"+val)
            } else {
                sortDB = append(sortDB, "-"+val)
            }

        }
    }
    if pageOffset != "" {
        t, err := strconv.Atoi(pageOffset)
        if err != nil {
            httputils.WriteError(w, http.StatusBadRequest, "Invalid page offset")
            return
        }
        offset = t
    }
    if pageSize != "" {
        t, err := strconv.Atoi(pageSize)
        if err != nil {
            httputils.WriteError(w, http.StatusBadRequest, "Invalid page size")
            return
        }
        size = t
    }

    res, err := e.tradeService.GetTradesUserHistory(common.HexToAddress(addr), &tradeSpec, sortDB, offset*size, size)
    if err != nil {
        logger.Error(err)
        httputils.WriteError(w, http.StatusInternalServerError, "")
        return
    }

    if res == nil {
        httputils.WriteJSON(w, http.StatusOK, []types.Trade{})
        return
    }

    httputils.WriteJSON(w, http.StatusOK, res)

}

func (e *tradeEndpoint) tradeWebsocket(input interface{}, c *ws.Client) {
    b, _ := json.Marshal(input)
    var ev *types.WebsocketEvent
    if err := json.Unmarshal(b, &ev); err != nil {
        logger.Error(err)
        return
    }

    socket := ws.GetTradeSocket()

    if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
        logger.Info("Event Type", ev.Type)
        err := map[string]string{"Message": "Invalid payload"}
        socket.SendErrorMessage(c, err)
        return
    }

    b, _ = json.Marshal(ev.Payload)
    var p *types.SubscriptionPayload
    err := json.Unmarshal(b, &p)
    if err != nil {
        logger.Error(err)
        return
    }

    if ev.Type == types.SUBSCRIBE {
        if (p.BaseToken == common.Address{}) {
            err := map[string]string{"Message": "Invalid base token"}
            socket.SendErrorMessage(c, err)
            return
        }

        if (p.QuoteToken == common.Address{}) {
            err := map[string]string{"Message": "Invalid quote token"}
            socket.SendErrorMessage(c, err)
            return
        }

        e.tradeService.Subscribe(c, p.BaseToken, p.QuoteToken)
    }

    if ev.Type == types.UNSUBSCRIBE {
        if p == nil {
            e.tradeService.Unsubscribe(c)
            return
        }

        e.tradeService.UnsubscribeChannel(c, p.BaseToken, p.QuoteToken)
    }
}
