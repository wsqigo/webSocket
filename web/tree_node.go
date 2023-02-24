package web

import (
	"strings"
)

const anyPath = "*"

const (
	// 根节点，只有跟用这个
	nodeTypeRoot = iota

	// *
	nodeTypeAny

	// 路径参数
	nodeTypeParam

	// 正则
	nodeTypeReg

	// 静态，即完全匹配
	nodeTypeStatic
)

// matchFunc 承担两个职责，一个是判断是否匹配，一个是在匹配之后
// 将必要的数据写入到 Context
// 所谓必要的数据，这里基本上是路径参数
// ==
// child.path = *
// child.path 是路径参数，路径参数写入到 Context 里面去
// child.path.match(reg)
type matchFunc func(ctx *Context, path string) bool

type node struct {
	children []*node

	// 如果这是叶子节点
	// 那么匹配上之后就可以调用之该方法
	handler   handlerFunc
	matchFunc matchFunc

	// 原始的 pattern。注意，它不是完整的pattern
	// 而是匹配到这个节点的pattern
	path string
	// 标记是怎样的节点
	nodeType int
}

// 静态节点
func newStaticNode(path string) *node {
	return &node{
		matchFunc: func(ctx *Context, p string) bool {
			return path == p && p != "*"
		},
		nodeType: nodeTypeStatic,
		path:     path,
	}
}

func newRootNode(method string) *node {
	return &node{
		matchFunc: func(ctx *Context, p string) bool {
			panic("never call me")
		},
		nodeType: nodeTypeRoot,
		path:     method,
	}
}

func newNode(path string) *node {
	if path == "*" {
		return newAnyNode()
	}
	if strings.HasPrefix(path, ":") {
		return newParamNode(path)
	}
	return newStaticNode(path)
}

// 通配符 * 节点
func newAnyNode() *node {
	return &node{
		// 因为我们不允许 * 后面还有节点，所以这里可以不用初始化
		//children: make([]*node, 0, 2),
		matchFunc: func(ctx *Context, p string) bool {
			return true
		},
		nodeType: nodeTypeAny,
		path:     anyPath,
	}
}

// 路径参数节点
func newParamNode(path string) *node {
	paraName := path[1:]
	return &node{
		matchFunc: func(ctx *Context, p string) bool {
			if ctx != nil {
				ctx.PathParams[paraName] = p
			}
			// 如果自身是一个参数路由,
			// 然后又来一个通配符，我们认为是不匹配的
			return p != anyPath
		},
		nodeType: nodeTypeParam,
		path:     path,
	}
}

// 正则节点
//func newRegNode(path string) *node {
//	// 依据你的规则拿到正则表达式
//	return &node{
//		matchFunc: func(p string, ctx *Context) bool {
//			// 怎么写？
//		},
//		path:     path,
//		nodeType: nodeTypeReg,
//	}
//}
