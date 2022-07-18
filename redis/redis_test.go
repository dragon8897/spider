package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(GetOptions())

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		t.Errorf("%+v", err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		t.Errorf("%+v", err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		t.Errorf("%+v", err)
	} else {
		fmt.Println("key2", val2)
	}
}
