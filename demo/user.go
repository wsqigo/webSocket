package demo

import (
	"fmt"
	"time"
	"webSocket/web"
)

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int    `json:"biz_code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
}

func SignUp(ctx *web.Context) {
	req := &signUpReq{}
	err := ctx.ReadJson(req)
	if err != nil {
		_ = ctx.BadRequestJson(&commonResponse{
			BizCode: 4, // 假如说我们这个代表输入参数错误
			// 注意这是demo，实际中你应该避免暴露 error
			Msg: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	_ = ctx.OkJson(&commonResponse{
		// 假设这个是新用户的ID
		Data: 123,
	})
}

func SlowService(ctx *web.Context) {
	time.Sleep(10 * time.Second)
	_ = ctx.OkJson(&commonResponse{
		Msg: "Hi, this is msg from slow service",
	})
}
