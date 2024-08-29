package cache

import (
	"context"
	"time"

	"github.com/kodernubie/gocommon/conf"
	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	client *redis.Client
}

func (o *RedisConnection) Set(key, value string, expiration ...time.Duration) {

	var targetExpiration time.Duration = 0

	if len(expiration) > 0 {
		targetExpiration = expiration[0]
	}

	o.client.Set(context.Background(), key, value, targetExpiration)
}

func (o *RedisConnection) SetNX(key, value string, expiration ...time.Duration) bool {

	var targetExpiration time.Duration = 0

	if len(expiration) > 0 {
		targetExpiration = expiration[0]
	}

	ret, _ := o.client.SetNX(context.Background(), key, value, targetExpiration).Result()

	return ret
}

func (o *RedisConnection) Get(key string) string {

	ret, _ := o.client.Get(context.Background(), key).Result()

	return ret
}

func (o *RedisConnection) Incr(key string, expiration ...time.Duration) int {

	if len(expiration) > 0 {

		pipe := o.client.TxPipeline()
		cmd := pipe.Incr(context.Background(), key)
		pipe.Expire(context.Background(), key, expiration[0])
		pipe.Exec(context.Background())

		return int(cmd.Val())
	}

	ret, _ := o.client.Incr(context.Background(), key).Result()
	return int(ret)
}

func (o *RedisConnection) Decr(key string, expiration ...time.Duration) int {

	if len(expiration) > 0 {

		pipe := o.client.TxPipeline()
		cmd := pipe.Incr(context.Background(), key)
		pipe.Expire(context.Background(), key, expiration[0])
		pipe.Exec(context.Background())

		return int(cmd.Val())
	}

	ret, _ := o.client.Decr(context.Background(), key).Result()
	return int(ret)
}

func (o *RedisConnection) Del(key string) {

	o.client.Del(context.Background(), key)
}

func init() {

	RegConnCreator("redis", func(configName string) (Connection, error) {

		ret := &RedisConnection{}

		ret.client = redis.NewClient(&redis.Options{
			Addr:     conf.Str("CACHE_"+configName+"_ADDR", "localhost:6379"),
			Password: conf.Str("CACHE_"+configName+"_PASS", ""),
			DB:       conf.Int("CACHE_"+configName+"_DB", 0),
		})

		return ret, nil
	})
}
