// http使用handler 的 w req 对象繁琐， 我们使用context对象对其进行封装， 并添加常用功能
package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// param
	Params map[string]string
	// response info
	StatusCode int
	// handlefunc 存放中间件函数
	index int // 存放当前执行到第几个中间件
	handlers []HandlerFunc // 存放中间件
	engine *Engine // 使用用户设置的html渲染方法等等

}

// 封装一个context，不需要自己装配
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index : -1,
	}
}

// 获取post信息
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 返回设置状态信息
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}


// 发送数据回客户端, code : 状态码， 格式化string
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}


// 通过json格式发送数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500) // 服务器出错
	}
}

// 普通发送数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 发送html格式的数据
func (c *Context) HTML(code int , name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplate.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
	// c.Writer.Write([]byte(html))
}

// 获取param数据
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// next 函数，往下执行中间件，对于部分中间件来说他不会直接执行完，可能中间会传出handler和其他中间件的执行， 所以在调用中间件的时候可以提供next函数用于执行后面的语句
func (c *Context) Next() {
	c.index++
	lenth := len(c.handlers)
	for ; c.index < lenth; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(status int, err string) {
	c.index = len(c.handlers)
	c.JSON(status, H{"message": err})
}
