package ioc

import (
	"boot-camp-homework/webook/internal/repository/cache"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

func InitUserCache(client redis.Cmdable) cache.UserCache {
	return cache.NewUserCache(client, time.Minute*15)
}

func InitLocalCodeCache() cache.CodeCache {
	lock := sync.RWMutex{}
	localmap := make(map[string]cache.CodeInfo, 10)
	return cache.NewLocalCodeCache(lock, localmap)
}

func InitCodeCache(client redis.Cmdable) cache.CodeCache {
	return cache.NewCodeCache(client)
}
