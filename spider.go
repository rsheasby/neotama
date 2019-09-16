package main

import (
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

func getUrl(url string, retryLimit int) (body string, fail bool) {
	for ; retryLimit >= 0; retryLimit-- {
		res, err := http.Get(url)
		if err != nil {
			time.Sleep(time.Millisecond * 333)
			continue
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return string(body), false
	}
	return "", true
}

func Spider(job JobConfig) {
	sem := make(chan bool, job.threads)
	wnl := CreateWebNodeList()
	wnl.InsertSorted([]WebNode{{pending, directory, false, 0, true, job.url, job.url, nil, 0, ""}}, "", false)
	pbw, pbwo, pbwe := CreateProgressBarWriter(&wnl, job.showProgress, job.colorOption == lol)
	wp := CreateWnlPrinter(&wnl, pbwo, job.outputFormat, job.colorOption)
	pbwe.Write([]byte(""))
	for {
		if wnl.IsDone() {
			wp.PrintDone()
			pbw.HideBar()
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
				html, fail := getUrl(node.path, job.retryLimit)
				<-sem
				if fail {
					node.nodeFail = true
					wnl.SetFailed(node.path)
				} else {
					nodes := ParseDirList(html, node.path, node.nodeDepth+1, job.depthLimit != -1 && node.nodeDepth >= job.depthLimit, job.pConfig)
					wnl.InsertSorted(nodes, node.path, !job.noSort)
				}
				wnl.SetStatus(node.path, done)
			}(pending)
		}
	}
}
