package tconf

import (
	"time"

	"github.com/toventang/tconf/backends"
)

type configManager struct {
	store backends.Store
}

func newConfigManager(machines []string) (backends.Store, error) {
	store, err := backends.New(machines, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &configManager{store}, nil
}

func (c *configManager) Get(key string) ([]byte, error) {
	value, err := c.store.Get(key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *configManager) Watch(key string, stop chan bool) <-chan *backends.Response {
	resp := make(chan *backends.Response, 0)
	backendResp := c.store.Watch(key, stop)
	go func() {
		for {
			select {
			case <-stop:
				return
			case r := <-backendResp:
				if r.Error != nil {
					resp <- &backends.Response{nil, r.Error}
					continue
				}
				resp <- &backends.Response{r.Value, nil}
			}
		}
	}()
	return resp
}
