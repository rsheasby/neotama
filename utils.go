package main

import (
	"regexp"
	"strings"
	"math"
	"strconv"
	"html"
)

var sizeMatcher *regexp.Regexp

func init() {
	sizeMatcher = regexp.MustCompile(`(?m)^([0-9\.]*)(.*)$`)
}

func lastChar(input string) (c string) {
	// TODO: Error handling
	c = input[len(input)-1:]
	return
}

func parseSize(input string) (bytes int64) {
	matches := sizeMatcher.FindStringSubmatch(input)
	// TODO: Error handling
	var multiplier float64 = 1000
	if len(matches[1]) > 0 {
		if strings.ContainsAny(matches[1], "iI") {
			multiplier = 1024
		}
		switch {
			case strings.ContainsAny(matches[1], "kK"): multiplier = math.Pow(multiplier, 1)
			case strings.ContainsAny(matches[1], "mM"): multiplier = math.Pow(multiplier, 2)
			case strings.ContainsAny(matches[1], "gG"): multiplier = math.Pow(multiplier, 3)
			case strings.ContainsAny(matches[1], "tT"): multiplier = math.Pow(multiplier, 4)
			case strings.ContainsAny(matches[1], "pP"): multiplier = math.Pow(multiplier, 5)
			case strings.ContainsAny(matches[1], "zZ"): multiplier = math.Pow(multiplier, 6)
			default: multiplier = 1
		}
	}
	// TODO: Error handling
	floatBytes, _ := strconv.ParseFloat(matches[0], 64)
	bytes = int64(floatBytes * multiplier)
	return
}

func cleanHtml(input string) (result string) {
	result = html.UnescapeString(input)
	result = strings.TrimSpace(result)
	return
}

func insertWebNode(slice *[]webNode, node webNode, index int) {
	arr := *slice
	arr = append(arr, webNode{})
	copy(arr[index+1:], arr[index:])
	arr[index] = node
	*slice = arr
}

func splitEntry(entry []string, config jobConfig) (result unparsedNode) {
	result.path = entry[config.json.Regex.PathGroup]
	result.time = entry[config.json.Regex.TimeGroup]
	result.size = entry[config.json.Regex.SizeGroup]
	result.description = entry[config.json.Regex.DescriptionGroup]
	return
}
