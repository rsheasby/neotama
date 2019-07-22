package main

import (
	"time"
	"regexp"
)

var dirParentRegex *regexp.Regexp
var fileParentRegex *regexp.Regexp

func init() {
	dirParentRegex = regexp.MustCompile(`(.*\/).*\/`)
	fileParentRegex = regexp.MustCompile(`(.*\/).*`)
}

type NodeType bool
const (
	file NodeType = false
	directory NodeType = true
)

type NodeStatus int8
const (
	pending NodeStatus = iota
	busy
	done
	failed
)

type WebNode struct {
	nodeStatus NodeStatus
	nodeType NodeType
	path string
	name string
	time *time.Time
	size int64
	description string
}

func (w *WebNode) GetParentPath() (result string) {
	var regexResult []string
	switch w.nodeType {
	case file:
		regexResult = fileParentRegex.FindStringSubmatch(w.path)
	case directory:
		regexResult = dirParentRegex.FindStringSubmatch(w.path)
	}

	if regexResult == nil {
		return ""
	} else {
		return regexResult[1]
	}
}
