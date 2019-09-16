package main

var BuiltinConfigs map[string]ParserConfig

func init() {
	BuiltinConfigs = make(map[string]ParserConfig)
	BuiltinConfigs["apache"] = parseParserConfig([]byte(`{"options": {"enableDescription": true, "timeFormat": "2006-01-02 15:04"}, "regex": {"lineMatch": "(?mU)^<tr>.*<a href=\"(.*)\">.*<\\\/a>.*<td.*>(.*)<\\\/td><td.*>(.*)<\\\/td><td>(.*)<\\\/td><\\\/tr>$", "pathGroup": 1, "timeGroup": 2, "sizeGroup": 3, "descriptionGroup": 4}}`))
	BuiltinConfigs["nginx"] = parseParserConfig([]byte(`{"options": {"enableDescription": false, "timeFormat": "02-Jan-2006 15:04"}, "regex": {"lineMatch": "(?mU)^<a href=\"([^\"]*)\">.*<\\\/a> *([^ ]* [^ ]*)\\s*(-|\\d+)\\s$", "pathGroup": 1, "timeGroup": 2, "sizeGroup": 3, "descriptionGroup": 0}}`))
}
