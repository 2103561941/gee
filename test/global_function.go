// 对基本功能进行简单的测试，部分地方跟gin稍有不同，但大体类似
// 测试代码写的比较简单，但对于严格的测试来说，最好每个功能都分模块测试

package main

import "github.com/2103561941/gee"

// 自定义中间件，用于测试Group Use的接口是否正确


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


	// 测试动态路由 Param
	group.GET("/dynamic/:name", func(c *gee.Context) {
		name := c.Param("name")
		c.String(200, "Hello %s\n", name)
	})

	// 测试动态路由*，查询文件
	group.GET("/dynamic/*filename", func(c *gee.Context) {
		filename := c.Param("filename")
		c.String(200, "GET file %s\n", filename)
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

/
Hello World!

/Post 
(post name : cyb, age : 20)
{
    "age": "20",
    "name": "cyb"
}

/dynamic/:name
URI(dynamic)


*/