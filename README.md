# tconf #

基于 viper，仅提供 etcdv3 APIs 的支持，原理是创建新的 etcd v3 的 provider，注入到 viper。
（后续会考虑支持 zookeeper、consul，但仍然推荐使用 etcd）

## 功能说明 ##

本项目仅提供从 etcd 配置中心读取配置，如果需要写入配置请在配置中心操作。
使用时至少指定一台配置中心节点，在程序第一次运行时，自动从指定的节点获取配置项，并替换本地配置。
开启服务器监听后，配置中心的任意改动都会被监听到，并替换本地配置。

## 示例 ##

写入 etcd 配置，下面是配置 mongodb 的数据库连接参数。
支持多个配置中心集群（或节点）地址，以逗号 "," 分隔。

```console
etcdctl put /conf_example/ '{"cluster":"192.168.50.190:2379,192.168.50.138:2379,192.168.50.159:2379","db":{"host":["localhost"],"port":"27017","username":"admin","password":"123456","database":"ecs"}}'
```

创建对象 configer，此时 configer 对象监听的是上面配置的 etcd 集群。

```go
configer := tconf.New(&tconf.Config{
  Path:       "/conf_example/",
  FileName:   "config.example.yaml",
  Cluster:    "127.0.0.1:2379",
  Renew:      true,
  ClusterKey: "cluster",
})
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
  Host     []string `json:"host"`
  Port     uint16   `json:"port"`
  UserName string   `json:"username"`
  Password string   `json:"password"`
}

func main(){
  db := &DBConfig{}
  // 把配置项 db 的值写入到 DBConfig 类型 db 对象
  configer.UnmarshalKey("db", db)
  fmt.Printf("database: %s, host: %v, port: %d, username: %s, password: %s", db.Database, db.Host, db.Port, db.UserName, db.Password)
}
```

如果需要使用配置中心，需要开启监听来同步配置项。
**注：监听会开启新的 goroutine。**

```go
configer.Watch()
```

具体可参考 [示例](examples/main.go)

## 需解决的问题 ##

- 在使用本地配置文件时，如果配置中心可用，能再次连接配置中心并监听配置项更改。
- 支持 zookeeper 和 consul。
