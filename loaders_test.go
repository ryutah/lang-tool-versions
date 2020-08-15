package main

import (
	"reflect"
	"testing"
)

func TestExtractAsVersions(t *testing.T) {
	type in struct {
		tags []tagResponse
		rule string
	}
	cases := []struct {
		name     string
		in       in
		expected []string
	}{
		{
			name: `rule: ^go(\d+\.\d+(?:\.\d+)?)$`,
			in: in{
				tags: []tagResponse{
					{Name: "go1.10.1"},
					{Name: "go1.15"},
					{Name: "notmatch-1.15"},
				},
				rule: `^go(\d+\.\d+(?:\.\d+)?)$`,
			},
			expected: []string{"1.10.1", "1.15"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gots := extractAsVersions(c.in.tags, c.in.rule)
			if !reflect.DeepEqual(gots, c.expected) {
				t.Errorf("want: %v. got: %v", c.expected, gots)
			}
		})
	}
}
