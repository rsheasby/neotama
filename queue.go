package main

import (
	"regexp"
	"sync"
	"fmt"
)

var nodeList []webNode
var listLock sync.Mutex

var dirParentRegex *regexp.Regexp
var fileParentRegex *regexp.Regexp

func init() {
	nodeList = make([]webNode, 0)
	dirParentRegex = regexp.MustCompile(`(.*\/).*\/`)
	fileParentRegex = regexp.MustCompile(`(.*\/).*`)
}

func (node webNode) getParentNode() (result string) {
	// fmt.Println("Getting parent node of: " + node.path)
	if node.nodeType == file {
		result = fileParentRegex.FindStringSubmatch(node.path)[1]
	} else {
		result = dirParentRegex.FindStringSubmatch(node.path)[1]
	}
	return
}

func compareNodes(lhs, rhs webNode) (validOrder bool) {
	switch {
	case lhs.path == rhs.getParentNode():
		validOrder = true
	case lhs.getParentNode() == lhs.getParentNode():
		switch {
		case lhs.nodeType == rhs.nodeType:
			validOrder = lhs.path < rhs.path
		case lhs.nodeType == file && rhs.nodeType != file:
			validOrder = true
		default:
			validOrder = false
		}
	default:
		validOrder = false
	}
	return
}

func queueInsert(node webNode) {
	fmt.Println("Inserting node into queue: " + node.path)
	listLock.Lock()
	defer listLock.Unlock()
	for k, v := range(nodeList) {
		if compareNodes(node, v) {
			insertWebNode(&nodeList, node, k)
			return
		}
	}
	nodeList = append(nodeList, node)
}

func queueDone() (result bool) {
	for _, v := range(nodeList) {
		if v.status != explored {
			return false
		}
	}
	return true
}

func queuePop() (node webNode, wait bool, done bool) {
	listLock.Lock()
	defer listLock.Unlock()
	done = queueDone()
	if done {
		done = true
		return
	}
	for k, v := range(nodeList) {
		if v.status == unexplored {
			nodeList[k].status = inProgress
			node = v
			return
		}
	}
	wait = true

	return
}

func markNodeDone(node webNode) {
	fmt.Println("Marking node as done: ", node.path)
	listLock.Lock()
	defer listLock.Unlock()
	for k, v := range(nodeList) {
		if node.path == v.path {
			nodeList[k].status = explored
		}
	}
}
