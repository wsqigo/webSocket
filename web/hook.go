package web

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Hook 是一个钩子函数。注意，
// ctx 是有一个超时机制的 context.Context
// 所以你必须处理超时的问题
type Hook func(ctx context.Context) error

// BuildCloseServerHook 这里其实可以考虑使用 errgroup
// 但是我们这里不用是希望每个server单独关闭
// 互相之间不影响
func BuildCloseServerHook(servers ...Server) Hook {
	return func(ctx context.Context) error {
		wg := sync.WaitGroup{}
		doneCh := make(chan struct{})
		wg.Add(len(servers))

		for _, s := range servers {
			go func(svr Server) {
				err := svr.Shutdown(ctx)
				if err != nil {
					fmt.Printf("server shutdown error: %v \n", err)
				}
				time.Sleep(time.Second)
				wg.Done()
			}(s)
		}

		// server有可能关不掉，为了支持超时机制
		// 开一个goroutine.
		go func() {
			// wait只会堵住当前goroutine
			wg.Wait()
			doneCh <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			fmt.Printf("closing servers timeout \n")
			return ErrorHookTimeout
		case <-doneCh:
			fmt.Printf("close all servers \n")
			return nil
		}
	}
}
