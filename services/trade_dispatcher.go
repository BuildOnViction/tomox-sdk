package services

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
)

// TradeDispatcherService depatch trade for services
type TradeDispatcherService struct {
	tradeDao                 interfaces.TradeDao
	dispatcherdb             dispatcherdb
	tradeFetchNotifyCallback []func(*types.Trade)
	tradeNotifyCallback      []func(*types.TradeChangeEvent)
	mutex                    sync.RWMutex
}

type timeframe struct {
	FirstTime int64 `json:"firstTime" bson:"firstTime"`
	LastTime  int64 `json:"lastTime" bson:"lastTime"`
}
type timeframes []*timeframe

type dispatcherdb struct {
	tframes timeframes
}

// NewTradeDispatcherService init new ohlcv service
func NewTradeDispatcherService(tradeDao interfaces.TradeDao) *TradeDispatcherService {
	return &TradeDispatcherService{
		tradeDao: tradeDao,
	}
}

// UnsubscribeTrade handles all the unsubscription messages
func (s *TradeDispatcherService) UnsubscribeTrade(fn func(ev *types.TradeChangeEvent)) {

}

// SubscribeTrade handles all the subscription messages
func (s *TradeDispatcherService) SubscribeTrade(fn func(ev *types.TradeChangeEvent)) {
	if fn != nil {
		s.tradeNotifyCallback = append(s.tradeNotifyCallback, fn)
	}
}

// UnsubscribeFetch handles all the unsubscription messages
func (s *TradeDispatcherService) UnsubscribeFetch(fn func(*types.Trade)) {
}

// SubscribeFetch handles all the subscription messages
func (s *TradeDispatcherService) SubscribeFetch(fn func(*types.Trade)) {
	if fn != nil {
		s.tradeFetchNotifyCallback = append(s.tradeFetchNotifyCallback, fn)
	}
}

// Start init cache
// ensure add current time frame before trade notify come
func (s *TradeDispatcherService) Start() {
	now := time.Now().Unix()
	datefrom := now - intervalMin
	s.loadDatabase()
	lastFrame := s.lastTimeFrame()
	if lastFrame != nil {
		logger.Info("last frame first time", time.Unix(lastFrame.FirstTime, 0))
		if now-lastFrame.LastTime < intervalMin {
			datefrom = lastFrame.LastTime
		}
	} else {
		// add start frame to list
		s.dispatcherdb.tframes = append(s.dispatcherdb.tframes, &timeframe{
			FirstTime: now - intervalMax,
			LastTime:  now - intervalMax,
		})
	}
	// add current frame to list
	s.dispatcherdb.tframes = append(s.dispatcherdb.tframes, &timeframe{
		FirstTime: now,
		LastTime:  now,
	})

	lastFrame = s.lastTimeFrame()
	logger.Info("init fetch", time.Unix(datefrom, 0), time.Unix(now, 0))
	s.fetch(datefrom, now, lastFrame)
	s.commitDatabase()
	go s.continueFetch()
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.commitDatabase()
				if err != nil {
					logger.Error(err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	go s.watchChanges()
	logger.Info("OHLCV finished")
}

func (s *TradeDispatcherService) fetch(fromdate int64, todate int64, frame *timeframe) {

	pageOffset := 0
	size := 1000
	for {
		trades, err := s.tradeDao.GetTradeByTime(fromdate, todate, pageOffset*size, size)
		logger.Debug("FETCH DATA", pageOffset*size)
		if err != nil || len(trades) == 0 {
			break
		}
		sort.Slice(trades, func(i, j int) bool {
			return trades[i].CreatedAt.Unix() < trades[j].CreatedAt.Unix()
		})
		s.mutex.Lock()
		for i, trade := range trades {
			for _, fn := range s.tradeFetchNotifyCallback {
				fn(trade)
			}
			if i == 0 {
				s.updatefisttimeframe(trade.CreatedAt.Unix(), frame)
			}

		}
		s.mutex.Unlock()
		pageOffset = pageOffset + 1
	}
}

// ensure init cache finished before invoke
func (s *TradeDispatcherService) continueFetch() {
	if len(s.dispatcherdb.tframes) > 1 {
		for i := len(s.dispatcherdb.tframes) - 1; i > 0; i-- {
			currentframe := s.dispatcherdb.tframes[i]
			preframe := s.dispatcherdb.tframes[i-1]
			if currentframe.FirstTime > preframe.LastTime {
				logger.Debug("continue fetch", time.Unix(preframe.LastTime, 0), time.Unix(currentframe.FirstTime, 0))
				s.fetch(preframe.LastTime, currentframe.FirstTime, currentframe)
			}

		}
	}
	logger.Debug("continueFetch finished")
}
func (s *TradeDispatcherService) lastTimeFrame() *timeframe {
	if len(s.dispatcherdb.tframes) > 0 {
		return s.dispatcherdb.tframes[len(s.dispatcherdb.tframes)-1]
	}
	return nil
}
func (s *TradeDispatcherService) updatefisttimeframe(firsttime int64, frame *timeframe) {
	logger.Info("updatefisttimeframe", time.Unix(firsttime, 0))
	if frame != nil {
		frame.FirstTime = firsttime
	}
}

func (s *TradeDispatcherService) updatelasttimeframe(lasttime int64, frame *timeframe) {
	if frame != nil {
		frame.LastTime = lasttime
	}

}

func (s *TradeDispatcherService) commitDatabase() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	data, err := json.Marshal(s.dispatcherdb)
	if err != nil {
		return err
	}
	file, err := os.Create("dispatcher.db")
	defer file.Close()
	if err == nil {
		_, err = file.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TradeDispatcherService) loadDatabase() error {
	file, err := os.Open("dispatcher.db")
	defer file.Close()
	if err != nil {
		return err
	}
	stats, statsErr := file.Stat()
	if statsErr != nil {
		return statsErr
	}

	size := stats.Size()
	bytes := make([]byte, size)
	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	var db dispatcherdb
	err = json.Unmarshal(bytes, &db)
	if err != nil {
		return err
	}
	s.dispatcherdb.tframes = db.tframes
	return nil
}

// WatchChanges watch record
func (s *TradeDispatcherService) watchChanges() {

	ct, sc, err := s.tradeDao.Watch()

	if err != nil {
		logger.Error("Failed to open change stream")
		return //exiting func
	}

	defer ct.Close()
	defer sc.Close()

	// Watch the event again in case there is error and function returned
	defer s.watchChanges()

	ctx := context.Background()

	//Handling change stream in a cycle
	for {
		select {
		case <-ctx.Done(): // if parent context was cancelled
			err := ct.Close() // close the stream
			if err != nil {
				logger.Error("Change stream closed")
			}
			return //exiting from the func
		default:
			ev := types.TradeChangeEvent{}

			//getting next item from the steam
			ok := ct.Next(&ev)

			//if data from the stream wasn't un-marshaled, we get ok == false as a result
			//so we need to call Err() method to get info why
			//it'll be nil if we just have no data
			if !ok {
				err := ct.Err()
				if err != nil {
					logger.Error(err)
					return
				}
			}

			//if item from the stream un-marshaled successfully, do something with it
			if ok {
				logger.Debugf("Operation Type: %s", ev.OperationType)
				for _, fn := range s.tradeNotifyCallback {

					if ev.OperationType == types.OPERATION_TYPE_INSERT {
						lastFrame := s.lastTimeFrame()
						s.updatelasttimeframe(ev.FullDocument.CreatedAt.Unix(), lastFrame)
					}
					fn(&ev)
				}
			}
		}
	}
}
