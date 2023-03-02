package filters

import (
	"fmt"
	"webSocket/web"
)

func init() {
	web.RegisterFilter("my-custom", myFilterBuilder)
}

func myFilterBuilder(next web.Filter) web.Filter {
	return func(ctx *web.Context) {
		fmt.Println("假装这是我自定义的 filter")
		next(ctx)
	}
}
