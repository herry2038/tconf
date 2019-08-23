package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/toventang/tconf"
)

func main() {
	var path, filename, cluster, clusterKey string
	var renew bool
	flag.StringVar(&path, "path", "/conf_example/", "a key for k/v store")
	flag.StringVar(&filename, "filename", "../config.example.yaml", "local file cache for all configuration sections")
	flag.StringVar(&cluster, "cluster", "192.168.50.190:2379", "etcd cluster, mutil cluster can be use with ','")
	flag.BoolVar(&renew, "renew", true, "create a new TConf with the new node when the k/v store endpoints are inconsistent with the local")
	flag.StringVar(&clusterKey, "key", "cluster", "a endpoints config section in remote k/v store")
	flag.Parse()

	configer := tconf.New(&tconf.Config{
		Path:       path,
		FileName:   filename,
		Cluster:    cluster,
		Renew:      renew,
		ClusterKey: clusterKey,
	})
	val, err := configer.Get("db")
	if err != nil {
		log.Fatal(`get "db" error: `, err)
	}
	fmt.Println(`db: `, val)

	db := &DBConfig{}
	configer.UnmarshalKey("db", db)
	fmt.Println(db)

	// If you want to get those configurations when changed each time, you can use the Watch() method
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			err := configer.Watch()
			if err != nil {
				wg.Done()
			}
		}
	}()
	fmt.Println(strings.Repeat("=", 37))
	wg.Wait()
}

// DBConfig present a database config
type DBConfig struct {
	Database string   `json:"database"`
	Host     []string `json:"host"`
	Port     uint16   `json:"port"`
	UserName string   `json:"username"`
	Password string   `json:"password"`
}
