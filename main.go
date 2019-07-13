package main

import (
	"time"
	"runtime"
	"fmt"
)

func main() {
	config := parseArgs()
	parseConfigFile(&config)
	parseChan := startParser(config)
	sem := make(chan bool, config.json.Options.Threads)
	for _, v := range(config.urls) {
		queueInsert(webNode{
			directory,
			unexplored,
			v,
			v,
			time.Time{},
			0,
			"",
		})
	}
	for {
		sem<-true // Punch in
		// fmt.Println("Starting Goroutine")
		n, wait, done := queuePop()
		if wait {
			<-sem
			runtime.Gosched()
		} else if done {
			for _, v := range(nodeList) {
				fmt.Printf("%#v\n", v)
			}
			<-sem
			return
		} else {
			go func(node webNode){
				fmt.Println("Starting Goroutine")
				parseChan<-queryUrl(node)
				fmt.Println("Finishing Goroutine")
				<-sem
			}(n)
		}
	}
}
