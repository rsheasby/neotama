package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"github.com/akamensky/argparse"
	"regexp"
)

type jobConfig struct {
	urls []string
	configFile string
	compiledRegexp *regexp.Regexp
	json struct {
		Options struct {
			EnableDescription bool
			Threads int
			TimeFormat string
		}
		Regex struct {
			LineMatch string
			PathGroup int
			TimeGroup int
			SizeGroup int
			DescriptionGroup int
		}
	}
}

func parseArgs() (config jobConfig) {
	parser := argparse.NewParser("", "Safely and quickly crawls a directory listing and outputs a pretty tree.")

	urls := parser.List("u", "url", &argparse.Options{Required: true, Help: "URLs to crawl"})
	configFile := parser.String("c", "config", &argparse.Options{Required: false, Help: "Config file to use for parsing the directory listing."})

	parser.Parse(os.Args)

	config.urls = *urls
	config.configFile = *configFile
	return
}

func parseConfigFile(config *jobConfig) {
	// TODO: Add error handling
	fileContents, _ := ioutil.ReadFile(config.configFile)
	json.Unmarshal([]byte(fileContents), &config.json)
	config.compiledRegexp, _ = regexp.Compile(config.json.Regex.LineMatch)
}
