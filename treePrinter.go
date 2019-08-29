package main

import (
	"fmt"
)

type TreePrinter struct {
	wnl *WebNodeList
	treeLines []bool
	lastDepth int
}

func CreateTreePrinter(wnl *WebNodeList) (tp TreePrinter) {
	tp.wnl = wnl
	return
}

func (tp *TreePrinter) printNode(index int) {
	if index < 0 || index >= len(tp.wnl.list) {
		return
		// TODO: Proper error return
	}
	node := tp.wnl.list[index]
	if node.nodeDepth >= len(tp.treeLines) {
		tp.treeLines = append(tp.treeLines, false)
	}
	for i := 0; i < node.nodeDepth - 1; i++ {
		if i < len(tp.treeLines) && tp.treeLines[i] {
			fmt.Print("│   ")
		} else {
			fmt.Print("    ")
		}
	}
	if node.nodeDepth > 0 {
		if node.nodeLastSibling {
			fmt.Print("└── ")
			tp.treeLines[node.nodeDepth - 1] = false
		} else {
			fmt.Print("├── ")
			tp.treeLines[node.nodeDepth - 1] = true
		}
	}
	fmt.Println(node.name)
}

func (tp * TreePrinter) PrintDone() {
	l := tp.wnl
	l.mux.Lock()
	defer l.mux.Unlock()
	for ;l.printPointer < len(l.list); l.printPointer++ {
		if l.list[l.printPointer].nodeStatus == done {
			tp.printNode(l.printPointer)
		} else {
			return
		}
	}
}
