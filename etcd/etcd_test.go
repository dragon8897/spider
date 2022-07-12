package etcd

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ETCD_TEST_KEY = "ETCD_TEST_KEY"
	ETCD_PREFIX   = "ETCD_PREFIX_"

	ETCD_LEASE_TTL = 2
)

func TestEtcd(t *testing.T) {
	cli, err := clientv3.New(GetConfig())
	if err != nil {
		t.Errorf("new err: %v", err)
	}
	defer cli.Close()
	const resultVal string = "hello"
	ctx := context.Background()
	_, err = cli.Put(ctx, ETCD_TEST_KEY, resultVal)
	if err != nil {
		t.Errorf("put err: %v", err)
	}
	resp, err := cli.Get(ctx, ETCD_TEST_KEY)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	if resp.Count != 1 {
		t.Errorf("get count != 1")
	}
	if !bytes.Equal(resp.Kvs[0].Value, []byte(resultVal)) {
		t.Errorf("get value %s != %s", resp.Kvs[0].Value, resultVal)
	}

	// 删除
	del, err := cli.Delete(ctx, ETCD_TEST_KEY)
	if err != nil {
		t.Errorf("delete err: %v", err)
	}
	if del.Deleted != 1 {
		t.Errorf("delete count != 1")
	}

	// prefix
	_, err = cli.Put(ctx, ETCD_PREFIX+ETCD_TEST_KEY, resultVal)
	if err != nil {
		t.Errorf("put err: %v", err)
	}
	resp, err = cli.Get(ctx, ETCD_PREFIX, clientv3.WithPrefix())
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	if resp.Count != 1 {
		t.Errorf("get count != 1")
	}
	if !bytes.Equal(resp.Kvs[0].Value, []byte(resultVal)) {
		t.Errorf("get value %s != %s", resp.Kvs[0].Value, resultVal)
	}

	// lease
	lease, err := cli.Grant(ctx, ETCD_LEASE_TTL)
	if err != nil {
		t.Errorf("grant err: %v", err)
	}
	_, err = cli.Put(ctx, ETCD_TEST_KEY, resultVal, clientv3.WithLease(lease.ID))
	if err != nil {
		t.Errorf("put err: %v", err)
	}
	resp, err = cli.Get(ctx, ETCD_TEST_KEY)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	for _, v := range resp.Kvs {
		fmt.Printf("key: %s value: %s\n", v.Key, v.Value)
	}
	time.Sleep((ETCD_LEASE_TTL + 1) * time.Second)
	resp, err = cli.Get(ctx, ETCD_TEST_KEY)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	if resp.Count != 0 {
		t.Errorf("get count != 0")
	}
}

func TestWatch(t *testing.T) {
	cli, err := clientv3.New(GetConfig())
	if err != nil {
		t.Errorf("new err: %v", err)
	}
	const resultVal string = "watching"
	ctx := context.Background()

	// watch
	wch := cli.Watch(ctx, ETCD_TEST_KEY)
	go func() {
		time.Sleep((ETCD_LEASE_TTL + 1) * time.Second)
		cli.Close()
	}()
	go func() {
		cli.Put(ctx, ETCD_TEST_KEY, resultVal)
		cli.Delete(ctx, ETCD_TEST_KEY)
		cli.Put(ctx, ETCD_TEST_KEY, resultVal)
	}()
	for v := range wch {
		for _, event := range v.Events {
			fmt.Printf("%v , %s, %s\n", event.Type, event.Kv.Key, event.Kv.Value)
		}
	}
}
