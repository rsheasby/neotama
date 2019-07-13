package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"time"
	"fmt"
)

func queryUrl(node webNode) (response queryRes) {
	url := node.path
	fmt.Println("Querying URL: " + url)
	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	response.url = url
	response.body = string(body)
	response.parent = node
	return
}

func startParser(config jobConfig) (input chan queryRes) {
	input = make(chan queryRes, config.json.Options.Threads)
	go func() {
		parseListing(<-input, config)
	}()
	return
}

func parseListing(listing queryRes, config jobConfig) (results []webNode) {
	// TODO: Selective description parsing
	res := config.compiledRegexp.FindAllStringSubmatch(listing.body, -1)
	results = make([]webNode, 0, len(res))
	for _, v := range(res) {
		node, skip := parseEntry(v, listing.url, config)
		if !skip {
			queueInsert(node)
		}
	}
	markNodeDone(listing.parent)
	return
}

func parseEntry(entry []string, parentUrl string, config jobConfig) (result webNode, skip bool) {
	e := splitEntry(entry, config)
	cleanTime := cleanHtml(e.time)
	if cleanTime == "" {
		skip = true
		return
	}
	result.path = parentUrl + e.path
	result.name, _ = url.QueryUnescape(e.path)
	if (lastChar(result.path) == "/") {
		result.nodeType = directory
		result.status = unexplored
	} else {
		result.nodeType = file
		result.status = explored
	}
	result.time, _ = time.Parse(config.json.Options.TimeFormat, cleanTime)
	if config.json.Options.EnableDescription {
		result.description = cleanHtml(e.description)
	} else {
		result.description = ""
	}
	result.size = parseSize(cleanHtml(e.size))
	return
}
