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

var builderMap = make(map[string]FilterBuilder, 4)

func RegisterFilter(name string, builder FilterBuilder) {
	// 情况1 有些时候你可能不允许重复注册，那么你要先检测是否已经注册过了
	// 情况2 你会在并发的环境下调用这个方法，那么你应该
	builderMap[name] = builder
}

func GetFilterBuilder(name string) FilterBuilder {
	// 如果你觉得名字必须是正确的，那么你同样需要检测
	return builderMap[name]
}
