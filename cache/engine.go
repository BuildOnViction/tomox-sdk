package cache

// EngineCache struct
type EngineCache struct {
	fiatCacheClient     FiatCacheClient
	pairDataCacheClient PairDataCacheClient
}

// MarketCacheClient struct
type MarketCacheClient struct {
	fiatCacheClient     FiatCacheClient
	pairDataCacheClient PairDataCacheClient
}
