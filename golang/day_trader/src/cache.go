package main

import (
	"github.com/rainycape/memcache"
)

var cache *memcache.Client

func initCache() {
	cache, _ = memcache.New("daytrader_cache:11211")
}

func (user *User) setCache() error {
	bytes, _ := user.MarshalJSON()
	item := &memcache.Item{
		Key:   user.Id,
		Value: bytes,
	}
	return cache.Set(item)
}

func (stock *Stock) setCache() error {
	bytes, _ := stock.MarshalJSON()
	item := &memcache.Item{
		Key:   stock.Symbol,
		Value: bytes,
	}
	return cache.Set(item)
}

func getCacheStock(key string) (*Stock, error) {
	stock := &Stock{}
	item, err := cache.Get(key)
	if err != nil {
		return nil, err
	}
	err = stock.UnmarshalJSON(item.Value)
	return stock, err
}

func getCacheUser(key string) (*User, error) {
	user := &User{Id: key}
	item, err := cache.Get(key)
	if err != nil {
		return user, err
	}
	err = user.UnmarshalJSON(item.Value)
	if user.StockMap == nil {
		user.StockMap = make(map[string]int)
	}
	return user, err
}
