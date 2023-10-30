package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	//"github.com/redis/go-redis/v9"
	"github.com/coocood/freecache"
)

var (
	////go:embed lua/set_code.lua
	//luaSetCode string
	////go:embed lua/verify_code.lua
	//luaVerifyCode string

	ErrCodeSendTooMany   = errors.New("Code Sent Too Frequent")
	ErrCodeVerifyTooMany = errors.New("Code Verified Too Frequent")

	cacheSizeMB   int = 100
	expireSeconds int = 60
)

// In Memory Verification Code Cache
type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type FreecacheCodeCache struct {
	cache *freecache.Cache
}

func NewCodeCache(cache freecache.Cache) CodeCache {
	return &FreecacheCodeCache{
		cache: freecache.NewCache(cacheSizeMB * 1024 * 1024),
	}
}

func (c *FreecacheCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	key := c.key(biz, phone)
	got, _ := c.cache.Get([]byte(key))
	if got != nil {
		lastSent := time.Unix(0, c.cache.Touch([]byte(key), expireSeconds))
		if time.Since(lastSent) < time.Minute {
			return ErrCodeSendTooMany
		}
	}

	err := c.cache.Set([]byte(key), []byte(code), expireSeconds)
	if err != nil {
		return errors.New("code exist, but no expiration time")
	}
	return nil
}

func (c *FreecacheCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	key := c.key(biz, phone)
	got, _ := c.cache.Get([]byte(key))
	if got == nil {
		return false, errors.New("code doesn't exist")
	}

	if string(got) != code {
		return false, errors.New("code doesn't match")
	}

	lastVerified := time.Unix(0, c.cache.Touch([]byte(key), expireSeconds))
	if time.Since(lastVerified) < time.Minute {
		return false, ErrCodeVerifyTooMany
	}

	c.cache.Set([]byte(key), []byte(code), expireSeconds)
	return true, nil
}

func (c *FreecacheCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

// Redis Verification Code Cache

//type CodeRedisCache interface {
//	Set(ctx context.Context, biz, phone, code string) error
//	Verify(ctx context.Context, biz, phone, code string) (bool, error)
//}
//
//type RedisCodeCache struct {
//	cmd redis.Cmdable
//}
//
//func NewCodeCache(cmd redis.Cmdable) CodeRedisCache {
//	return &RedisCodeCache{
//		cmd: cmd,
//	}
//}
//
//func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
//	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
//	// print log
//	if err != nil {
//		// issue with redis
//		return err
//	}
//	switch res {
//	case -2:
//		return errors.New("code exist, but no expiration time")
//	case -1:
//		return ErrCodeSendTooMany
//	default:
//		return nil
//	}
//}
//
//func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
//	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
//	if err != nil {
//		// issue with redis
//		return false, err
//	}
//	switch res {
//	case -2:
//		return false, nil
//	case -1:
//		return false, ErrCodeVerifyTooMany
//	default:
//		return true, nil
//	}
//}
//
//func (c *RedisCodeCache) key(biz, phone string) string {
//	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
//}
