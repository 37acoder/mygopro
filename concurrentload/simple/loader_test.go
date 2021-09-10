package simple

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/37acoder/mygopro/concurrentload"
)

func TestLoad(t *testing.T) {
	ctx := context.Background()
	c := make(chan int64, 0)
	go func() {
		for d := range c {
			fmt.Println(strconv.FormatInt(d, 10))
		}
	}()

	m := LoadManager{
		Loaders: []concurrentload.Loader{
			NewHttpWgetLoader("https://baidu.com", c),
			NewHttpWgetLoader("https://baidu.com", c),
			NewHttpWgetLoader("https://baidu.com", c),
			NewHttpWgetLoader("https://baidu.com", c),
			NewHttpWgetLoader("https://baidu.com", c),
		},
	}
	e := m.LoadAll(ctx)
	if e != nil {
		panic(e)
	}
	time.Sleep(2)
}
