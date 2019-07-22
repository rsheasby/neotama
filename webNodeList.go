package main

import (
	"sync"
)

func CreateWebNodeList() (result WebNodeList) {
	result.list = make([]WebNode, 0)
	return
}

type WebNodeList struct {
	mux sync.Mutex
	list []WebNode
}

func compareWebNodes(lhs, rhs WebNode) (validOrder bool) {
	rhsParent := rhs.GetParentPath()
	switch {
	case lhs.path == rhsParent:
		return true
	case lhs.GetParentPath() != rhsParent:
		return false
	case lhs.nodeType == rhs.nodeType:
		return lhs.path < rhs.path
	case lhs.nodeType == file && rhs.nodeType != file:
		return true
	default:
		return false
	}
}

func (l *WebNodeList) insertAtIndex(n WebNode, i int) {
	l.list = append(l.list, WebNode{})
	copy(l.list[i+1:], l.list[i:])
	l.list[i] = n
}

func (l *WebNodeList) containsPath(path string) (result bool) {
	for _, v := range(l.list) {
		if v.path == path {
			return true
		}
	}
	return false
}

func (l * WebNodeList) InsertSorted(nodes []WebNode) {
	l.mux.Lock()
	defer l.mux.Unlock()
	NODES:
	for _, node := range(nodes) {
		if !l.containsPath(node.path) {
			for k, v := range(l.list) {
				if compareWebNodes(node, v) {
					l.insertAtIndex(node, k)
					continue NODES
				}
			}
			l.list = append(l.list, node)
		}
	}
}

func (l *WebNodeList) IsDone() (bool) {
	l.mux.Lock()
	defer l.mux.Unlock()
	for _, v := range(l.list) {
		if v.nodeStatus == pending || v.nodeStatus == busy {
			return false
		}
	}
	return true
}

func (l *WebNodeList) GetPending() (result WebNode, wait bool) {
	l.mux.Lock()
	defer l.mux.Unlock()
	for k, v := range(l.list) {
		if v.nodeStatus == pending {
			l.list[k].nodeStatus = busy
			result = v
			return
		}
	}
	wait = true
	return
}

func (l *WebNodeList) SetStatus(path string, status NodeStatus) {
	l.mux.Lock()
	defer l.mux.Unlock()
	for k, v := range(l.list) {
		if v.path == path {
			l.list[k].nodeStatus = status
			return
		}
	}
}
