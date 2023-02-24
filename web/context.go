package web

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request

	PathParams map[string]string
}

func (c *Context) ReadJson(data any) error {
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, data)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
		// 一般路径参数都是一个，所以容量1就可以了
		PathParams: make(map[string]string, 1),
	}
}
