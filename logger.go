// 日志打印中间件
package gee

import (
	"log"
	"time"

)

func Logger() HandlerFunc {
	return func(c *Context) {
		log.Printf("%q, URI = %s, StatusCode = %d\n", time.Now(), c.Req.RequestURI, c.StatusCode)
	}
}
