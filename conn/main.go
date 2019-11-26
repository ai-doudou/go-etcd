package main

import (
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

func main() {
    // 客户端配置
    config := clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: 5 * time.Second,
    }
    fmt.Println("[etcd] clientv3 config:", config.Endpoints, config.DialTimeout)

    // 建立连接
    client, err := clientv3.New(config)
    if err != nil {
        fmt.Println("[etcd] clientv3 new", err)
        return
    }
    fmt.Println("connect success")

    // 关闭连接
    defer client.Close()
}
