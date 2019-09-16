package main

import (
	"encoding/json"
	"github.com/akamensky/argparse"
	"github.com/junegunn/go-isatty"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type OutputFormat int

const (
	tree OutputFormat = iota
	list
	urlencoded
)

type ColorOption int

const (
	on ColorOption = iota
	off
	lol
)

type JobConfig struct {
	url          string
	threads      int
	retryLimit   int
	depthLimit   int
	noSort       bool
	showProgress bool
	colorOption  ColorOption
	colorValues  string
	pConfig      ParserConfig
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

func detectServer(url string, retryLimit int) (server string, fail bool) {
	for ; retryLimit >= 0; retryLimit-- {
		res, err := http.Head(url)
		if err != nil {
			continue
		}
		server := res.Header.Get("Server")
		if server == "" {
			return "", true
		}
		if strings.Contains(server, "Apache") {
			return "apache", false
		}
		if strings.Contains(server, "nginx") {
			return "nginx", false
		}
	}
	return "", true
}

func readParserConfig(filename string) (result ParserConfig) {
	fileContents, fileErr := ioutil.ReadFile(filename)
	if fileErr != nil {
		log.Fatalf("Could not read parser config file \"%s\".", filename)
	}
	return parseParserConfig(fileContents)
}

// I know the name is ass, but I can't think of a better one. If whoever is reading this can think of a better name, please open a PR.
func parseParserConfig(jsonConfig []byte) (result ParserConfig) {
	jsonErr := json.Unmarshal(jsonConfig, &result)
	if jsonErr != nil {
		log.Fatalf("Config does not contain valid JSON: %s", jsonErr)
	}
	compRegexp, regexpErr := regexp.Compile(result.Regex.LineMatch)
	if regexpErr != nil {
		log.Fatalf("Invalid regex in parser config file.")
	}
	result.CompiledRegexp = compRegexp
	return
}

func ReadConfig() (config JobConfig) {
	parser := argparse.NewParser("", "Safely and quickly crawls a directory listing, outputting a pretty tree.")

	url := parser.String("u", "url", &argparse.Options{Required: true, Help: "URL to crawl"})
	threads := parser.Int("t", "threads", &argparse.Options{Required: false, Help: "Maximum number of concurrent connections", Default: 10})
	retryLimit := parser.Int("r", "retry", &argparse.Options{Required: false, Help: "Maximum amount of times to retry a failed query", Default: 3})
	depthLimit := parser.Int("d", "depth", &argparse.Options{Required: false, Help: "Maximum depth to traverse. Depth of 0 means only query the provided URL. Value of -1 means unlimited", Default: -1})
	noSort := parser.Flag("", "disable-sorting", &argparse.Options{Required: false, Help: "Disables sorting. Default behavior is to sort by path alphabetically, with files above directories"})
	progress := parser.Selector("", "progress", []string{"auto", "on", "off"}, &argparse.Options{Required: false, Default: "auto", Help: "Whether to show the stderr progress bar or not"})
	color := parser.Selector("", "color", []string{"auto", "on", "off", "lol"}, &argparse.Options{Required: false, Default: "auto", Help: "Whether to output color codes or not. Color codes will be read from LS_COLORS if it exists, and will fallback to some basic defaults otherwise"})
	server := parser.Selector("s", "server", []string{"auto", "apache", "nginx"}, &argparse.Options{Required: false, Default: "auto", Help: "Server type to use for parsing. Auto will detect the server based on the HTTP headers"})
	configFile := parser.String("p", "parser-config", &argparse.Options{Required: false, Help: "Config file to use for parsing the directory listing"})
	outputFormat := parser.Selector("o", "output", []string{"tree", "list", "urlencoded"}, &argparse.Options{Required: false, Default: "tree", Help: "Output format of results"})

	if err := parser.Parse(os.Args); err != nil {
		log.Fatalf("Command line argument error: %s", err)
	}

	config.url = *url
	if config.url[len(config.url)-1] != '/' {
		config.url += "/"
	}
	config.threads = *threads
	config.retryLimit = *retryLimit
	config.depthLimit = *depthLimit
	config.noSort = *noSort
	if *configFile != "" {
		config.pConfig = readParserConfig(*configFile)
	} else {
		if *server == "auto" {
			server, err := detectServer(config.url, config.retryLimit)
			if err {
				log.Fatalf("Unable to detect the server type. Please manually specify the server type using -s or specify a config file using -p.")
			} else {
				config.pConfig = BuiltinConfigs[server]
			}
		} else {
			config.pConfig = BuiltinConfigs[*server]
		}
	}

	switch *progress {
	case "auto":
		config.showProgress = isatty.IsTerminal(os.Stderr.Fd())
	case "on":
		config.showProgress = true
	case "off":
		config.showProgress = false
	}

	switch *color {
	case "auto":
		if isatty.IsTerminal(os.Stdout.Fd()) {
			config.colorOption = on
		} else {
			config.colorOption = off
		}
	case "on":
		config.colorOption = on
	case "off":
		config.colorOption = off
	case "lol":
		config.colorOption = lol
	}

	switch *outputFormat {
	case "tree":
		config.outputFormat = tree
	case "list":
		config.outputFormat = list
	case "urlencoded":
		config.outputFormat = urlencoded
	}

	return
}
