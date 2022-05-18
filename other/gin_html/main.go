// 测试使用html模板格式

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// func main() {
// 	engine := gin.Default()
// 	engine.LoadHTMLGlob("./template/*.tmpl")
// 	engine.GET("/", func(c *gin.Context) {
// 		c.String(200, "success")
// 		c.HTML(200, "title.tmpl", gin.H{
// 			"title": "hello, my name is cyb",
// 		})
// 	})
// 	engine.Run(":9999")
// }

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl")
	// 如果使用 LoadHTMLFiles 的话这么做（需要列举所有需要加载的文件，不如上述 LoadHTMLGlob 模式匹配方便）：
	// router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/title", func(c *gin.Context) {
		c.HTML(http.StatusOK, "123.tmpl", gin.H{
			"title": "hello, my name is cyb",
		})
	})
	router.Run(":8080")
}
