package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request

	PathParams map[string]string
}

func (ctx *Context) ReadJson(data any) error {
	body, err := io.ReadAll(ctx.R.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, data)
}

func (ctx *Context) OkJson(data any) error {
	// http库里面提前定义好了各种响应码
	return ctx.WriteJson(http.StatusOK, data)
}

func (ctx *Context) WriteJson(status int, data any) error {
	ctx.W.WriteHeader(status)
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = ctx.W.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) BadRequestJson(data any) error {
	// http库里面提前定义好了各种响应码
	return ctx.WriteJson(http.StatusBadRequest, data)
}

func (ctx *Context) Reset(w http.ResponseWriter, r *http.Request) {
	ctx.W = w
	ctx.R = r
	ctx.PathParams = make(map[string]string, 1)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
		// 一般路径参数都是一个，所以容量1就可以了
		PathParams: make(map[string]string, 1),
	}
}

func newContext() *Context {
	fmt.Println("creat new context")
	return &Context{}
}
