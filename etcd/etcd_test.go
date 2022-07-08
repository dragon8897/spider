package etcd

import "testing"

func TestGet(t *testing.T) {
	v, err := Get("greeting")
	if err != nil {
		t.Errorf("get : %v", err)
	}
	println(string(v))
}
