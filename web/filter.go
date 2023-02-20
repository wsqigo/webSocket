package web

import (
	"fmt"
	"time"
)

type FilterBuilder func(next Filter) Filter

type Filter func(ctx *Context)

func MetricFilterBuilder(next Filter) Filter {
	return func(ctx *Context) {
		// 执行前的时间
		startTime := time.Now()
		next(ctx)
		// 执行后的时间
		fmt.Println("run time:", time.Since(startTime))
	}
}
