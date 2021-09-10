package main

import (
	"context"
	"sync"
)

//NoConcurrentResolver 无并发解法
type NoConcurrentResolver struct {
}

func (n NoConcurrentResolver) Load(urls []string) (r map[string]int64) {
	r = map[string]int64{} // wrong
	//r = make(map[string]int64, len(urls)) // right
	for _, url := range urls {
		r[url] = RealHttpGet(context.Background(), url)
	}
	return r
}

// JustUseGoKeyWord 仅使用go关键字并发
type JustUseGoKeyWord struct {
}

func (j JustUseGoKeyWord) Load(urls []string) map[string]int64 {
	r := map[string]int64{}
	for _, url := range urls {
		go func() {
			r[url] = RealHttpGet(context.Background(), url)
		}()
	}
	return r
}

//UseMutexButWrongWaitGroup 错误使用使用互斥锁和WaitGroup保证并发安全和同步
type UseMutexButWrongWaitGroup struct {
}

func (j UseMutexButWrongWaitGroup) Load(urls []string) map[string]int64 {
	r := map[string]int64{}
	mu := &sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, url := range urls {
		go func(url string, wg sync.WaitGroup) {
			wg.Add(1)
			result := RealHttpGet(context.Background(), url)
			mu.Lock()
			r[url] = result
			mu.Unlock()
			wg.Done()
		}(url, wg)
	}
	wg.Wait()
	return r
}

//UseMutexAndWaitGroup 正确使用互斥锁和WaitGroup保证并发安全和同步
type UseMutexAndWaitGroup struct {
}

func (j UseMutexAndWaitGroup) Load(urls []string) map[string]int64 {
	r := make(map[string]int64, len(urls))
	mu := &sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string, wg *sync.WaitGroup) {
			result := RealHttpGet(context.Background(), url)
			mu.Lock()
			r[url] = result
			mu.Unlock()
			wg.Done()
		}(url, &wg)
	}
	wg.Wait()
	return r
}
