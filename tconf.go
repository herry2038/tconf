package tconf

import (
	"context"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/toventang/tconf/backends"
)

type TConf struct {
	client             backends.Store
	kvstore            map[string]interface{}
	current            *Response
	onChanged          func(*Response)
	configFileProvider ConfigFileProvider
	filename           string
}

// Config configurations
type Config struct {
	// k/v store center
	Client backends.Store
	// local file name.
	FileName string
	// config file processor
	ConfigFileProvider ConfigFileProvider
}

// Response represents the configurations from the k/v store center
type Response struct {
	Path  string
	Value []byte
}

// New creates a new TConf using config params
func New(config Config) *TConf {
	return &TConf{
		client:             config.Client,
		configFileProvider: config.ConfigFileProvider,
		filename:           config.FileName,
	}
}

// Get gets configurations for key
func (c *TConf) Get(key string) interface{} {
	if c.current == nil || c.configFileProvider == nil {
		return nil
	}
	if c.kvstore == nil || len(c.kvstore) == 0 {
		kvs, err := c.configFileProvider.Unmarshal(c.current.Value)
		if err != nil {
			c.kvstore = nil
			return nil
		}
		c.kvstore = kvs
	}

	return c.kvstore[key]
}

// Fetch gets configurations for Path
func (c *TConf) Fetch(ctx context.Context, val *Response) *TConf {
	if c.current != nil && c.current.Value != nil {
		val.Path = c.current.Path
		val.Value = c.current.Value
		return c
	}
	if val.Path == "" {
		panic("path must be a string value")
	}

	var v []byte
	var err error
	if c.client != nil {
		v, err = c.client.Get(ctx, val.Path)
	}
	if err != nil {
		// as an error occured when getting the configuration from the K/V Store Center,
		// read the configuration from the local file
		v, err = c.configFileProvider.ReadConfig()
	}

	val.Value = v
	c.current = val

	c.saveConfig()

	return c
}

func (c *TConf) saveConfig() {
	if c.configFileProvider == nil {
		p := &YAMLProvider{FileName: c.filename}
		c.configFileProvider = p
	}
	m, err := c.configFileProvider.Unmarshal(c.current.Value)
	if err != nil {
		log.Fatalln("an error occured when save config to file, ", err.Error())
	} else {
		c.kvstore = m
		go c.configFileProvider.WriteConfig(m)
	}
}

func (c *TConf) WithOnConfigChanged(onChanged func(*Response)) *TConf {
	c.onChanged = onChanged
	return c
}

// Watch config changes
func (c *TConf) Watch(ctx context.Context, stop chan bool) *TConf {
	if c.client == nil {
		return c
	}
	c.watch(ctx, c.current.Path, stop)
	return c
}

func (c *TConf) watch(ctx context.Context, path string, stop chan bool) {
	//在这里监听配置中心的修改，并实时保存到本地配置文件，同时更新缓存，以免客户端仍然使用过期的配置
	response := c.client.Watch(ctx, path, stop)
	go func() {
		for {
			select {
			case <-stop:
				return
			case r := <-response:
				if r.Error != nil {
					log.Fatalln("etcd watcher error: ", r.Error)
					continue
				}
				c.current.Path = path
				c.current.Value = r.Value
				go c.saveConfig()
				if c.onChanged != nil {
					go c.onChanged(&Response{
						Path:  path,
						Value: r.Value,
					})
				}
			}
		}
	}()
}

// UnmarshalAll decode the all configurations into a struct
func (c *TConf) UnmarshalAll(out interface{}) error {
	return mapstructure.Decode(c.kvstore, out)
}

// UnmarshalKey decode the config of the key into a struct
func (c *TConf) UnmarshalKey(key string, out interface{}) error {
	return mapstructure.Decode(c.kvstore[key], out)
}

// Close client connection
func (c *TConf) Close() error {
	return c.client.Close()
}
