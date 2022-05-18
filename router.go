package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node      
	handlers map[string]HandlerFunc 
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc)}
}

func parsePattern(Pattern string) []string {
	vs := strings.Split(Pattern, "/")
	var parts []string
	for _, s := range vs {
		if s != "" {
			parts = append(parts, s)
			if s[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, Pattern string, handler HandlerFunc) {
	parts := parsePattern(Pattern)
	key := method + "-" + Pattern

	_, ok := r.roots[method]
	if !ok { 
		r.roots[method] = &node{}
	}
	r.roots[method].insert(Pattern, parts, 0)
	r.handlers[key] = handler
	log.Printf("Route %4s - %s", method, Pattern) 
}

func (r *router) getRouter(method string, path string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	searchPattern := parsePattern(path)

	n := root.search(searchPattern, 0)
	if n != nil {
		parts := parsePattern(n.pattern)

		log.Println(75, parts)

		params := make(map[string]string) 
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchPattern[index]
			}
			if (part[0] == '*') && (len(part) > 1) { 
				log.Println(83, part)
				params[part[1:]] = strings.Join(searchPattern[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}


func (r *router) handle(c *Context) {
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
	c.Next() 
}
