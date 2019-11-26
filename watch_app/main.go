package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "github.com/coreos/etcd/mvcc/mvccpb"
    "time"
)

// watch 监听
/*
监听kv变化:常用作与集群中配置下发，状态同步 非常有价值
 */
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

    keyName := "/cron/jobs/wages/status"
    kv := clientv3.NewKV(client)

    // 模拟etcd中KV的变化,这里用来测试
    go func() {
        for {
            _, err = kv.Put(context.TODO(), keyName, "jobVal "+time.Now().Format("2006.01.02.15.04.05"))
            if err != nil {
                fmt.Println("kv put err:", err)
            }
            _, err = kv.Delete(context.TODO(), keyName)
            if err != nil {
                fmt.Println("kv Delete err:", err)
            }
            time.Sleep(1 * time.Second)
        }
    }()


    // 先GET到当前的值，并监听后续变化
    getResponse, err := kv.Get(context.TODO(), keyName)
    if  err != nil {
        fmt.Println(err)
        return
    }
    // 现在key是存在的
    if len(getResponse.Kvs) != 0 {
        fmt.Println("当前值:", string(getResponse.Kvs[0].Value))
    }

    // 当前etcd集群事务ID, 单调递增的
    watchStartRevision := getResponse.Header.Revision + 1

    // 创建一个watcher
    watcher := clientv3.NewWatcher(client)
    // 启动监听
    fmt.Println("从该版本向后监听:", watchStartRevision)
    ctx, cancelFunc := context.WithCancel(context.TODO())
    time.AfterFunc(15 * time.Second, func() {
        cancelFunc()
    })

    watchRespChan := watcher.Watch(ctx, keyName, clientv3.WithRev(watchStartRevision))
    // 处理kv变化事件
    for watchResp := range watchRespChan {
        for _, event := range watchResp.Events {
            switch event.Type {
            case mvccpb.PUT:
                fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
            case mvccpb.DELETE:
                fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
            }
        }
    }


}
