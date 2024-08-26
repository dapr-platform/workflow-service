package redis_op

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestZadd(t *testing.T) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := Rdb.ZAdd(context.Background(), "workflow:测试:TimeDelay:65d1341b-5cd6-4032-8bdc-d19a4c8ab10f", redis.Z{1, "1"}, redis.Z{2, "2"}, redis.Z{3, "3"}).Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val)
}
