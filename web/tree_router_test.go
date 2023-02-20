package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerBasedOnTree_Route(t *testing.T) {
	handler := NewHandlerBasedOnTree().(*HandlerBasedOnTree)
	assert.NotNil(t, handler.root)

	handler.Route(http.MethodPost, "/user", func(ctx *Context) {})

	// 开始做断言，这个时候我们应该确认，在根节点之下只有一个user节点
	assert.Equal(t, 1, len(handler.root.children))

	n := handler.root.children[0]
	assert.NotNil(t, n)
	assert.Equal(t, "user", n.path)
	assert.NotNil(t, n.handler)
	assert.Empty(t, n.children)

	// 我们只有
	//  user -> profile
	handler.Route(http.MethodPost, "/user/profile", func(c *Context) {})
	assert.Equal(t, 1, len(n.children))
	profileNode := n.children[0]
	assert.NotNil(t, profileNode)
	assert.Equal(t, "profile", profileNode.path)
	assert.NotNil(t, profileNode.handler)
	assert.Empty(t, profileNode.children)

	// 试试重复
	handler.Route(http.MethodPost, "/user", func(ctx *Context) {})
	n = handler.root.children[0]
	assert.NotNil(t, n)
	assert.Equal(t, "user", n.path)
	assert.NotNil(t, n.handler)
	// 有profile节点
	assert.Equal(t, 1, len(n.children))

	// 给 user 再加一个节点
	handler.Route(http.MethodPost, "/user/home", func(ctx *Context) {})
	assert.Equal(t, 2, len(n.children))
	homeNode := n.children[1]
	assert.NotNil(t, homeNode)
	assert.Equal(t, "home", homeNode.path)
	assert.NotNil(t, homeNode.handler)
	assert.Empty(t, homeNode.children)
}
