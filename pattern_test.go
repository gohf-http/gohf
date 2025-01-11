package gohf

import (
	"fmt"
	"testing"
)

func TestPattern(t *testing.T) {
	type result struct {
		pattern pattern
		err     bool
	}

	tests := map[string]result{
		"GET google.com/foo/bar/": {
			pattern: pattern{
				method: "GET",
				host:   "google.com",
				path:   "foo/bar/",
			},
			err: false,
		},
		"GET google.com/foo/bar": {
			pattern: pattern{
				method: "GET",
				host:   "google.com",
				path:   "foo/bar",
			},
			err: false,
		},
		"google.com/foo/bar": {
			pattern: pattern{
				method: "",
				host:   "google.com",
				path:   "foo/bar",
			},
			err: false,
		},
		"/foo/bar": {
			pattern: pattern{
				method: "",
				host:   "",
				path:   "foo/bar",
			},
			err: false,
		},
		"google.com/": {
			pattern: pattern{
				method: "",
				host:   "google.com",
				path:   "",
			},
			err: false,
		},
		"  ": {
			pattern: pattern{},
			err:     true,
		},
		"GET google.com": {
			pattern: pattern{},
			err:     true,
		},
	}

	for patternString, result := range tests {
		t.Run("test parsePattern: "+patternString, func(t *testing.T) {
			pattern, err := parsePattern(patternString)
			if got, want := err != nil, result.err; got != want {
				t.Errorf("expect error exists:%v got:%v", want, err)
			}
			if got, want := pattern, result.pattern; got != want {
				t.Errorf("pattern mismatch. got:%v want:%v", got, want)
			}
		})
	}
}

func TestMergePattern(t *testing.T) {
	type testcase struct {
		p1 pattern
		p2 pattern
	}
	type result struct {
		pattern pattern
		err     bool
	}

	tests := map[testcase]result{
		{
			p1: pattern{
				method: "GET",
				host:   "google.com",
				path:   "foo",
			},
			p2: pattern{
				method: "",
				host:   "",
				path:   "bar",
			},
		}: {
			pattern: pattern{
				method: "GET",
				host:   "google.com",
				path:   "foo/bar",
			},
			err: false,
		},
		{
			p1: pattern{
				method: "GET",
				host:   "google.com",
				path:   "foo",
			},
			p2: pattern{
				method: "GET",
				host:   "",
				path:   "",
			},
		}: {
			pattern: pattern{
				method: "GET",
				host:   "google.com",
				path:   "foo",
			},
			err: false,
		},
		{
			p1: pattern{
				method: "GET",
				host:   "google.com",
				path:   "foo",
			},
			p2: pattern{
				method: "PUT",
				host:   "",
				path:   "",
			},
		}: {
			pattern: pattern{},
			err:     true,
		},
	}

	i := 0
	for testcase, result := range tests {
		i++
		name := fmt.Sprintf("Test mergePattern %d", i)
		t.Run(name, func(t *testing.T) {
			pattern, err := mergePattern(testcase.p1, testcase.p2)
			if got, want := err != nil, result.err; got != want {
				t.Errorf("expect error exists:%v got:%v", want, err)
			}
			if got, want := pattern, result.pattern; got != want {
				t.Errorf("pattern mismatch. got:%v want:%v", got, want)
			}
		})
	}
}

func TestPatternString(t *testing.T) {
	tests := map[pattern]string{
		{
			method: "GET",
			host:   "google.com",
			path:   "foo",
		}: "GET google.com/foo",
		{
			method: "GET",
			host:   "google.com",
			path:   "/foo",
		}: "GET google.com/foo",
		{
			method: "GET",
			host:   "google.com",
			path:   "//foo",
		}: "GET google.com/foo",
		{
			method: "",
			host:   "google.com",
			path:   "//foo",
		}: "google.com/foo",
		{
			method: "",
			host:   "",
			path:   "//foo",
		}: "/foo",
		{
			method: "",
			host:   "",
			path:   "",
		}: "/",
	}

	i := 0
	for testcase, result := range tests {
		i++
		name := fmt.Sprintf("Test pattern.String %d", i)
		t.Run(name, func(t *testing.T) {
			s := testcase.String()
			if got, want := s, result; got != want {
				t.Errorf("pattern string mismatch. got:%v want:%v", got, want)
			}
		})
	}
}
