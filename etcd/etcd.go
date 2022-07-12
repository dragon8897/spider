package etcd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetConfig() clientv3.Config {
	return clientv3.Config{
		Endpoints:   []string{"d.sduang.top:19379"},
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "xxxxxx",
	}
}
