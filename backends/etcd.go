package backends

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/mvcc/mvccpb"

	etcd "go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	client *etcd.Client
}

func New(endpoints []string, dialTimeout time.Duration) (*EtcdClient, error) {
	client, err := etcd.New(etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("Connect to etcd server for tconf: %v", err)
	}
	return &EtcdClient{client: client}, nil
}

// Get retrieves a value from a K/V store for the provided key.
func (c *EtcdClient) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 50*time.Second)
	resp, err := c.client.Get(ctx, key)
	defer cancel()
	if err != nil {
		return nil, fmt.Errorf(`Get key "%s" from etcd server: %v`, key, err)
	}
	if len(resp.Kvs) > 0 {
		return resp.Kvs[0].Value, nil
	}
	return nil, fmt.Errorf(`Get key "%s" not found`, key)
}

// Watch monitors a K/V store for changes to key.
func (c *EtcdClient) Watch(key string, stop chan bool) <-chan *Response {
	respChan := make(chan *Response, 0)
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-stop
			cancel()
		}()

		rch := c.client.Watch(ctx, key, etcd.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				respChan <- &Response{[]byte(ev.Kv.Value), nil}
			}
		}
	}()
	return respChan
}

func convertToKVPairs(kvs []*mvccpb.KeyValue) []*KVPair {
	var nodes []*KVPair

	for _, kv := range kvs {
		nodes = append(nodes, &KVPair{Key: string(kv.Key), Value: kv.Value})
	}
	return nodes
}
