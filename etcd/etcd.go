package etcd

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Get(key string) (value []byte, err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"d.sduang.top:19379"},
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "xxxxxx",
	})
	defer cli.Close()
	if err != nil {
		return
	}
	resp, err := cli.Get(context.Background(), key)
	if err != nil {
		return
	}
	for _, v := range resp.Kvs {
		value = v.Value
	}
	return
}
