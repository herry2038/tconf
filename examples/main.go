package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/toventang/tconf"
	"github.com/toventang/tconf/backends/etcd"
	clientv3 "go.etcd.io/etcd/clientv3"
)

type DB struct {
	Host, User, Password string
}

func main() {
	var cluster string
	flag.StringVar(&cluster, "cluster", "192.168.50.190:2379", "")
	endpoints := strings.Split(cluster, ",")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := etcd.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	stop := make(chan bool)

	// watch configurations for prefix "/config_example/"
	dbConfig := &tconf.Response{Path: "/config_example/"}
	configer := tconf.New(tconf.Config{
		Client:   cli,
		FileName: "config.example.yaml",
	})
	configer.Fetch(ctx, dbConfig).Watch(ctx, stop).WithOnConfigChanged(func(rsp *tconf.Response) {
		log.Println("onChanged: ", string(rsp.Value))
	})
	val := configer.Get("host")
	log.Println("CONFIG: ", string(dbConfig.Value), ",host: ", val)

	// watch configurations for prefix "/conf_example/"
	conf := &tconf.Response{Path: "/conf_example/"}
	example := tconf.New(tconf.Config{
		Client:   cli,
		FileName: "config.example2.yaml",
	})
	example.Fetch(ctx, conf).Watch(ctx, stop)
	host := example.Get("cluster")
	log.Println("CONFIG: ", string(conf.Value), ",host: ", host)

	<-stop
}
