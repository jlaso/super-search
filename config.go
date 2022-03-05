package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func cleanEmpties(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func cleanDuplicatesAndEmpties(inputList []string) []string {
	already := make(map[string]bool)
	var outputList []string
	for _, item := range inputList {
		if item != "" {
			if _, value := already[item]; !value {
				already[item] = true
				outputList = append(outputList, item)
			}
		}
	}
	return outputList
}

func readExcludeFile() *[]string {
	var result []string
	dirname, err := os.UserHomeDir()
	if err != nil {
		return &result
	}
	fname := dirname + "/.ss_except"
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return &result
	}
	result = cleanDuplicatesAndEmpties(strings.Split(string(content), "\n"))
	return &result
}
