package main

import (
	"context"
	"fmt"
	"time"
)

const MaxGoroutineNum = 100

//ConcurrentLimit 使用生产者消费者模型控制并发数量
type ConcurrentLimit struct {
}


func (j ConcurrentLimit) Load(urls []string) map[string]int64 {
	r := make(map[string]int64, len(urls))
	closeChan := make(chan struct{})
	go func() {
		// 超时时间到后，该goroutine会关闭closeChan
		<-time.After(MaxWaitTime)
		close(closeChan)
	}()

	resultReceiver := make(chan Result, 0)
	//taskProducer := make(chan string, len(urls)) // 同步做
	taskPublisher := make(chan string, 0) // 异步做
	go func() { // 也可以同步做，将channel的缓存设置为len(urls)即可
		for _, url := range urls {
			taskPublisher <- url
		}
		close(taskPublisher)
	}()
	goNum := MaxGoroutineNum
	if MaxGoroutineNum > len(urls) {
		goNum = len(urls)
	}
	for i := 0; i < goNum; i++ {
		// 起goroutine异步消费任务，产生result
		go j.taskExecutor(closeChan, taskPublisher, resultReceiver)
	}
	// 消费result
	count := 0
	for {
		select {
		case rs := <-resultReceiver:
			r[rs.Url] = rs.Status
			fmt.Println("received result: ", rs)
			count += 1
			if count == len(urls) {
				return r
			}
		case <-closeChan:
			return r
		}
	}
}

func (j ConcurrentLimit) taskExecutor(closeChan <-chan struct{}, taskPublisher <-chan string, resultReceiver chan<- Result) {
	for {
		var rs int64
		select {
		// 避免任务还未执行时就已经超时了
		case <-closeChan:
			return
		// 接收任务
		case task, ok := <-taskPublisher:
			// 已无任务，通道已关闭
			if !ok {
				//fmt.Println("no task, goroutine exit")
				return
			}
			// 执行任务
			rs = RealHttpGet(context.Background(), task)
			select {
			// 避免执行完毕后超时或者已无Result的接收者
			case <-closeChan:
				//fmt.Println("timeout after task executed, url:", task)
				return
			// 生产Result
			case resultReceiver <- Result{Url: task, Status: rs}:
				continue
			}
		}
	}
}
