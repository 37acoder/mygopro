package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	MaxWaitTime = time.Millisecond * 1000
	UrlNum      = 200
)

//func MockHttpGet(url string) int64 {
//	fmt.Println("start load url:", url)
//	time.Sleep(time.Duration((rand.Float32()-0.5)*float32(200*time.Millisecond)) + MaxWaitTime)
//	fmt.Println("loaded url:", url)
//	return 200
//}

func RealHttpGet(ctx context.Context, url string) int64 {
	fmt.Println("start load url:", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	fmt.Println("loaded url", url)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil {
		fmt.Println(err)
		return 0
	}
	return int64(resp.StatusCode)
}

type Resolver interface {
	Load(urls []string) map[string]int64
}

func main() {
	urls := make([]string, UrlNum)
	for i := 0; i < UrlNum; i++ {
		//urls[i] = strconv.Itoa(i)
		urls[i] = "https://baidu.com/?arg=" + strconv.Itoa(i)
	}
	//var resolver Resolver = NoConcurrentResolver{} // change it
	//var resolver Resolver = JustUseGoKeyWord{} // change it
	//var resolver Resolver = UseMutexButWrongWaitGroup{} // change it
	//var resolver Resolver = UseMutexAndWaitGroup{} // change it
	//var resolver Resolver = UseContext{} // change it
	//var resolver Resolver = UseContextAndChannel{} // change it
	var resolver Resolver = ConcurrentLimit{} // change it
	startTime := time.Now()
	r := resolver.Load(urls)
	sucCount := 0
	for _, v := range r {
		//fmt.Printf("Load url:%s\tstatus:%d\n", k, v)
		if v != 0 {
			sucCount += 1
		}
	}
	fmt.Println("success num:", sucCount)
	fmt.Println("time consume:", time.Since(startTime))
}
