package main

import (
	"testing"
)

func TestPgCaseToGoCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"first_name", "FirstName"},
		{"id", "ID"},
		{"person_id", "PersonID"},
		{"person_ideal", "PersonIdeal"},
	}

	for i, tt := range tests {
		actual := pgCaseToGoCase(tt.input)
		if actual != tt.expected {
			t.Errorf(`%d. Given "%s", expected "%s", but got "%s"`, i, tt.input, tt.expected, actual)
		}
	}
}
