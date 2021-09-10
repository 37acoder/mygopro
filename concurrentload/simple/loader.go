package simple

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/37acoder/mygopro/concurrentload"
)

type LoadManager struct {
	Loaders []concurrentload.Loader
}

func (s *LoadManager) Register(loader concurrentload.Loader) {
	s.Loaders = append(s.Loaders, loader)
}

func (s *LoadManager) LoadAll(ctx context.Context) error {
	ctx, canceler := context.WithCancel(ctx)
	wg := sync.WaitGroup{}
	for _, loader := range s.Loaders {
		err := loader.Prepare(ctx)
		if err != nil {
			if errors.Is(err, concurrentload.ErrorFatalToStopLoad) {
				canceler()
				return err
			} else {
				continue
			}
		}
		wg.Add(1)
		go func(loader concurrentload.Loader) {
			defer func() {
				wg.Done()
				if p := recover(); p != nil {
					fmt.Printf("panic when loading, loader:%v, reason:%v\n", loader, p)
				}
			}()
			loader.Load(ctx)
		}(loader)
	}
	wg.Wait()
	canceler()
	return nil
}

type HttpWgetLoader struct {
	Url             string
	resultCollector chan int64
	Error           error
}

func NewHttpWgetLoader(url string, resultCollector chan int64) *HttpWgetLoader {
	return &HttpWgetLoader{
		Url:             url,
		resultCollector: resultCollector,
		Error:           nil,
	}
}

func (h *HttpWgetLoader) Prepare(ctx context.Context) error {
	return nil
}

func (h *HttpWgetLoader) Load(ctx context.Context) {
	data, err := http.Get(h.Url)
	if err != nil {
		h.Error = err
		return
	}
	h.resultCollector <- int64(data.StatusCode)
}
