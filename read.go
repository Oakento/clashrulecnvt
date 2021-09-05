package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ClashConf struct {
	RuleProviders map[string]map[string]interface{} `yaml:"rule-providers"`
	Rules         []string                          `yaml:"rules"`
}

type RuleProvider struct {
	name     string
	ptype    string
	behavior string
	path     string
	url      string
}

func readClashConf(input string) ([]RuleProvider, []string, map[string]interface{}) {

	yamlFile, err := ioutil.ReadFile(input)
	if err != nil {
		log.Printf("Error occurred while reading input file: %v", err)
		os.Exit(1)
	}
	var conf ClashConf
	var fullConf map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Printf("Error occurred while reading input file: %v", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &fullConf)
	if err != nil {
		log.Printf("Error occurred while reading input file: %v", err)
		os.Exit(1)
	}
	ruleProviders := conf.RuleProviders
	// ruleProviders := fullConf["rule-providers"].(map[string]map[string]interface{})
	if ruleProviders == nil {
		log.Printf("Did not find any rule providers.")
		os.Exit(0)
	}
	var providers []RuleProvider
	for key, value := range ruleProviders {
		var provider *RuleProvider = new(RuleProvider)
		provider.name = key
		ptype, _ := value["type"].(string)
		behavior, _ := value["behavior"].(string)
		path, _ := value["path"].(string)
		url, _ := value["url"].(string)
		provider.ptype = ptype
		provider.behavior = behavior
		provider.path = path
		provider.url = url

		providers = append(providers, *provider)
	}
	log.Println("Raw configuration read.")
	return providers, conf.Rules, fullConf
}
