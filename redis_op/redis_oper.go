package redis_op

import (
	"context"
	"encoding/json"
	"github.com/dapr-platform/common"
	"github.com/redis/go-redis/v9"
	"time"
)

var Rdb *redis.Client

// 使用单独的redis client，不使用dapr的statestore. 是为了使用redis多种数据结构。
func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "things2024", // password set
		DB:       0,            // use default DB
	})
	common.Logger.Debug("Rdb inited")

}

func GetRedisVal[T any](ctx context.Context, key string) (t T, err error) {
	val, err := Rdb.Get(ctx, key).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(val), &t)
	return
}

func SetRedisVal(ctx context.Context, key string, val any, expiredSeconds int) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return Rdb.Set(ctx, key, string(b), time.Second*time.Duration(expiredSeconds)).Err()
}
