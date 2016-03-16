package main

import (
	"testing"
)

func test_parseArgsForOSSDownload(t *testing.T, args, expectedRemotes, expectedLocals []string) {
	actualRemotes, actualLocals := parseArgsForOSSDownload(args)
	for i, actuals := range [][]string{actualRemotes, actualLocals} {
		for _, actual := range actuals {
			found := false
			for _, expected := range map[int][]string{0: expectedRemotes, 1: expectedLocals}[i] {
				if expected == actual {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s is unexpected", actual)
			}
		}
	}
}

func TestParseArgsForOSSDownload(t *testing.T) {
	test_parseArgsForOSSDownload(t, []string{"/"}, []string{"/"}, []string{})
	test_parseArgsForOSSDownload(t, []string{"oss.go"}, []string{"/oss.go"}, []string{})
	test_parseArgsForOSSDownload(t, []string{"//oss//main.go"}, []string{"/oss/main.go"}, []string{})
	test_parseArgsForOSSDownload(t, []string{"//oss//main.go", "foo"}, []string{"/oss/main.go"}, []string{"foo"})
	test_parseArgsForOSSDownload(t, []string{"//1//1.go", "/2/2.go", "test"}, []string{"/1/1.go", "/2/2.go"}, []string{"test/1.go", "test/2.go"})
}
