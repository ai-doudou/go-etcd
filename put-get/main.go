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

    // 存值
    _, err = client.Put(context.TODO(), "/api/wages/key", "1001")
    if err != nil {
        fmt.Println("put failed, err:", err)
        return
    }
    // 取值
    resp, err := client.Get(context.TODO(), "/api/wages/key")
    if err != nil {
        fmt.Println("get failed err:", err)
        return
    }
    // 遍历结果,返回key的列表
    for _, item := range resp.Kvs {
        fmt.Printf("%s : %s \n", item.Key, item.Value)
    }

    // 创建kv客户端
    kv := clientv3.NewKV(client)
    // 设置key,并获取值
    putResp, err := kv.Put(context.TODO(), "/cron/jobs/job1", "bye", clientv3.WithPrevKV())
    if err != nil {
        fmt.Println("kv put err:", err)
    }
    // 获取存取信息
    if putResp.PrevKv != nil {
        fmt.Println("key:", string(putResp.PrevKv.Key))
        fmt.Println("Value:", string(putResp.PrevKv.Value))
        fmt.Println("Version:", string(putResp.PrevKv.Version))
    }
    // 获取版本信息
    fmt.Println("Revision:", resp.Header.Revision)
    fmt.Println("ClusterId:", resp.Header.ClusterId)
    fmt.Println("MemberId:", resp.Header.MemberId)
}
