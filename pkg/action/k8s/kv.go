package k8s

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"knativetest/pkg/action"
	"knativetest/pkg/client"
)

type KvAction struct {
	action.Action
	Key   string
	Value string
}

func (c *KvAction) Process(ctx context.Context) interface{} {
	result, err := client.GetClient().EtcdClient.KV.Put(context.TODO(), c.Key, c.Value)
	if err != nil {
		fmt.Printf("result: %v, error info: %v", result, err)
	}
	var res, err1 = client.GetClient().EtcdClient.KV.Get(context.TODO(), c.Key, clientv3.WithPrefix())
	if err1 != nil {
		fmt.Printf("result: %v, error info: %v", res, err1)
	}
	return res
}
