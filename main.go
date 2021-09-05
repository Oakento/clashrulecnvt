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
var DEFAULT_CONFIG_DIR string = path.Join(HOME, ".config/clashrulecnvt/")
var DEFAULT_CONFIG string = path.Join(DEFAULT_CONFIG_DIR, "config.yaml")
var DEFAULT_CLASH_CONFIG_DIR string = path.Join(HOME, ".config/clash/")
var DEFAULT_CLASH_CONFIG string = path.Join(DEFAULT_CLASH_CONFIG_DIR, "config.yaml")

func main() {

	var (
		inputFile     string
		outputFile    string
		help          bool
		cmdConfigPath string
		cmdInputFile  string
		cmdOutputFile string
		cmdProxy      string
	)

	conf := Conf{
		RawConfigFile:    DEFAULT_CLASH_CONFIG,
		OutputConfigFile: path.Join(DEFAULT_CLASH_CONFIG_DIR, "newconfig.yaml"),
		HttpProxy:        HTTP_PROXY,
	}
	initConf(conf)

	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&cmdConfigPath, "c", "", "set `config` to read.")
	flag.StringVar(&cmdInputFile, "i", "", "set `input` raw configuration file")
	flag.StringVar(&cmdOutputFile, "o", "", "set `output` path")
	flag.StringVar(&cmdProxy, "p", "", "set `http proxy` to download rule files.")

	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	if cmdConfigPath != "" {
		conf = readConf(cmdConfigPath)
	} else {
		conf = readConf(DEFAULT_CONFIG)
	}
	inputFile = conf.RawConfigFile
	outputFile = conf.OutputConfigFile
	HTTP_PROXY = conf.HttpProxy

	if cmdInputFile != "" {
		inputFile = cmdInputFile
	}

	if cmdOutputFile != "" {
		outputFile = cmdOutputFile
	}

	if cmdProxy != "" {
		HTTP_PROXY = cmdProxy
	}

	ruleProviders, rawRules, fullConf := readClashConf(inputFile)
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

	writeClashConf(fullConf, parsedRules, outputFile)
}
