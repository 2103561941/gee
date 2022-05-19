package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 使用runtime包，打印报错路径上的信息
func trace(messgae string) string{
	var pcs [32]uintptr
	// 只打印前三层调用的文件名和所在行数
	n := runtime.Callers(3, pcs[:])
	// 使用strings.Builder 累加速度更快，内存分配更少
	var str strings.Builder
	str.WriteString(messgae + "\nTraceBack:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s  %d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc{
	return func(c *Context) {
		defer func() {
			// 恢复错误的情况，打印错误日志
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error\n")
			}	
		}()
		c.Next()
	}
}