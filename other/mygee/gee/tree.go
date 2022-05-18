// 使用前缀树匹配part，实现动态路由
package gee

import (
	"strings"
)

// 创建路由前缀树模型
type node struct {
	pattern  string  // path : /p/:lang , 只有在最后一个part，这个值才会是path，否则这个值就是空，用来判断是不是匹配到了最后一个节点
	part     string  // a part of path : p or :lang
	children []*node // child parts
	isWild   bool    // 是否是动态匹配，动态匹配值为true
}

// 节点插入匹配, 匹配到相同节点，直接结束, 将整个path的匹配按 / 分割成多个part的匹配
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if (child.part == part) || (child.isWild) { // part 相同 或者 动态匹配 就直接返回
			return child // 匹配到了就返回匹配到的节点，进行该节点的子节点的匹配
		}
	}
	return nil // 没匹配到就返回nil， 用于后续判断有无匹配结果
}

// 函数寻找匹配，用于与用户输入url的匹配, 可能有多个匹配正确的结果，都要返回
func (n *node) matchChildren(part string) []*node {
	var nodes []*node
	for _, child := range n.children {
		if (child.part == part) || (child.isWild) {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 寻找匹配的节点, 插入对应的
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height { // 判断是否已经匹配到最后一个part
		// 最后一个node 保存Pattern
		n.pattern = pattern
		return
	}

	part := parts[height] // 待匹配的节点的part
	child := n.matchChild(part)
	if child == nil { // 没有匹配上, 创建一个新的分支
		child = &node{
			part:   part,
			isWild: (part[0] == ':') || (part[0] == '*'), // 判断是否是动态路由
		}
	}
	n.children = append(n.children, child)

	// log.Println("*************************")

	child.insert(pattern, parts, height+1) // 通过已经匹配的子节点去递归搜搜
}

// 寻找匹配的节点，并返回，用于执行对应的handler
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { // 完全匹配成功, *在末尾表示匹配任意内容
		if n.pattern == "" { // 不完全匹配， 因为没有匹配到路由的全部路径，只匹配了部分
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	// 对于成功匹配的child， 进行 11 递归
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil { // 存在匹配成功的结果
			return result
		}
	}

	return nil
}
