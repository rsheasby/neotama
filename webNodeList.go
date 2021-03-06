package main

import (
	"sort"
	"sync"
	"unicode"
)

func CreateWebNodeList() (result WebNodeList) {
	result.list = make([]WebNode, 0)
	return
}

type WebNodeList struct {
	mux  sync.Mutex
	list []WebNode

	pendingCount int
	busyCount    int
	doneCount    int
	failCount    int

	busyPointer int
}

func compareStrings(lhs, rhs string) (validOrder bool) {
	// Create left and right rune slices
	lrs, rrs := []rune(lhs), []rune(rhs)
	for i := 0; i < len(lrs) && i < len(lrs); i++ {
		if lhs != rhs {
			lhsLower := unicode.ToLower(lrs[i])
			rhsLower := unicode.ToLower(rrs[i])
			if lhsLower == rhsLower {
				return lrs[i] < rrs[i]
			} else {
				return lhsLower < rhsLower
			}
		}
	}
	return len(lhs) < len(rhs)
}

func compareWebNodes(lhs, rhs *WebNode) (validOrder bool) {
	switch {
	case lhs.nodeType == rhs.nodeType:
		return compareStrings(lhs.name, rhs.name)
	case lhs.nodeType == file:
		return true
	default:
		return false
	}
}

func (l *WebNodeList) insertAtIndex(n []WebNode, index int) {
	insertLen := len(n)
	l.list = append(l.list, make([]WebNode, insertLen)...)
	copy(l.list[index+insertLen:], l.list[index:])
	for i := 0; i < insertLen; i++ {
		l.list[index+i] = n[i]
	}
}

func (l *WebNodeList) InsertSorted(nodes []WebNode, parentPath string, sortNodes bool) {
	if len(nodes) == 0 {
		return
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	if sortNodes {
		sort.Slice(nodes, func(i, j int) bool { return compareWebNodes(&nodes[i], &nodes[j]) })
	}
	for i := 0; i < len(nodes); i++ {
		if nodes[i].nodeStatus == pending {
			l.pendingCount++
		}
	}
	for l.busyPointer < len(l.list) && l.list[l.busyPointer].nodeStatus == done {
		l.busyPointer++
	}
	for i := l.busyPointer; i < len(l.list); i++ {
		if l.list[i].nodeStatus == busy && l.list[i].path == parentPath {
			for k := range nodes {
				nodes[k].nodeDepth = l.list[i].nodeDepth + 1
			}
			nodes[len(nodes)-1].nodeLastSibling = true
			l.insertAtIndex(nodes, i+1)
			return
		}
	}
	l.list = append(l.list, nodes...)
}

func (l *WebNodeList) IsDone() bool {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.pendingCount == 0 && l.busyCount == 0
}

func (l *WebNodeList) GetPending() (result WebNode, wait bool) {
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.pendingCount == 0 {
		wait = true
		return
	}
	for i := range l.list {
		if l.list[i].nodeStatus == pending {
			l.setStatusByIndex(i, busy)
			result = l.list[i]
			return
		}
	}
	wait = true
	return
}

func (l *WebNodeList) SetFailed(path string) {
	l.mux.Lock()
	defer l.mux.Unlock()
	for i := range l.list {
		if l.list[i].path == path {
			if !l.list[i].nodeFail {
				l.failCount++
				l.list[i].nodeFail = true
			}
			return
		}
	}
}

func (l *WebNodeList) SetStatus(path string, status NodeStatus) {
	l.mux.Lock()
	defer l.mux.Unlock()
	for i := range l.list {
		if l.list[i].path == path {
			l.setStatusByIndex(i, status)
			return
		}
	}
}

func (l *WebNodeList) setStatusByIndex(index int, status NodeStatus) {
	switch l.list[index].nodeStatus {
	case pending:
		l.pendingCount--
	case busy:
		l.busyCount--
	case done:
		l.doneCount--
	}
	switch status {
	case pending:
		l.pendingCount++
	case busy:
		l.busyCount++
	case done:
		l.doneCount++
	}
	l.list[index].nodeStatus = status
	return
}

func (l *WebNodeList) GetStats() (done, fail, total int) {
	return l.doneCount + l.failCount, l.failCount, l.doneCount + l.pendingCount + l.busyCount + l.failCount
}
