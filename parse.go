package main

import (
	"strings"
)

type RuleEntity struct {
	index    int
	name     string
	strategy string
}

func parseRule(rulesets map[string]RuleSet, rawRules []string) []string {
	var ruleentities []RuleEntity
	for i, value := range rawRules {
		rulesplit := strings.Split(value, ",")
		if rulesplit[0] == "RULE-SET" {
			ruleentities = append(ruleentities, RuleEntity{
				index:    i,
				name:     rulesplit[1],
				strategy: rulesplit[2],
			})
		}
	}
	newRules := make([]string, len(rawRules))
	copy(newRules, rawRules)

	for i := len(ruleentities) - 1; i >= 0; i-- {
		leftRules := newRules[:ruleentities[i].index]
		rightRules := newRules[ruleentities[i].index+1:]
		swapRules := rulesets[ruleentities[i].name].Payload
		prefix := ""
		if rulesets[ruleentities[i].name].behavior == "ipcidr" {
			prefix = "IP-CIDR,"
		} else if rulesets[ruleentities[i].name].behavior == "domain" {
			prefix = "DOMAIN,"
		}

		suffixes := "," + ruleentities[i].strategy
		var newSwapRules []string
		for _, v := range swapRules {
			newSwapRules = append(newSwapRules, prefix+v+suffixes)
		}
		newRules = append(leftRules, newSwapRules...)
		newRules = append(newRules, rightRules...)
	}
	return newRules
}
