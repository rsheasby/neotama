package main

import (
	// "github.com/davecgh/go-spew/spew"
)

func main() {
	job := ReadConfig()
	Spider(job)
}
