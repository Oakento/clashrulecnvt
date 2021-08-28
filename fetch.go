package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"gopkg.in/yaml.v2"
)

func fetchProvider(provider RuleProvider, ch chan RuleSet) {
	if provider.ptype == "file" {
		ruleFile, err := ioutil.ReadFile(provider.path)
		if err != nil {
			log.Printf("Error occurred while reading local rule file: %v", err)
		}
		var rule RuleSet
		err = yaml.Unmarshal(ruleFile, &rule)
		if err != nil {
			ch <- RuleSet{name: provider.name, Payload: nil, behavior: ""}
			log.Printf("Error occurred while reading local rule file: %v", err)
		}
		rule.name = provider.name
		rule.behavior = provider.behavior
		ch <- rule

	} else if provider.ptype == "http" {
		var client *http.Client
		if HTTP_PROXY != "" {
			log.Printf("Start fetching remote data through http_proxy: %v", provider.url)
			proxyURL, err := url.Parse(HTTP_PROXY)
			if err != nil {
				log.Panic("Having problem parsing http_proxy")
			}
			client = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
			}
		} else {
			log.Printf("Start fetching remote data: %v", provider.url)
			client = &http.Client{}
		}
		req, err := http.NewRequest("GET", provider.url, nil)
		if err != nil {
			ch <- RuleSet{name: provider.name, Payload: nil, behavior: ""}
			log.Printf("Failed to fetch remote rule data due to connection failure")
		}
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("Accept-Language", "en-US,zh-TW;q=0.8,zh;q=0.5,en;q=0.3")
		// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Accept-Charset", "utf-8")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")

		resp, err := client.Do(req)

		if err != nil {
			ch <- RuleSet{name: provider.name, Payload: nil, behavior: ""}
			log.Printf("Failed to fetch remote rule data: %v", provider.url)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panicf("Error occurred while fetching remote data: %v", err)
		}
		var rule RuleSet
		err = yaml.Unmarshal(body, &rule)
		if err != nil {
			ch <- RuleSet{name: provider.name, Payload: nil, behavior: ""}
			log.Panicf("Error occurred while fetching remote data: %v", err)
		}
		rule.name = provider.name
		rule.behavior = provider.behavior
		ch <- rule
	}

}
