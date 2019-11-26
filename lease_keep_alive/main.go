package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

// 存储 读取value
func main() {

    // 配置
    config := clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: time.Second * 5,
    }
    // 连接 床见一个客户端
    client, err := clientv3.New(config)
    if err != nil {
        fmt.Println(err)
        return
    }

    // 申请一个lease 租约
    lease := clientv3.NewLease(client)

    // 申请一个10秒的租约
    leaseGrantResp, err := lease.Grant(context.TODO(), 10)
    if err != nil {
        fmt.Println(err)
        return
    }

    // 拿到租约id
    leaseid := leaseGrantResp.ID

    // 获得kv api子集
    kv := clientv3.NewKV(client)

    // 自动续租
    keepRestChan, err := lease.KeepAlive(context.TODO(), leaseid)
    if err != nil {
        fmt.Println(err)
        return
    }

    // 处理续租应答的协程
    var keepresp *clientv3.LeaseKeepAliveResponse
    go func() {

        for {
            select {
            case keepresp = <-keepRestChan:
                if keepRestChan == nil {
                    fmt.Println("租约已失效了")
                    goto END
                } else { // 每秒会续租一次，所以就会收到一次应答
                    fmt.Println("收到自动续租的应答")
                }
            }
        }
    END:
    }()

    // put一个kv 让它与租约关联起来 从而实现10秒自动过期
    putResp, err := kv.Put(context.TODO(), "cron/lock/job1", "v5", clientv3.WithLease(leaseid))
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("写入成功", putResp.Header.Revision)

    // 定时的看一下key过期了没有
    for {
        getResp, err := kv.Get(context.TODO(), "cron/lock/job1")
        if err != nil {
            fmt.Println(err)
            return
        }
        if getResp.Count == 0 {
            fmt.Println("kv过期了")
            break
        }
        fmt.Println("还没过期：", getResp.Kvs)
        time.Sleep(time.Second * 2)
    }
}
