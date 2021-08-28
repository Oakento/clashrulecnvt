package main

import (
	"flag"
	"os"
	"path"
)

type RuleSet struct {
	name     string
	behavior string
	Payload  []string `yaml:"payload"`
}

var HOME string = os.Getenv("HOME")
var HTTP_PROXY string = os.Getenv("HTTP_PROXY")

func main() {

	var (
		help       bool
		inputFile  string
		outputFile string
	)

	DEFAULT_CONFIG_DIR := path.Join(HOME, ".config/clash/")
	DEFAULT_CONFIG := path.Join(DEFAULT_CONFIG_DIR, "config.yaml")

	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&inputFile, "i", DEFAULT_CONFIG, "set `input` raw configuration file")
	flag.StringVar(&outputFile, "o", path.Join(DEFAULT_CONFIG_DIR, "newconfig.yaml"), "set `output` path")

	flag.Parse()
	if help {
		flag.Usage()
	}

	ruleProviders, rawRules, fullConf := readConf(inputFile)
	chs := make([]chan RuleSet, len(ruleProviders))

	for i, provider := range ruleProviders {
		chs[i] = make(chan RuleSet)
		go fetchProvider(provider, chs[i])
	}

	rulesets := make(map[string]RuleSet)
	for _, ch := range chs {
		ruleset := <-ch
		rulesets[ruleset.name] = ruleset
	}
	parsedRules := parseRule(rulesets, rawRules)

	writeConf(fullConf, parsedRules, outputFile)
}
