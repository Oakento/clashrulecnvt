package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	RawConfigFile    string `yaml:"raw-config-file"`
	OutputConfigFile string `yaml:"output-config-file"`
	HttpProxy        string `yaml:"http-proxy"`
}

func initConf(conf Conf) {
	if _, err := os.Stat(DEFAULT_CONFIG); err != nil {
		err := os.MkdirAll(DEFAULT_CONFIG_DIR, 0755)

		if err != nil {
			log.Printf("Failed to create config directory.")
			os.Exit(1)
		}

		confYAML, err := yaml.Marshal(&conf)
		log.Println(string(confYAML))
		if err != nil {
			log.Println("Wrong format of configuration.")
			os.Exit(1)
		}
		err = ioutil.WriteFile(DEFAULT_CONFIG, confYAML, 0644)
		if err != nil {
			log.Println("Error occurred while saving default configuration file.")
			os.Exit(1)
		}
	}
}

func readConf(confPath string) Conf {
	if _, err := os.Stat(confPath); err != nil {
		log.Println("No such configuration file.")
		os.Exit(1)
	}
	confFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Println("IO error: cannot read the file.")
		os.Exit(1)
	}
	var conf Conf
	err = yaml.Unmarshal(confFile, &conf)
	if err != nil {
		log.Println("Format error with config file.")
		os.Exit(1)
	}
	return conf
}
