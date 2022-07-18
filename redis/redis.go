package redis

import (
	"github.com/go-redis/redis/v8"
)

func GetOptions() *redis.Options {
	return &redis.Options{
		Addr: "d.sduang.top:19679",
	}
}
