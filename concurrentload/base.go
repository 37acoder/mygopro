package concurrentload

import (
	"context"
	"errors"
)

var (
	ErrorFatalToStopLoad = errors.New("fatal error happened, stop load all")
)

type Loader interface {
	Prepare(ctx context.Context) error // 返回了错误则不进行之后的步骤
	Load(context.Context)              // 并行Load
}

type LoadManager interface {
	Register(loader Loader)
	LoadAll(ctx context.Context) error
}
