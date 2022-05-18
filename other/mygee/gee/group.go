// 设计分组
package gee

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc //支持的中间件
	// parent *RouterGroup //目前没有用了
	engine *Engine
}

// 创建group， Engine ”继承”了他，可以直接使用他的方法
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix, // 实现分组累加
		engine: engine,
		// parent: rg,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 通过group添加路由，实际上也是往engine的router中添加，只不过增加了group的prefix
func (rg *RouterGroup) addGroup(method string, prefix string, handler HandlerFunc) {
	rg.engine.addRoute(method, rg.prefix+prefix, handler)
}

func (rg *RouterGroup) GET(prefix string, handler HandlerFunc) {
	rg.addGroup("GET", prefix, handler)
}

func (rg *RouterGroup) POST(prefix string, handler HandlerFunc) {
	rg.addGroup("POST", prefix, handler)
}

// 加载中间件
func (rg *RouterGroup) Use(middlewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middlewares...)
}

// 添加支持访问的html文件夹root, 通过指定路径，打开文件
func (rg *RouterGroup) creatStaticHtml(relativepath string, fs http.FileSystem) HandlerFunc {
	// 解析请求地址，反应到真实的物理地址，使用http函数处理文件

	//1. 获取完全的虚拟地址
	abstructPath := path.Join(rg.prefix, relativepath)
	fileServer := http.StripPrefix(abstructPath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filename")
		if _, err := fs.Open(file); err != nil { // 文件打开失败
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 根据指定文件, 创建html的对外路由接口（对用户开放）
func (rg *RouterGroup) Static(relativepath string, root string) {
	handler := rg.creatStaticHtml(relativepath, http.Dir(root))
	// 加上前缀，创建路由func
	urlPattern := relativepath + "/*filename"
	rg.GET(urlPattern, handler)
}
