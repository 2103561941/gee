// 对新增前缀树进行测试
package gee

import (
	"testing"
)

// init router ，创建几个路由
func initTestRouter()  *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/nihao", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/assert/*filename", nil)
	return r
}

// func TestParse(t *testing.T) {
// 	ports := parsePattern("/hello/nihao///:name/*filename//*xx")
// 	// s := []string{"1", "2", "3"}
// 	t.Log(ports)
// 	// t.Log(s)
// 	// t.Logf("%T", &node{})
// }

// func TestMatchChild(t *testing.T) {
// 	n := &node{}
// 	n.matchChild("/hello")
// }

// func TestInsert(t *testing.T) {
// 	n := &node{
// 		children: make([]*node, 0),
// 	}
// 	pattern := "/hello/nihao/:name"
// 	parts := parsePattern(pattern)
// 	n.insert(pattern, parts, 0)
// 	t.Log(n)
// }

// func TestGetRouter(t *testing.T) {
// 	r := initTestRouter()
// 	t.Log(r)

// 	// n, _ := r.getRouter("GET", "/")
// 	// if n == nil {
// 	// 	t.Log(n, "n == nil")
// 	// } else {
// 	// 	t.Log(n.part)
// 	// }
// 	n, _ := r.getRouter("GET", "/hello/nihao")
// 	if n == nil {
// 		t.Fatal("n == nil, error")
// 	}

// 	if n.pattern != "/hello/nihao" {
// 		t.Fatal("n.pattern != /hello/nihao, error")
// 	}

// 	// if parmas["name"] != "cyb" {
// 	// 	t.Fatal("parmas[name] != cyb, error")
// 	// }

// 	t.Log("getRouter success")
// }


// func TestGetRouter1(t *testing.T) {
// 	r := initTestRouter()
// 	t.Log(r)

// 	n, parmas := r.getRouter("GET", "/hello/cyb")
// 	if n == nil {
// 		t.Fatal("n == nil, error")
// 	}

// 	if n.pattern != "/hello/:name" {
// 		t.Fatal("n.pattern != /hello/:name, error")
// 	}

// 	if parmas["name"] != "cyb" {
// 		t.Fatal("parmas[name] != cyb, error")
// 	}

// 	t.Log("getRouter success")
// }


func TestGetRouter2(t *testing.T) {
	r := initTestRouter()
	t.Log(r)

	n, params := r.getRouter("GET", "/assert/image.bmp")
	if n == nil {
		t.Fatal("n == nil, error")
	}

	if n.pattern != "/assert/*filename" {
		t.Fatal("n.pattern != /assert/*filename, error")
	}
	t.Log(params["filename"])

	t.Log("getRouter success")
}