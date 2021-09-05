package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func writeClashConf(fullConf map[string]interface{}, newRules []string, output string) {
	fullConf["rules"] = newRules
	delete(fullConf, "rule-providers")
	data, err := yaml.Marshal(fullConf)
	if err != nil {
		log.Println("Have trouble with Marshal")
		os.Exit(1)
	}
	err = ioutil.WriteFile(output, data, 0644)
	if err != nil {
		log.Println("Have trouble saving modified configuration file.")
		os.Exit(1)
	}
	log.Printf("Modified configuration saved to %v", output)
}
