package main

import (
	"fmt"
	"net/url"
)

type WnlPrinter struct {
	wnl          *WebNodeList
	outputFormat OutputFormat
	treeLines    []bool
	lastDepth    int
	printIndex   int
}

func CreateWnlPrinter(wnl *WebNodeList, of OutputFormat) (wp WnlPrinter) {
	wp.wnl = wnl
	wp.outputFormat = of
	return
}

func (wp *WnlPrinter) treePrintNode(index int) {
	if index < 0 || index >= len(wp.wnl.list) {
		return
		// TODO: Proper error return
	}
	node := wp.wnl.list[index]
	if node.nodeDepth >= len(wp.treeLines) {
		wp.treeLines = append(wp.treeLines, false)
	}
	for i := 0; i < node.nodeDepth-1; i++ {
		if i < len(wp.treeLines) && wp.treeLines[i] {
			fmt.Print("│   ")
		} else {
			fmt.Print("    ")
		}
	}
	if node.nodeDepth > 0 {
		if node.nodeLastSibling {
			fmt.Print("└── ")
			wp.treeLines[node.nodeDepth-1] = false
		} else {
			fmt.Print("├── ")
			wp.treeLines[node.nodeDepth-1] = true
		}
	}
	fmt.Println(node.name)
}

func (wp *WnlPrinter) PrintDone() {
	l := wp.wnl
	l.mux.Lock()
	defer l.mux.Unlock()
	for ; wp.printIndex < len(l.list); wp.printIndex++ {
		if l.list[wp.printIndex].nodeStatus == done {
			switch wp.outputFormat {
			case tree:
				wp.treePrintNode(wp.printIndex)
			case urlencoded:
				fmt.Println(wp.wnl.list[wp.printIndex].path)
			case list:
				path, _ := url.PathUnescape(wp.wnl.list[wp.printIndex].path)
				fmt.Println(path)
			}
		} else {
			return
		}
	}
}
