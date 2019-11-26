package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

// 事务txn实现分布式锁
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

    // if 不存在key then 设置它 else 抢锁失败
    kv := clientv3.NewKV(client)

    // 创建事务
    txn := kv.Txn(context.TODO())

    // 定义事务
    // 如果key不存在
    txn.If(clientv3.Compare(clientv3.CreateRevision("mutex1"), "=", 0)).
        Then(clientv3.OpPut("mutex", "yes")).
        Else(clientv3.OpGet("mutex")) // 否则抢锁失败

    // 提交事务
    txnResp, err := txn.Commit()
    if err != nil {
        fmt.Println(err)
        return
    }

    // 判断是否抢到了锁
    if !txnResp.Succeeded {
        fmt.Println("锁被占用", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
        return
    }

    // 2 处理业务
    fmt.Println("处理任务")
    time.Sleep(3 * time.Second)
    // 在锁内 很安全

    // 3 释放锁
    // defer 会把租约释放掉，关联的kv就被删除了
    op:=clientv3.OpDelete("mutex1")
    kv.Do(context.TODO(),op)

    if resp, err := client.Get(context.TODO(), "mutex1"); err != nil {
        fmt.Println("client Get err:", err)
    } else {

        fmt.Println("client get len",len(resp.Kvs))
    }


}
