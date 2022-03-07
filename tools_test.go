package main

import "testing"

func TestStringInSlice(t *testing.T) {
	var sample = []string{
		"abc",
		"def",
		"hij",
	}
	if !stringInSlice("abc", &sample) {
		t.Errorf("Unexpected result for stringInSlice(%s)", "abc")
	}
	if stringInSlice("cba", &sample) {
		t.Errorf("Unexpected result for stringInSlice(%s)", "cba")
	}
}
