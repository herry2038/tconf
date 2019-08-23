package tconf

import (
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

// New initialize TConf.
func New(conf *Config) *TConf {
	// First start
	fstConfiger := newConf(conf.Path, conf.FileName, conf.Cluster)
	if conf.Renew {
		// get the new clusters
		result, err := fstConfiger.Get(conf.ClusterKey)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Clusters: ", result)

		defer func() {
			fstConfiger = nil
		}()

		configer := newConf(conf.Path, conf.FileName, result.(string))
		// configer.configer.ReadInConfig()
		return configer
	}
	// fstConfiger.configer.ReadInConfig()
	return fstConfiger
}

func newConf(path, filename, cluster string) *TConf {
	var configer = viper.New()

	if filename == "" {
		panic("Filename need a value")
	}
	configer.SetConfigFile(filename)

	if cluster != "" {
		log.Println("Used etcd")

		fetchRemoteConfig(configer, path, cluster)
	}

	return &TConf{configer}
}

func fetchRemoteConfig(configer *viper.Viper, path, cluster string) {
	configer.AddRemoteProvider("etcd", cluster, path)
	err := configer.ReadRemoteConfig()
	if err != nil {
		log.Fatalln("Read remote config error: ", err)
	}
	configer.WriteConfig()
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
