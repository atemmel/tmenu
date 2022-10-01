package main

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

	result := filter(input, alternatives)

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

	result := filter(input, alternatives)

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

	result := filter(input, alternatives)

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

	result := filter(input, alternatives)

	if len(result) != len(expected) {
		fail(t, expected, result)
	}

	for i := range result {
		if result[i] != expected[i] {
			fail(t, expected, result)
		}
	}
}
