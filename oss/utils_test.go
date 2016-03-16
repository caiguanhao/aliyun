package main

import "testing"

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
