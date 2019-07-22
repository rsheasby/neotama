package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"runtime"
	// "github.com/davecgh/go-spew/spew"
)

var _ = fmt.Println

func getUrl(url string) (html string) {
	// fmt.Println("Getting url", url)
	res, _ := http.Get(url)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}

func getListing(url string, pConfig ParserConfig) (nodes []WebNode) {
	html := getUrl(url)
	nodes = ParseDirList(html, url, pConfig)
	return
}

func Spider(job JobConfig) {
	sem := make(chan bool, job.threads)
	wnl := CreateWebNodeList()
	wnl.InsertSorted([]WebNode {WebNode {pending, directory, job.url, job.url, nil, 0, ""}})
	for {
		if wnl.IsDone() {
			// spew.Dump(wnl.list)
			for _, v := range(wnl.list) {
				fmt.Println(v.path)
			}
			return
		} else {
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
					wnl.InsertSorted(nodes)
					wnl.SetStatus(node.path, done)
				}(pending)
			}
		}
	}
}
