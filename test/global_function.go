// 对基本功能进行简单的测试，部分地方跟gin稍有不同，但大体类似
// 测试代码写的比较简单，但对于严格的测试来说，最好每个功能都分模块测试

// logger 和 recovery 中间件会打印终端
package main

import (
	"fmt"
	"html/template"
	"time"

	"github.com/2103561941/gee"
)

// 自定义中间件，用于测试Group Use的接口是否正确
// 简单往客户端发送一条信息确认是否添加成功
func SendMessage() gee.HandlerFunc{
	return func(c *gee.Context) {
		c.String(200, "SendMessage Successfully!\n")
	}
}

// 自定义FunMap，对HTML模板title.tmpl的now进行格式化输出
func FormatAsDate(t time.Time) string{
	return fmt.Sprintf("%d,%d,%d\n", t.Year(), t.Month(), t.Day())
}


func main() {
	// 测试USE功能
	r := gee.Defalut()

	// 测试GET/String
	r.GET("/", func(c *gee.Context) {
		c.String(200, "Hello World!\n")
	})

	// 测试POST/JSON
	r.POST("/Post", func(c *gee.Context) {
		c.JSON(200, gee.H{
			"name": c.PostForm("name"),
			"age":  c.PostForm("age"),
		})
	})

	// 测试分组功能
	group := r.Group("/dynamic")

	// 测试分组添加中间件
	group.Use(SendMessage())

	// 测试动态路由 Param
	group.GET("/:name", func(c *gee.Context) {
		name := c.Param("name")
		c.String(200, "Hello %s\n", name)
	})

	// 测试动态路由*，查询文件
	group.GET("//floder/*filename", func(c *gee.Context) {
		filename := c.Param("filename")
		c.String(200, "GET file %s\n", filename)
	})

	// 测试HTML渲染功能
	// 自定义funcMap
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate" : FormatAsDate,
	})


	// 加载文件夹下所有的tmpl文件 (注意上下两步不能颠倒，tmpl要求的函数必须先存在在funcmap中才可以引入文件)
	r.LoadHTMLGlob("./css/*")


	// 使用模板文件进行渲染输出
	r.GET("/html", func(c *gee.Context) {
		c.HTML(200, "title.tmpl", gee.H{
			"title" : "This is Title",
			"now" : time.Date(2022, 5, 19, 18, 8, 0, 0, time.UTC),
		})
	})

	r.GET("/panic", func(c *gee.Context) {
		panic("man made panic")
	})
	

	// 测试Run功能
	r.Run(":9999")
}


/*
我使用的是postman进行测试

Output:
(after run....)
2022/05/19 14:51:30 Route  GET - /
2022/05/19 14:51:30 Route POST - /Post
2022/05/19 14:51:30 Route  GET - /dynamic/dynamic/:name
2022/05/19 14:51:30 Route  GET - /dynamic/dynamic/*filename

下面统一省略localhost:9999 的前缀
----------------------------
/
Hello World!

/Post 
(post name : cyb, age : 20)
{
    "age": "20",
    "name": "cyb"
}
----------------------------
/dynamic/:name
(URI:dynamic/cyb)
SendMessage Successfully!
Hello cyb
----------------------------
/dynamic/floder/*filename
(URI:dynamic/floder/css/1.jpg)
SendMessage Successfully!
GET file css/1.jpg
----------------------------
不加载FormatAsDate时html渲染的输出
hello, This is Title

Date: 2022-05-19 18:08:00 +0000 UTC
----------------------------
加载时的输出
hello, This is Title

Date: 2022,5,19
----------------------------
关于logger的测试就显示在终端里
2022/05/19 18:16:35 "2022-05-19 18:16:35.622826484 +0800 CST m=+47.857380262", URI = /html, StatusCode = 0

关于recovery的测试主要是通过在路由为/painc函数里增加panic函数进行测试
TraceBack:
        /usr/local/go/src/runtime/panic.go838
        /home/cyb/go_learn_code/goTinyProject/gee/test/global_function.go84
        /home/cyb/go_learn_code/goTinyProject/gee/context.go90
        /home/cyb/go_learn_code/goTinyProject/gee/recovery.go38
        /home/cyb/go_learn_code/goTinyProject/gee/context.go90
        /home/cyb/go_learn_code/goTinyProject/gee/logger.go14
        /home/cyb/go_learn_code/goTinyProject/gee/context.go90
        /home/cyb/go_learn_code/goTinyProject/gee/router.go93
        /home/cyb/go_learn_code/goTinyProject/gee/gee.go58
        /usr/local/go/src/net/http/server.go2917
        /usr/local/go/src/net/http/server.go1967
        /usr/local/go/src/runtime/asm_amd64.s1572

web:
{"message":"Internal Server Error\n"}

*/