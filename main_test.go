package main

import (
	"reflect"
	"testing"
)

func TestSortVersions(t *testing.T) {
	cases := []struct {
		name     string
		in       []string
		expected []string
	}{
		{
			name: "sort sermantic versions",
			in: []string{
				"1.0.0", "2.1.0", "1.2.0", "1.10.2", "1.10.20", "1.10.3",
			},
			expected: []string{
				"2.1.0", "1.10.20", "1.10.3", "1.10.2", "1.2.0", "1.0.0",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sortVersions(c.in)
			if !reflect.DeepEqual(c.expected, c.in) {
				t.Errorf("expected: %v, got: %v", c.expected, c.in)
			}
		})
	}
}
