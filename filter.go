package main

import (
	"strings"
)

func filter(input string, alternatives []string) []string {
	result := make([]string, 0, len(alternatives))

	for i := range alternatives {
		if matches(input, alternatives[i]) {
			result = append(result, alternatives[i])
		}
	}

	return result
}

func matches(matchString string, input string) bool {
	matchString = strings.ToUpper(matchString)
	input = strings.ToUpper(input)

	for i := range matchString {
		this := matchString[i]
		found := -1
		for j := range input {
			that := input[j]

			if this == that {
				found = j
				break
			}
		}

		if found == -1 {
			return false
		}

		input = strings.Replace(input, string(this), "", 1)
	}
	return true
}
