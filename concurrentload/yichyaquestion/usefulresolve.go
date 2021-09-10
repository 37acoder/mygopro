package main

import (
	"context"
	"fmt"
	"time"
)

//UseContextAndChannel 使用Channel做并发同步和超时控制
type UseContextAndChannel struct {
}

type Result struct {
	Url    string
	Status int64
}

func (j UseContextAndChannel) Load(urls []string) (r map[string]int64) {
	r = make(map[string]int64, len(urls))
	resultReceiver := make(chan Result, len(urls))
	closeChan := make(chan struct{})
	go func() {
		// 超时时间到后，该goroutine会关闭closeChan
		<-time.After(MaxWaitTime)
		close(closeChan)
	}()
	for _, url := range urls {
		// 并发加载
		go j.get(closeChan, url, resultReceiver)
	}
	count := 0
	for {
		select {
		// 收集响应
		case rs := <-resultReceiver:
			r[rs.Url] = rs.Status
			count += 1
			// 全部加载完毕
			if count == len(urls) {
				return
			}
		// 避免超时
		case <-closeChan:
			return
		}
	}
}

func (j UseContextAndChannel) get(closeChan <-chan struct{}, url string, resultReceiver chan<- Result) {
	// 不受控制的阻塞操作
	result := RealHttpGet(context.Background(), url)
	select {
	case resultReceiver <- Result{Url: url, Status: result}:
	case <-closeChan: // 避免resultReceiver的接收者停止接收而一直阻塞
		fmt.Println("timeout url:", url)
		return
	}
}
