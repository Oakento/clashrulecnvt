package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func writeConf(fullConf map[string]interface{}, newRules []string, output string) {
	fullConf["rules"] = newRules
	delete(fullConf, "rule-providers")
	data, err := yaml.Marshal(fullConf)
	if err != nil {
		log.Println("Have trouble with Marshal")
	}
	err = ioutil.WriteFile(output, data, 0644)
	if err != nil {
		log.Println("Have trouble saving modified configuration file.")
	}
	log.Printf("Modified configuration saved to %v", output)
}
