package flow_dsl

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestFlowLoop(t *testing.T) {
	var cacheData = &SecondInterval{
		Second: time.Now().Unix(),
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "demo.lq-tec.cn:36379",
		Password: "things2023", // no password set
		DB:       0,            // use default DB
	})
	err1 := rdb.Set(context.Background(), "test", cacheData, time.Duration(60)*time.Second).Err()
	if err1 != nil {
		panic(err1)
	}

	var newData SecondInterval
	err2 := rdb.Get(context.Background(), "test").Scan(&newData)
	if err2 != nil {
		panic(err2)
	}

}
