package main

import (
	"io/ioutil"
	"net/http"
	"runtime"
)

func getUrl(url string) (body string) {
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

func Spider(job JobConfig) {
	sem := make(chan bool, job.threads)
	wnl := CreateWebNodeList()
	wp := CreateWnlPrinter(&wnl, job.outputFormat)
	wnl.InsertSorted([]WebNode{{pending, directory, false, 0, true, job.url, job.url, nil, 0, ""}}, "", false)
	for {
		if wnl.IsDone() {
			wp.PrintDone()
			return
		}
		wp.PrintDone()
		sem <- true
		pending, wait := wnl.GetPending()
		if wait {
			<-sem
			runtime.Gosched()
			continue
		} else {
			go func(node WebNode) {
				html := getUrl(node.path)
				<-sem
				nodes := ParseDirList(html, node.path, node.nodeDepth+1, job.pConfig)
				wnl.InsertSorted(nodes, node.path, true)
				wnl.SetStatus(node.path, done)
			}(pending)
		}
	}
}
