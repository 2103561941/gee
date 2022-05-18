package gee

import (
	"testing"
)

func TestGroup(t *testing.T) {
	r := New()
	r.GET("/index", func(c *Context) {
		c.String(200, "success")
	})

	v := r.Group("/hello")
	v.GET("/123", func(c *Context) {
		c.String(200, "success")
	})

}
