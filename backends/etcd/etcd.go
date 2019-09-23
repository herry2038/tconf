package etcd

import (
	"context"
	"fmt"

	"github.com/toventang/tconf/backends"
	clientv3 "go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	client *clientv3.Client
}

func New(config clientv3.Config) (*EtcdClient, error) {
	client, err := clientv3.New(config)
	if err != nil {
		return nil, fmt.Errorf("Connect to etcd server for tconf: %v", err)
	}
	return &EtcdClient{client: client}, nil
}

// Get retrieves a value from a K/V store for the provided key.
func (c *EtcdClient) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf(`Get key "%s" from etcd server: %v`, key, err)
	}
	if len(resp.Kvs) > 0 {
		return resp.Kvs[0].Value, nil
	}
	return nil, fmt.Errorf(`Get key "%s" not found`, key)
}

// Watch a K/V store for changes to key.
func (c *EtcdClient) Watch(ctx context.Context, key string, stop chan bool) <-chan *backends.Response {
	respChan := make(chan *backends.Response, 0)
	go func() {
		go func() {
			<-stop
		}()

		rch := c.client.Watch(ctx, key, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				respChan <- &backends.Response{Value: []byte(ev.Kv.Value), Error: nil}
			}
		}
	}()
	return respChan
}

func (c *EtcdClient) Close() error {
	return c.client.Close()
}
