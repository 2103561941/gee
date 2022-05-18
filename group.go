// 采用分组的形式，

package gee

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc 
	engine *Engine
}

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix, 
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (rg *RouterGroup) addGroup(method string, prefix string, handler HandlerFunc) {
	rg.engine.addRoute(method, rg.prefix+prefix, handler)
}

func (rg *RouterGroup) GET(prefix string, handler HandlerFunc) {
	rg.addGroup("GET", prefix, handler)
}

func (rg *RouterGroup) POST(prefix string, handler HandlerFunc) {
	rg.addGroup("POST", prefix, handler)
}

func (rg *RouterGroup) Use(middlewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middlewares...)
}

func (rg *RouterGroup) creatStaticHTML(relativepath string, fs http.FileSystem) HandlerFunc {
	abstructPath := path.Join(rg.prefix, relativepath)
	fileServer := http.StripPrefix(abstructPath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filename")
		if _, err := fs.Open(file); err != nil { 
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (rg *RouterGroup) Static(relativepath string, root string) {
	handler := rg.creatStaticHTML(relativepath, http.Dir(root))
	urlPattern := relativepath + "/*filename"
	rg.GET(urlPattern, handler)
}
