package web

import (
	"net/http"
	"strings"
)

type HandlerBasedOnTree struct {
	root *node
}

func NewHandlerBasedOnTree() Handler {
	return &HandlerBasedOnTree{
		root: &node{},
	}
}

// ServerHTTP 就是从树里面找节点
// 找到了就执行
func (h *HandlerBasedOnTree) ServeHTTP(ctx *Context) {
	handler, ok := h.findRouter(ctx.R.URL.Path)
	if !ok {
		ctx.W.WriteHeader(http.StatusNotFound)
		ctx.W.Write([]byte("Not Found"))
		return
	}
	handler(ctx)
}

func (h *HandlerBasedOnTree) findRouter(pattern string) (handlerFunc, bool) {
	// 去掉头尾可能有的/，然后按照/切割成段
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")
	cur := h.root

	for _, path := range paths {
		// 从子节点里边找一个到了当前 path 的节点
		matchNode, ok := h.findMatchChild(cur, path)
		if !ok {
			return nil, false
		}
		cur = matchNode
	}
	// 到这里,应该是找完了
	return cur.handler, cur.handler != nil
}

// Route 就相当于往树里面插入节点
func (h *HandlerBasedOnTree) Route(method string, pattern string, handlerFund handlerFunc) {
	// 将pattern按照URL的分隔符切割
	// 例如，/user/friends 将变成 [user, friends]
	// 将前后的/去掉，统一格式
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")
	// 当前指向根节点
	cur := h.root
	for index, path := range paths {
		// 从子节点里边找一个匹配到了当前 path 的节点
		matchChild, ok := h.findMatchChild(cur, path)
		if ok {
			cur = matchChild
		} else {
			// 为当前节点根据
			h.createSubTree(cur, paths[index:], handlerFund)
			return
		}
	}
	// 离开了循环，说明我们加入的是短路径
	// 比如说我们先加入了 /order/detail
	// 再加入/order，那么会走到这里
	cur.handler = handlerFund
}

func (h *HandlerBasedOnTree) findMatchChild(root *node, path string) (*node, bool) {
	for _, child := range root.children {
		if child.path == path {
			return child, true
		}
	}
	return nil, false
}

func (h *HandlerBasedOnTree) createSubTree(root *node, paths []string, handlerFunc handlerFunc) {
	cur := root
	for _, path := range paths {
		n := newNode(path)
		cur.children = append(cur.children, n)
		cur = n
	}
	cur.handler = handlerFunc
}

type node struct {
	path     string
	children []*node

	// 如果这是叶子节点
	// 那么匹配上之后就可以调用之该方法
	handler handlerFunc
}

func newNode(path string) *node {
	return &node{
		path: path,
	}
}
