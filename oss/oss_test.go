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

func TestWalkFiles__pathsForFile(t *testing.T) {
	testWalkFiles(t, []string{"test/fixtures/a/b/c"}, "/", [][]string{
		{"test/fixtures/a/b/c", "/c"},
	})

	testWalkFiles(t, []string{"test/fixtures/a/b/c"}, "/e", [][]string{
		{"test/fixtures/a/b/c", "/e"},
	})

	testWalkFiles(t, []string{"test/fixtures/a/b/c"}, "/e/", [][]string{
		{"test/fixtures/a/b/c", "/e/c"},
	})
}

func TestWalkFiles__pathsForDirectory(t *testing.T) {
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
}

func TestWalkFiles__mix(t *testing.T) {
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

func TestWalkFiles__parents__pathsForFile(t *testing.T) {
	parentsPath = true

	testWalkFiles(t, []string{"test/fixtures/a/b/c"}, "/", [][]string{
		{"test/fixtures/a/b/c", "/test/fixtures/a/b/c"},
	})

	testWalkFiles(t, []string{"test/fixtures/a/b/c"}, "/e", [][]string{
		{"test/fixtures/a/b/c", "/e/test/fixtures/a/b/c"},
	})

	testWalkFiles(t, []string{"test/fixtures/a/b/c"}, "/e/", [][]string{
		{"test/fixtures/a/b/c", "/e/test/fixtures/a/b/c"},
	})
}

func TestWalkFiles__parents__pathsForDirectory(t *testing.T) {
	parentsPath = true

	testWalkFiles(t, []string{"test/fixtures"}, "/", [][]string{
		{"test/fixtures/a/b/c", "/test/fixtures/a/b/c"},
		{"test/fixtures/d/e", "/test/fixtures/d/e"},
		{"test/fixtures/f", "/test/fixtures/f"},
	})

	testWalkFiles(t, []string{"test/fixtures"}, "/e", [][]string{
		{"test/fixtures/a/b/c", "/e/test/fixtures/a/b/c"},
		{"test/fixtures/d/e", "/e/test/fixtures/d/e"},
		{"test/fixtures/f", "/e/test/fixtures/f"},
	})

	testWalkFiles(t, []string{"test/fixtures"}, "/e/", [][]string{
		{"test/fixtures/a/b/c", "/e/test/fixtures/a/b/c"},
		{"test/fixtures/d/e", "/e/test/fixtures/d/e"},
		{"test/fixtures/f", "/e/test/fixtures/f"},
	})

	testWalkFiles(t, []string{"test/fixtures/"}, "/e/", [][]string{
		{"test/fixtures/a/b/c", "/e/test/fixtures/a/b/c"},
		{"test/fixtures/d/e", "/e/test/fixtures/d/e"},
		{"test/fixtures/f", "/e/test/fixtures/f"},
	})
}

func TestWalkFiles__parents__mix(t *testing.T) {
	parentsPath = true

	testWalkFiles(t, []string{"test/fixtures/a", "test/fixtures/f"}, "/", [][]string{
		{"test/fixtures/a/b/c", "/test/fixtures/a/b/c"},
		{"test/fixtures/f", "/test/fixtures/f"},
	})

	testWalkFiles(t, []string{"test/fixtures/a", "test/fixtures/f"}, "/e", [][]string{
		{"test/fixtures/a/b/c", "/e/test/fixtures/a/b/c"},
		{"test/fixtures/f", "/e/test/fixtures/f"},
	})

	testWalkFiles(t, []string{"test/fixtures/a", "test/fixtures/f"}, "/e/", [][]string{
		{"test/fixtures/a/b/c", "/e/test/fixtures/a/b/c"},
		{"test/fixtures/f", "/e/test/fixtures/f"},
	})
}

func testHumanBytes(t *testing.T, size float64, expected string) {
	human := humanBytes(size)
	if human != expected {
		t.Errorf("human bytes error when size is %f: %s should be %s", size, human, expected)
	}
}

func TestHumanBytes(t *testing.T) {
	testHumanBytes(t, 1, "1 bytes")
	testHumanBytes(t, 512, "512 bytes")
	testHumanBytes(t, 1024, "1 KB")
	testHumanBytes(t, 1920, "1.875 KB")
	testHumanBytes(t, 1920*10+1, "18.751 KB")
	testHumanBytes(t, 1024*1023, "1023 KB")
	testHumanBytes(t, -1024*1023, "-1023 KB")
	testHumanBytes(t, 1024*1024, "1 MB")
	testHumanBytes(t, 1024*1024*1024, "1 GB")
	testHumanBytes(t, 1024*1024*1024*1024, "1 TB")
	testHumanBytes(t, 1024*1024*1024*1024*1024, "1024 TB")
}
