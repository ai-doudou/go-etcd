package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

// 存储 读取value
func main() {

    // 客户端配置
    config := clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: 5 * time.Second,
    }
    fmt.Println("[etcd] clientv3 config:", config.Endpoints, config.DialTimeout)

    // 建立连接
    client, err := clientv3.New(config)
    // 关闭连接
    defer client.Close()
    if err != nil {
        fmt.Println("[etcd] clientv3 new", err)
        return
    }
    fmt.Println("[etcd] connect success")

    kv := clientv3.NewKV(client)
    resp, err := kv.Get(context.TODO(), "/api/wages/key", clientv3.WithPrefix())
    if err != nil {
        fmt.Println("kv get err:", err)
        return
    }

    fmt.Println(resp.Kvs, resp.Count)
}
