package main

import (
	"regexp"
	"strings"
)

func stringInSlice(a string, list *[]string) bool {
	if list == nil || len(*list) == 0 {
		return false
	}
	for _, b := range *list {
		if b == a {
			return true
		}
		if strings.Contains(a, b) {
			return true
		}
		r, _ := regexp.Compile(b)
		if r.FindString(a) != "" {
			return true
		}
	}
	return false
}
