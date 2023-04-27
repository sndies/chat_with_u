package cache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/sndies/chat_with_u/middleware/log"
	"time"
)

var cacheClient *cache.Cache

func Init() {
	cacheClient = cache.New(10*time.Minute, 5*time.Minute)
}

func Add(ctx context.Context, key string, val interface{}, expire time.Duration) error {
	err := cacheClient.Add(key, val, expire)
	if err != nil {
		log.Errorf(ctx, "cache set key: %s, err: %v", key, err)
	}
	return err
}

func Get(ctx context.Context, key string) (interface{}, bool) {
	val, exist := cacheClient.Get(key)
	log.Infof(ctx, "cache get key: %s, val: %+v", key, val)
	return val, exist
}

func Del(ctx context.Context, key string) {
	cacheClient.Delete(key)
}