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

    // 1, 上锁 (创建租约, 自动续租, 拿着租约去抢占一个key)
    lease := clientv3.NewLease(client)

    // 申请一个5秒的租约
    leaseGrantResp, err := lease.Grant(context.TODO(), 15)
    if  err != nil {
        fmt.Println(err)
        return
    }

    // 拿到租约的ID
    leaseId := leaseGrantResp.ID

    // 准备一个用于取消自动续租的context
    ctx, cancelFunc := context.WithCancel(context.TODO())
    // 确保函数退出后, 自动续租会停止
    defer cancelFunc()

    // 5秒后会取消自动续租
    keepRespChan, err := lease.KeepAlive(ctx, leaseId)
    if  err != nil {
        fmt.Println(err)
        return
    }

    // 处理续约应答的协程
    go func() {

        var keepResp *clientv3.LeaseKeepAliveResponse
        for {
            select {
            case keepResp = <- keepRespChan:
                if keepRespChan == nil {
                    fmt.Println("租约已经失效了")
                    goto END
                } else {    // 每秒会续租一次, 所以就会受到一次应答
                    fmt.Println("收到自动续租应答:", keepResp.ID)
                }
            }
        }
    END:
    }()


    //  if 不存在key， then 设置它, else 抢锁失败
    kv := clientv3.NewKV(client)

    // 创建事务
    txn := kv.Txn(context.TODO())

    // 如果key不存在
    txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
        Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseId))).
        Else(clientv3.OpGet("/cron/lock/job9")) // 否则抢锁失败

    // 提交事务
    txnResp, err := txn.Commit()
    if  err != nil {
        fmt.Println(err)
        return
    }

    // 判断是否抢到了锁
    if !txnResp.Succeeded {
        fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
        return
    }

    // 2, 处理业务

    fmt.Println("处理任务")


    // 3, 释放锁(取消自动续租, 释放租约)
    // defer 会把租约释放掉, 关联的KV就被删除了
    defer lease.Revoke(context.TODO(), leaseId)
}
