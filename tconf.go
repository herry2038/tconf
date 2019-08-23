package tconf

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type TConf struct {
	configer *viper.Viper
}
type Config struct {
	Path, FileName, Cluster string
	Renew                   bool
	ClusterKey              string
}

// New initialize TConf
func New(conf *Config) *TConf {
	// First start
	fstConfiger := newConf(conf.Path, conf.FileName, conf.Cluster)
	if conf.Renew {
		// get the new clusters
		result, err := fstConfiger.Get(conf.ClusterKey)
		if err != nil {
			log.Fatalln(err)
		}

		defer func() {
			fstConfiger = nil
		}()

		configer := newConf(conf.Path, conf.FileName, result.(string))
		return configer
	}
	return fstConfiger
}

func newConf(path, filename, cluster string) *TConf {
	var configer = viper.New()
	var err error

	if filename == "" {
		panic("Filename need a value")
	}
	configer.SetConfigFile(filename)

	if cluster == "" {
		err = fmt.Errorf("No cluster is set")
	}

	if err == nil {
		log.Println("Used etcd clusters: ", cluster)

		err = fetchRemoteConfig(configer, path, cluster)
	}
	if err != nil {
		log.Println("Used local file to config: ", filename)
		// read local file when k/v store error occurred
		err = configer.ReadInConfig()
		if err != nil {
			log.Panicln(err)
		}
	}

	return &TConf{configer}
}

func fetchRemoteConfig(configer *viper.Viper, path, cluster string) error {
	configer.AddRemoteProvider("etcd", cluster, path)
	err := configer.ReadRemoteConfig()
	if err != nil {
		return err
	}
	configer.WriteConfig()

	return nil
}

// Get value
func (c *TConf) Get(key string) (interface{}, error) {
	return c.configer.Get(key), nil
}

// Unmarshal unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func (c *TConf) Unmarshal(in interface{}) {
	c.configer.Unmarshal(in)
}

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func (c *TConf) UnmarshalKey(key string, in interface{}) {
	c.configer.UnmarshalKey(key, in)
}

// Watch watching for remote server
func (c *TConf) Watch() error {
	err := c.configer.WatchRemoteConfig()
	if err == nil {
		c.configer.WriteConfig()
	}
	return err
}
