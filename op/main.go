package main

import (
    "context"
    "fmt"
    "github.com/coreos/etcd/clientv3"
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
    // 关闭连接
    defer client.Close()
    if err != nil {
        fmt.Println("[etcd] clientv3 new", err)
        return
    }
    fmt.Println("[etcd] connect success")

    kv := clientv3.NewKV(client)

    // 创建op
    op := clientv3.OpPut("/cron/jobs/job8", "111111111111111")

    // 执行op
    response, err := kv.Do(context.TODO(), op)
    if err != nil {
        fmt.Println("kv Do err:", err)
        return
    }
    // 获取version
    fmt.Println("写入Revision:", response.Put().Header.Revision)

    op=clientv3.OpGet("/cron/jobs/job8")
    // 执行OP
    opResp, err := kv.Do(context.TODO(), op)
    if  err != nil {
        fmt.Println("kv do err",err)
        return
    }

    // 打印
    fmt.Println("数据Revision:", opResp.Get().Kvs[0].ModRevision)    // create rev == mod rev
    fmt.Println("数据value:", string(opResp.Get().Kvs[0].Value))
}
