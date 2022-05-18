// 日志打印中间件
package gee

import (
	"log"
	"time"

)

func Logger() HandlerFunc {
	return func(c *Context) {
		log.Printf("%q, URL.Path = %s\n", time.Now(), c.Path)
	}
}
