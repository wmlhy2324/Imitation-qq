package core

import (
	"context"
	"fmt"
	"testing"
)

func TestEtcd(t *testing.T) {
	client := InitEtcd("127.0.0.1:2379")
	res, err := client.Put(context.Background(), "auth", "127.0.0.1")
	fmt.Println(res, err)
	getResponse, err := client.Get(context.Background(), "auth")
	if err == nil && len(getResponse.Kvs) > 0 {
		fmt.Println(string(getResponse.Kvs[0].Value))
	}
}
