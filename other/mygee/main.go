package main

import (
	"fmt"
	"gee"
	"net/http"
	"time"
)

// 简单设计一个中间件, 分别在handler两端输出时间, 被g1组别使用,
func MiddUse4g1() gee.HandlerFunc {
	return func(c *gee.Context) {
		c.String(200, "time is %q\n", time.Now())
		// log.Printf("MiddUse4g1 time :%q\n", time.Now())
		c.Next()
		c.String(200, "time is %q\n", time.Now())
		// log.Printf("MiddUse4g1 time :%q\n", time.Now())
	}
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	// r.SetMap(template.FuncMap{
	// 	"FormatAsDate": FormatAsDate,
	// })
	r.LoadHtmlGlob("templates/*")
	r.Static("/assets", "./static")
	// r.GET("/", func(c *gee.Context) {
	// 	// c.String(http.StatusOK, "URL.Path = %q\n", c.Path)
	// 	c.HTML(200, "", "<h1>Hello, world/h1>")
	// })

	// r.GET("/hello", func(c *gee.Context) {
	// 	c.String(200, "Hello, %q, you are at %q\n", c.Query("name"), c.Path)
	// })

	// r.GET("/hello/:name", func(c *gee.Context) {
	// 	c.String(200, "Hello, %s\n", c.Params["name"])
	// })

	// r.GET("/assert/*filename", func(c *gee.Context) {
	// 	c.String(200, "Get file %s\n", c.Params["filename"])
	// })

	// r.POST("/login", func(c *gee.Context) {
	// 	c.JSON(200, gee.H{
	// 		"name": c.PostForm("name"),
	// 		"age":  c.PostForm("age"),
	// 	})
	// })

	// v := r.Group("/g1")
	// v.Use(MiddUse4g1())

	// v.GET("/f1", func(c *gee.Context) {
	// 	c.String(200, "URL.Path = %q\n", c.Path)
	// })
	// v2 := v.Group("/gg1")
	// // v2.Use(MiddUse4g1()) // 父子组同时use了一个方法的情况, 会调用两次中间件，算是一个坑吧， 一般人也不会这么做
	// v2.GET("/ff1", func(c *gee.Context) {
	// 	c.String(200, "URL.Path = %q\n", c.Path)
	// })
	r.LoadHtmlGlob("./template/*.tmpl")
	r.SetMap(nil)

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.Run(":9999")
}
