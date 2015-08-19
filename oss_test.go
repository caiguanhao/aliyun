package main

import "testing"

func TestGetLastPartOfPath(t *testing.T) {
	if getLastPartOfPath("root") != "root" {
		t.Error("it should be root")
	}
	if getLastPartOfPath("/root/my/file") != "file" {
		t.Error("it should be file")
	}
}

func testWalkFiles(t *testing.T, _source []string, _target string, expects [][]string) {
	source = _source
	target = _target
	done := make(chan struct{})
	defer close(done)
	paths, _ := walkFiles(done)
	for path := range paths {
		exists := false
		for _, item := range expects {
			if item[0] == path[0] {
				if item[1] != path[1] {
					t.Errorf("%s should be %s", path[1], item[1])
				}
				exists = true
			}
		}
		if !exists {
			t.Errorf("%s should not exist", path[0])
		}
	}
}

func TestWalkFiles(t *testing.T) {
	testWalkFiles(t, []string{"test/fixtures/a/b/c"}, "/", [][]string{
		{"test/fixtures/a/b/c", "/c"},
	})

	testWalkFiles(t, []string{"test/fixtures/a/b/c"}, "/e", [][]string{
		{"test/fixtures/a/b/c", "/e"},
	})

	testWalkFiles(t, []string{"test/fixtures/a/b/c"}, "/e/", [][]string{
		{"test/fixtures/a/b/c", "/e/c"},
	})

	testWalkFiles(t, []string{"test/fixtures"}, "/", [][]string{
		{"test/fixtures/a/b/c", "/fixtures/a/b/c"},
		{"test/fixtures/d/e", "/fixtures/d/e"},
		{"test/fixtures/f", "/fixtures/f"},
	})

	testWalkFiles(t, []string{"test/fixtures"}, "/e", [][]string{
		{"test/fixtures/a/b/c", "/e/a/b/c"},
		{"test/fixtures/d/e", "/e/d/e"},
		{"test/fixtures/f", "/e/f"},
	})

	testWalkFiles(t, []string{"test/fixtures"}, "/e/", [][]string{
		{"test/fixtures/a/b/c", "/e/fixtures/a/b/c"},
		{"test/fixtures/d/e", "/e/fixtures/d/e"},
		{"test/fixtures/f", "/e/fixtures/f"},
	})

	testWalkFiles(t, []string{"test/fixtures/"}, "/e/", [][]string{
		{"test/fixtures/a/b/c", "/e/fixtures/a/b/c"},
		{"test/fixtures/d/e", "/e/fixtures/d/e"},
		{"test/fixtures/f", "/e/fixtures/f"},
	})

	testWalkFiles(t, []string{"test/fixtures/a", "test/fixtures/f"}, "/", [][]string{
		{"test/fixtures/a/b/c", "/a/b/c"},
		{"test/fixtures/f", "/f"},
	})

	testWalkFiles(t, []string{"test/fixtures/a", "test/fixtures/f"}, "/e", [][]string{
		{"test/fixtures/a/b/c", "/e/b/c"},
		{"test/fixtures/f", "/e/f"},
	})

	testWalkFiles(t, []string{"test/fixtures/a", "test/fixtures/f"}, "/e/", [][]string{
		{"test/fixtures/a/b/c", "/e/a/b/c"},
		{"test/fixtures/f", "/e/f"},
	})
}
