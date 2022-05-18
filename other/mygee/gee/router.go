package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       //  存放method方法对应路由的父节点
	handlers map[string]HandlerFunc // 存放整个path对应的路由的handler
}

// 创建一个路由集
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc)}
}

// 将完整的pattern 拆分成parts数组，方便进行比对操作
func parsePattern(Pattern string) []string {
	vs := strings.Split(Pattern, "/")

	var parts []string

	for _, s := range vs {
		if s != "" {
			parts = append(parts, s)
			if s[0] == '*' { // 不考虑后面的内容
				break
			}
		}
	}

	return parts
}

// 装载一个路由
func (r *router) addRoute(method string, Pattern string, handler HandlerFunc) {
	parts := parsePattern(Pattern)
	key := method + "-" + Pattern

	// 查看是否roots里面是否已经有该method的root节点
	_, ok := r.roots[method]
	if !ok { // 没有就建一个
		r.roots[method] = &node{}
	}

	// 找到的话就执行插入操作，把当前方法对应的pattern插入到root中, 对前缀树进行一个扩容
	r.roots[method].insert(Pattern, parts, 0)

	r.handlers[key] = handler
	log.Printf("Route %4s - %s", method, Pattern) // 打印日志
}

// 获取路由信息,
func (r *router) getRouter(method string, path string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok { // 没有设置该方法
		return nil, nil
	}

	// log.Println("********1")

	searchPattern := parsePattern(path)
	// var params map[string]string

	n := root.search(searchPattern, 0)
	if n != nil { // 有相应的节点
		// 针对 ： 和 * 返回相应的param, 因为一个path上可能有很多个 ：，所以需要都放入param供用户使用
		parts := parsePattern(n.pattern)

		log.Println(75, parts)

		params := make(map[string]string) // 这里必须先给他分配make， 而不能声明直接在nil的map上进行操作
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchPattern[index]
			}
			if (part[0] == '*') && (len(part) > 1) { // * -> /../../..
				log.Println(83, part)
				params[part[1:]] = strings.Join(searchPattern[index:], "/")
				break
			}
		}
		// log.Println(n, "********2")
		return n, params
	}

	// log.Println(n, "********3")
	return nil, nil
}

// 对指定url进行处理
func (r *router) handle(c *Context) {
	// 由于用户输入的路由和我们存放的路由是不一样的（动态路由，虽然效果一样）， 但是我们要把roots里的pattern拿出来，作为handlers 的key
	// n -> *node (pattern),
	n, params := r.getRouter(c.Method, c.Path)

	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		handler := r.handlers[key]
		c.handlers = append(c.handlers, handler)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND, URL.Path = %s\n", c.Path)
		})

	}
	// c.String(200, "c has %d middlewares\n", len(c.handlers) - 1)
	c.Next() // 按序执行中间件和handler函数
}
