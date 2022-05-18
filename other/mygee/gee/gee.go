// 简单实现gee框架
package gee

import (
	"html/template"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup // 用来集成group	
	router *router
	groups []*RouterGroup // 保存已有的分组
	htmlTemplate *template.Template
	funcMap template.FuncMap 
}

// 创建一个engine
func New() *Engine {
	engine := &Engine{
		router : newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine : engine}
	engine.groups = []*RouterGroup{engine.RouterGroup} // ??
	return engine
}

// 将创建的方法存放到“默认路由”里，在http listen时可以快速查询使用
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}


// 使用group结构的add，就可以不再使用单独的method了， 这段可以不需要
// // 存放设置的GET方法， 注意全大写
// func (engine *Engine) GET(pattern string, handler HandlerFunc) {
// 	engine.addRoute("GET", pattern, handler)
// }

// // 存放设置的POST方法
// func (engine *Engine) POST(pattern string, handler HandlerFunc) {
// 	engine.addRoute("POST", pattern, handler)
// }

// 使用listen
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 实现 SeverHTTP 的接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var middlewares []HandlerFunc

	// 找到跟reo.path同组的路由，保存其中间件
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			// 如果用户在父子关系的组别都使用了Use， 由于使用了slice和append，不会去重，导致中间件会被执行多次, 如果考虑这种情况我认为可以用map和struct{}来实现简单的set
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}


// 添加用户自定义的html渲染方法
func (engine *Engine) SetMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHtmlGlob(pattern string) {
	engine.htmlTemplate = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

//TODO
// func (engine *Engine) LoadHtmlFiles(files... string)