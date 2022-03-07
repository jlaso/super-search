package main

import "testing"

func TestClean(t *testing.T) {
	var sample = []string{
		"",
		"",
		"hello",
		"bye",
		"hello",
	}
	result := cleanEmpties(sample)
	if len(result) != 3 {
		t.Errorf("Unexpected length of cleanEmpties(sample) = %d", len(result))
	}
	result = cleanDuplicatesAndEmpties(sample)
	if len(result) != 2 {
		t.Errorf("Unexpected length of cleanDuplicatesAndEmpties(sample) = %d", len(result))
	}
}
