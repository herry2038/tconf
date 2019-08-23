package backends

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var endpoints = []string{"192.168.50.190:2379"}

var cli *EtcdClient

func TestMain(t *testing.M) {
	cli, _ = New(endpoints, 5*time.Second)

	os.Exit(t.Run())
}

func TestGet(t *testing.T) {
	defer cli.client.Close()

	value, err := cli.Get("/conf_example/")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("get value: %v\n", string(value))
}

func TestWatch(t *testing.T) {
	stop := make(chan bool)
	resp := cli.Watch("/conf_example/", stop)

	go func() {
		for {
			select {
			case <-stop:
				return
			case r := <-resp:
				if r.Error != nil {
					t.Error(r.Error)
					continue
				}
				t.Log(resp)
			}
		}
	}()
}
