package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

// 存储 读取value
func main() {

    var (
        grantResponse *clientv3.LeaseGrantResponse
        putResponse *clientv3.PutResponse
        getResponse *clientv3.GetResponse
    )


    // 客户端配置
    config := clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: 5 * time.Second,
    }

    // 建立连接
    client, err := clientv3.New(config)
    // 关闭连接
    defer client.Close()
    if err != nil {
        fmt.Println("clientv3 new err", err)
        return
    }

    // 创建续租
    lease := clientv3.NewLease(client)

    // 申请一个10秒的租约
    if grantResponse, err = lease.Grant(context.TODO(), 15);err != nil {
        fmt.Println("lease grant err:", err)
        return
    }

    leaseId := grantResponse.ID
    kv := clientv3.NewKV(client)

    // put一个kv 让它与租约关联起来 从而实现10秒自动过期
    if putResponse, err = kv.Put(context.TODO(), "cron/wages/lock", "yes", clientv3.WithLease(leaseId));err != nil {
        fmt.Println("kv put err", err)
        return
    }

    fmt.Println("写入成功", putResponse.Header.Revision)

    // 定时的看一下key过期了没有
    for {

        if getResponse, err = kv.Get(context.TODO(), "cron/wages/lock");err != nil {
            fmt.Println(err)
            return
        }
        if getResponse.Count == 0 {
            fmt.Println("kv过期了")
            break
        }
        fmt.Println("还没过期：", getResponse.Kvs)
        time.Sleep(time.Second * 2)
    }
}
