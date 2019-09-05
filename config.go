package main

import (
	"encoding/json"
	"github.com/akamensky/argparse"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type OutputFormat int
const (
	tree OutputFormat = iota
	list
	urlencoded
)

type JobConfig struct {
	url     string
	threads int
	depthLimit int
	pConfig ParserConfig
	outputFormat OutputFormat
}

type ParserConfig struct {
	Options struct {
		EnableDescription bool
		TimeFormat        string
	}
	Regex struct {
		LineMatch        string
		PathGroup        int
		TimeGroup        int
		SizeGroup        int
		DescriptionGroup int
	}
	CompiledRegexp *regexp.Regexp
}

func readParserConfig(filename string) (result ParserConfig) {
	fileContents, fileErr := ioutil.ReadFile(filename)
	if fileErr != nil {
		log.Fatalf("Could not read parser config file \"%s\".", filename)
	}
	jsonErr := json.Unmarshal(fileContents, &result)
	if jsonErr != nil {
		log.Fatalf("Parser config file \"%s\" is not valid JSON.", filename)
	}
	compRegexp, regexpErr := regexp.Compile(result.Regex.LineMatch)
	if regexpErr != nil {
		log.Fatalf("Invalid regex in parser config file.")
	}
	result.CompiledRegexp = compRegexp
	return
}

func ReadConfig() (config JobConfig) {
	parser := argparse.NewParser("", "Safely and quickly crawls a directory listing and outputs a pretty tree.")

	url := parser.String("u", "url", &argparse.Options{Required: true, Help: "URL to crawl"})
	threads := parser.Int("t", "threads", &argparse.Options{Required: false, Help: "Maximum number of concurrent connections", Default: 10})
	depthLimit := parser.Int("d", "depth", &argparse.Options{Required: false, Help: "Maximum depth to traverse. Depth of 0 means only query the provided URL. Value of -1 means unlimited", Default: -1})
	configFile := parser.String("c", "config", &argparse.Options{Required: false, Help: "Config file to use for parsing the directory listing"})
	outputFormat := parser.Selector("o", "output", []string{"tree", "list", "urlencoded"}, &argparse.Options{Required: false, Default: "tree", Help: "Output format of results"})

	// TODO: Not sure if you care about the error or how you want to log it: https://github.com/akamensky/argparse#usage
	if err := parser.Parse(os.Args); err != nil {
		log.Printf("Command line argument error: %s", err)
	}

	config.url = *url
	config.threads = *threads
	config.depthLimit = *depthLimit
	config.pConfig = readParserConfig(*configFile)
	switch *outputFormat {
	case "tree": config.outputFormat = tree
	case "list": config.outputFormat = list
	case "urlencoded": config.outputFormat = urlencoded
	}

	return
}
