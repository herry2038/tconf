# tconf #

TConf 默认提供 etcdv3 的支持，并通过[接口](backends/store.go)可实现对 zookeeper、consul、etcdv2 的支持。
TConf 启动时以及开启服务器监听后，配置中心的任意改动都会被监听到，并替换本地配置。

默认支持 yaml 文件格式，提供[文件读写接口](provider.go)，可自行实现对系统环境变量(ENV)或其他文件格式的支持，如：toml、json、xml 等。
(* TConf 仅提供从配置中心读取配置，如果需要修改配置请在配置中心操作)

## 示例 ##

写入 etcd 配置，下面是配置 mongodb 的数据库连接参数。
支持多个配置中心集群（或节点）地址，以逗号 "," 分隔。

```console
etcdctl put /conf_example/ '{"cluster":"192.168.50.190:2379,192.168.50.138:2379,192.168.50.159:2379","db":{"host":"localhost","port":"27017","username":"admin","password":"123456","database":"ecs"}}'
```

创建对象 configer，此时 configer 对象会监听上面配置的 etcd 集群。

```go
cli, err := etcd.New(clientv3.Config{
  Endpoints:   "127.0.0.1:2379",
  DialTimeout: 5 * time.Second,
})
configer := tconf.New(&tconf.Config{
  FileName:   "config.example.yaml",
  Client:    cli,
}).Fetch(context.Background(), "/conf_example/")
```

获取配置项 cluster

```go
val, err := configer.Get("cluster")
if err != nil {
  log.Fatal(`get "cluster" error: `, err)
}
fmt.Println(`new clusters: `, val)
```

如果需要将配置转换到另一个自定义类型：

```go
type DBConfig struct {
  Database string   `json:"database"`
  Host     string   `json:"host"`
  Port     uint16   `json:"port"`
  UserName string   `json:"username"`
  Password string   `json:"password"`
}

func main(){
  db := &DBConfig{}
  // decode the config of the key "db" into db struct
  configer.UnmarshalKey("db", db)
  fmt.Printf("database: %s, host: %v, port: %d, username: %s, password: %s", db.Database, db.Host, db.Port, db.UserName, db.Password)
}
```

如果需要使用配置中心，需要开启监听来同步配置项。
**注：监听会开启新的 goroutine。**

```go
stop := make(chan bool)
configer.Watch(context.Background(), stop)
<-stop
```

使用回调做其他事：

```go
configer.WithOnConfigChanged(func(rsp *tconf.Response) {
  log.Printf("configurations was onchanged. prefix: %s, value: %s ", rsp.Path, string(rsp.Value))
})
```

具体可参考 [示例](examples/main.go)
