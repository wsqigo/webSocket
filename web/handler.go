package web

type Handler interface {
	ServeHTTP(ctx *Context)
	Routable
}

type handlerFunc func(ctx *Context)
