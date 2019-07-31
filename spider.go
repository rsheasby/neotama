package main

import (
	"net/http"
	"io/ioutil"
	// "fmt"
	"runtime"
)

func getUrl(url string) (html string) {
	for i := 0; i < 3; i++ {
		res, err := http.Get(url)
		if err != nil {
			continue
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return string(body)
	}
	return ""
}

func getListing(url string, pConfig ParserConfig) (nodes []WebNode) {
	html := getUrl(url)
	nodes = ParseDirList(html, url, pConfig)
	return
}

func Spider(job JobConfig) {
	sem := make(chan bool, job.threads)
	wnl := CreateWebNodeList()
	wnl.InsertSorted([]WebNode {WebNode {pending, directory, false, job.url, job.url, nil, 0, ""}}, "", true)
	for {
		if wnl.IsDone() {
			wnl.PrintDone()
			return
		} else {
			wnl.PrintDone()
			sem<-true
			pending, wait := wnl.GetPending()
			if wait {
				<-sem
				runtime.Gosched()
				continue
			} else {
				go func(node WebNode) {
					html := getUrl(node.path)
					<-sem
					nodes := ParseDirList(html, node.path, job.pConfig)
					wnl.InsertSorted(nodes, node.path, true)
					wnl.SetStatus(node.path, done)
				}(pending)
			}
		}
	}
}
