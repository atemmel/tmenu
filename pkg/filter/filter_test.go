package filter

import (
	"testing"
)

func fail(t *testing.T, expected []string, result []string) {
	t.Fatal("Expected output and actual output mismatch, expected:", expected, "result:", result)
}

func TestEmptyFilter(t *testing.T) {
	input := ""
	alternatives := []string{
		"a",
		"b",
		"c",
	}

	expected := []string{
		"a",
		"b",
		"c",
	}

	result := Filter(input, alternatives)

	if len(result) != len(expected) {
		fail(t, expected, result)
	}

	for i := range result {
		if result[i] != expected[i] {
			fail(t, expected, result)
		}
	}
}

func TestFilter1(t *testing.T) {
	input := "a"
	alternatives := []string{
		"a",
		"b",
		"c",
	}

	expected := []string{
		"a",
	}

	result := Filter(input, alternatives)

	if len(result) != len(expected) {
		fail(t, expected, result)
	}

	for i := range result {
		if result[i] != expected[i] {
			fail(t, expected, result)
		}
	}
}

func TestFilter2(t *testing.T) {
	input := "A"
	alternatives := []string{
		"a",
		"b",
		"c",
	}

	expected := []string{
		"a",
	}

	result := Filter(input, alternatives)

	if len(result) != len(expected) {
		fail(t, expected, result)
	}

	for i := range result {
		if result[i] != expected[i] {
			fail(t, expected, result)
		}
	}
}


func TestFilter3(t *testing.T) {
	input := "Aa"
	alternatives := []string{
		"a",
		"b",
		"c",
	}

	expected := []string{
	}

	result := Filter(input, alternatives)

	if len(result) != len(expected) {
		fail(t, expected, result)
	}

	for i := range result {
		if result[i] != expected[i] {
			fail(t, expected, result)
		}
	}
}

func TestFilter4(t *testing.T) {
	input := "al"
	alternatives := []string{
		"Alfons",
		"Alban",
		"Niklas",
		"Oskar",
	}

	expected := []string{
		"Alfons",
		"Alban",
		"Niklas",
	}

	result := Filter(input, alternatives)

	if len(result) != len(expected) {
		fail(t, expected, result)
	}

	for i := range result {
		if result[i] != expected[i] {
			fail(t, expected, result)
		}
	}
}

func TestFilter5(t *testing.T) {
	input := "asdf"
	alternatives := []string{
		"Alfons",
		"Alban",
		"Niklas",
		"Oskar",
	}

	expected := []string{
	}

	result := Filter(input, alternatives)

	if len(result) != len(expected) {
		fail(t, expected, result)
	}

	for i := range result {
		if result[i] != expected[i] {
			fail(t, expected, result)
		}
	}
}
