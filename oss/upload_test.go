package main

import (
	"os"
	"path/filepath"
	"testing"
)

func test_parseArgsForOSSUpload(t *testing.T, args []string, expectedLocals []string, expectedRemote string) {
	actualLocals, actualRemote := parseArgsForOSSUpload(args)
	for _, local := range actualLocals {
		found := false
		for _, expected := range expectedLocals {
			if expected == local {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s is unexpected local", local)
		}
	}
	if actualRemote != expectedRemote {
		t.Errorf("%s should be %s", actualRemote, expectedRemote)
	}
}

func TestParseArgsForOSSUpload(t *testing.T) {
	test_parseArgsForOSSUpload(t, []string{"/"}, []string{}, "/")
	test_parseArgsForOSSUpload(t, []string{"file"}, []string{}, "/file")
	test_parseArgsForOSSUpload(t, []string{"local", "remote"}, []string{"local"}, "/remote")
	test_parseArgsForOSSUpload(t, []string{"local", "///remote//file////location/"}, []string{"local"}, "/remote/file/location/")
	test_parseArgsForOSSUpload(t, []string{"local", "local2//file", "/remote/"}, []string{"local", "local2/file"}, "/remote/")
}

func TestGetLastPartOfPath(t *testing.T) {
	if getLastPartOfPath("root") != "root" {
		t.Error("it should be root")
	}
	if getLastPartOfPath("/root/my/file") != "file" {
		t.Error("it should be file")
	}
}

func test_localPathToRemotePath(t *testing.T, locals []string, remote string, parentsPath bool, expects map[string]string) {
	for _, local := range locals {
		expectedRemote := expects[local]
		if actualRemote := localPathToRemotePath(local, remote, len(locals) > 1, parentsPath); expectedRemote != actualRemote {
			t.Errorf("%s should be %s", actualRemote, expectedRemote)
		}
	}
}

func TestLocalPathToRemotePath(t *testing.T) {
	test_localPathToRemotePath(t, []string{"test/fixtures/a/b/c"}, "/", false, map[string]string{
		"test/fixtures/a/b/c": "/c",
	})

	test_localPathToRemotePath(t, []string{"test/fixtures/a/b/c"}, "/e", false, map[string]string{
		"test/fixtures/a/b/c": "/e",
	})

	test_localPathToRemotePath(t, []string{"test/fixtures/a/b/c"}, "/e/", false, map[string]string{
		"test/fixtures/a/b/c": "/e/c",
	})

	test_localPathToRemotePath(t, []string{"test/fixtures/a/b/c", "test/fixtures/f"}, "/e/", false, map[string]string{
		"test/fixtures/a/b/c": "/e/c",
		"test/fixtures/f":     "/e/f",
	})

	test_localPathToRemotePath(t, []string{"test/fixtures/a/b/c"}, "/", true, map[string]string{
		"test/fixtures/a/b/c": "/test/fixtures/a/b/c",
	})

	test_localPathToRemotePath(t, []string{"test/fixtures/a/b/c"}, "/e", true, map[string]string{
		"test/fixtures/a/b/c": "/e/test/fixtures/a/b/c",
	})

	test_localPathToRemotePath(t, []string{"test/fixtures/a/b/c"}, "/e/", true, map[string]string{
		"test/fixtures/a/b/c": "/e/test/fixtures/a/b/c",
	})

	test_localPathToRemotePath(t, []string{"test/fixtures/a/b/c", "test/fixtures/f"}, "/e/", true, map[string]string{
		"test/fixtures/a/b/c": "/e/test/fixtures/a/b/c",
		"test/fixtures/f":     "/e/test/fixtures/f",
	})
}

func test_localDirectoryToRemotePath(t *testing.T, locals []string, remote string, parentsPath bool, expects map[string]string) {
	for _, root := range locals {
		err := filepath.Walk(root, func(local string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			expectedRemote := expects[local]
			if actualRemote := localDirectoryToRemotePath(root, local, remote, parentsPath); expectedRemote != actualRemote {
				t.Errorf("%s should be %s", actualRemote, expectedRemote)
			}
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	}
}

func TestLocalDirectoryToRemotePath(t *testing.T) {
	test_localDirectoryToRemotePath(t, []string{"test/fixtures"}, "/", false, map[string]string{
		"test/fixtures/a/b/c": "/fixtures/a/b/c",
		"test/fixtures/d/e":   "/fixtures/d/e",
		"test/fixtures/f":     "/fixtures/f",
	})

	test_localDirectoryToRemotePath(t, []string{"test/fixtures"}, "/e", false, map[string]string{
		"test/fixtures/a/b/c": "/e/a/b/c",
		"test/fixtures/d/e":   "/e/d/e",
		"test/fixtures/f":     "/e/f",
	})

	test_localDirectoryToRemotePath(t, []string{"test/fixtures"}, "/e/", false, map[string]string{
		"test/fixtures/a/b/c": "/e/fixtures/a/b/c",
		"test/fixtures/d/e":   "/e/fixtures/d/e",
		"test/fixtures/f":     "/e/fixtures/f",
	})

	test_localDirectoryToRemotePath(t, []string{"test/fixtures/"}, "/e/", false, map[string]string{
		"test/fixtures/a/b/c": "/e/fixtures/a/b/c",
		"test/fixtures/d/e":   "/e/fixtures/d/e",
		"test/fixtures/f":     "/e/fixtures/f",
	})

	test_localDirectoryToRemotePath(t, []string{"test/fixtures"}, "/", true, map[string]string{
		"test/fixtures/a/b/c": "/test/fixtures/a/b/c",
		"test/fixtures/d/e":   "/test/fixtures/d/e",
		"test/fixtures/f":     "/test/fixtures/f",
	})

	test_localDirectoryToRemotePath(t, []string{"test/fixtures"}, "/e", true, map[string]string{
		"test/fixtures/a/b/c": "/e/test/fixtures/a/b/c",
		"test/fixtures/d/e":   "/e/test/fixtures/d/e",
		"test/fixtures/f":     "/e/test/fixtures/f",
	})

	test_localDirectoryToRemotePath(t, []string{"test/fixtures"}, "/e/", true, map[string]string{
		"test/fixtures/a/b/c": "/e/test/fixtures/a/b/c",
		"test/fixtures/d/e":   "/e/test/fixtures/d/e",
		"test/fixtures/f":     "/e/test/fixtures/f",
	})

	test_localDirectoryToRemotePath(t, []string{"test/fixtures/"}, "/e/", true, map[string]string{
		"test/fixtures/a/b/c": "/e/test/fixtures/a/b/c",
		"test/fixtures/d/e":   "/e/test/fixtures/d/e",
		"test/fixtures/f":     "/e/test/fixtures/f",
	})
}
