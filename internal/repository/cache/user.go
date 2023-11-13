package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/webook/internal/domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
	Del(ctx context.Context, id int64) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd: cmd,
		//key expire at 15 minutes
		expiration: time.Minute * 15,
	}
}

func (cache *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := cache.key(uid)
	// assume using Json to serialize the returned data
	data, err := cache.cmd.Get(ctx, key).Result() // data is a string
	//data, err := cache.cmd.Get(ctx, firstKey).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := cache.key(du.Id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, key, data, cache.expiration).Err()
}

func (cache *RedisUserCache) Del(ctx context.Context, id int64) error {
	return cache.cmd.Del(ctx, cache.key(id)).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
