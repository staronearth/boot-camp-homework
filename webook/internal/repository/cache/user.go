package cache

import (
	"boot-camp-homework/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

//	type Cache interface {
//		GetUser(ctx *context.Context, id int64) (domain.User, error)
//		//读取文章
//		GetArticle(ctx *context.Context, aid int64)
//
//		//还有别的业务
//		//。。。
//	}
//
//	type CacheV1 interface {
//		//你的中间件团队去做
//		Get(ctx context.Context, key string) (any, error)
//	}
var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}
type RedisUserCache struct {
	//传单机redis可以
	//传cluster的redis也可以
	client     redis.Cmdable
	expiration time.Duration
}

// NewUserCache A用到了B,B一定是接口;A用到了B, B一定是A的字段；A用到了B,A绝对不初始化B，而是外面注入
func NewUserCache(client redis.Cmdable, expiration time.Duration) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: expiration,
	}
}

// Get 只要error为nil,就认为User一定在
// 如果没有数据返回一个特定的error
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.Key(id)
	res, err := cache.client.Get(ctx, key).Bytes()
	var u domain.User
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(res, &u)
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.Key(user.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
