package tconf

import (
	"bytes"
	"io"
	"strings"

	"github.com/spf13/viper"
	"github.com/toventang/tconf/backends"
)

type etcdConfigProvider struct{}

func (rc etcdConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	b, err := cm.Get(rp.Path())
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (rc etcdConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}

	quit := make(chan bool)
	quitwc := make(chan bool)
	// viperResponsCh := make(chan *viper.RemoteResponse)
	cryptoResponseCh := cm.Watch(rp.Path(), quit)
	resultCh := make(chan []byte)
	// need this function to convert the Channel response form crypt.Response to viper.Response
	go func(cr <-chan *backends.Response, vr chan<- []byte, quitwc <-chan bool, quit chan<- bool) {
		for {
			select {
			case <-quitwc:
				quit <- true
				return
			case resp := <-cr:
				vr <- resp.Value
			}
		}
	}(cryptoResponseCh, resultCh, quitwc, quit)

	return bytes.NewReader(<-resultCh), nil
}

func (rc etcdConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, nil
	}
	quit := make(chan bool)
	quitwc := make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)
	cryptoResponseCh := cm.Watch(rp.Path(), quit)
	// need this function to convert the Channel response form crypt.Response to viper.Response
	go func(cr <-chan *backends.Response, vr chan<- *viper.RemoteResponse, quitwc <-chan bool, quit chan<- bool) {
		for {
			select {
			case <-quitwc:
				quit <- true
				return
			case resp := <-cr:
				vr <- &viper.RemoteResponse{
					Error: resp.Error,
					Value: resp.Value,
				}
			}
		}
	}(cryptoResponseCh, viperResponsCh, quitwc, quit)

	return viperResponsCh, quitwc
}

func getConfigManager(rp viper.RemoteProvider) (backends.Store, error) {
	cm, err := newConfigManager(strings.Split(rp.Endpoint(), ","))
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func init() {
	viper.RemoteConfig = &etcdConfigProvider{}
}
