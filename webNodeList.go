package main

import (
	"sync"
	"unicode"
)

func CreateWebNodeList() (result WebNodeList) {
	result.list = make([]WebNode, 0)
	return
}

type WebNodeList struct {
	mux sync.Mutex
	list []WebNode
}

func compareChar(lhs, rhs byte) (validOrder bool) {
	if lhs == rhs {
		return true
	}
	switch {
		case unicode.ToLower(rune(lhs)) == unicode.ToLower(rune(rhs)): return lhs > rhs
		case rhs == '/': return true
		case lhs == '/': return false
		default: return unicode.ToLower(rune(lhs)) < unicode.ToLower(rune(rhs))
	}
}

func intMin(a, b int) (min int) {
	if a < b {
		return a
	} else {
		return b
	}
}

func lexCompare(lhs, rhs string) (validOrder bool) {
	length := intMin(len(lhs), len(rhs))
	for i := 0; i < length; i++ {
		if lhs[i] != rhs[i] {
			return compareChar(lhs[i], rhs[i])
		}
	}
	return len(lhs) < len(rhs)
}

func compareWebNodes(lhs, rhs WebNode) (validOrder bool) {
	return lexCompare(lhs.path, rhs.path)
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
