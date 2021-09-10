package main

import (
	"context"
	"sync"
)

//UseContext 使用Context做超时控制
type UseContext struct {
}

func (j UseContext) Load(urls []string) (r map[string]int64) {
	r = make(map[string]int64, len(urls))
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, MaxWaitTime)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	defer cancel()
	for _, url := range urls {
		wg.Add(1)
		go func(url string, wg *sync.WaitGroup) {
			result := RealHttpGet(ctx, url)
			mu.Lock()
			r[url] = result
			mu.Unlock()
			wg.Done()
		}(url, &wg)
	}
	wg.Wait()
	return r
}
