package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrKeyNotFound            = errors.New("key不存在")
	ErrCodeSendTooMany        = errors.New("发送太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnkownForCode          = errors.New("我也不知道发生了什么，反正跟Code有关")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputcode string) (bool, error)
}
type LocalCodeCache struct {
	lock       sync.RWMutex
	localcache map[string]CodeInfo
}

type CodeInfo struct {
	Code string
	Cnt  int64
	TTL  int64
}

func NewLocalCodeCache(lock sync.RWMutex, localcache map[string]CodeInfo) CodeCache {
	return &LocalCodeCache{
		lock:       lock,
		localcache: localcache,
	}
}
func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	keys := l.key(biz, phone)
	//首先 key不存在
	val, ok := l.localcache[keys]
	subttl := time.Now().Unix() - val.TTL
	//代表没有这个值
	if !ok || subttl < 540 {
		codeinfo := CodeInfo{
			Code: code,
			Cnt:  3,
			TTL:  time.Now().Unix(),
		}
		l.localcache[keys] = codeinfo
		return nil
	}
	if val.TTL == 0 {
		return errors.New("系统错误")
	}

	//发送太频繁
	return ErrCodeSendTooMany
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputcode string) (bool, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	keys := l.key(biz, phone)
	val, ok := l.localcache[keys]
	if !ok {
		return false, ErrKeyNotFound
	}
	if val.Cnt <= 0 {
		// 正常来说，如果频繁出现这个错误，你就要告警，因为有人在搞你
		return false, ErrCodeVerifyTooManyTimes
	}
	if val.Code == inputcode {
		l.localcache[keys] = CodeInfo{}
		return true, nil
	}
	val.Cnt--
	return false, nil
}

func (c *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
