package main

import (
	"time"
)

type nodeType bool

const (
	directory nodeType = false
	file nodeType = true
)

type nodeStatus int8

const (
	unexplored nodeStatus = iota
	inProgress
	explored
)

type queryRes struct {
	url string
	body string
	parent webNode
}

type unparsedNode struct {
	path string
	time string
	size string
	description string
}

type webNode struct {
	nodeType nodeType
	status nodeStatus
	path string
	name string
	time time.Time
	size int64
	description string
}
