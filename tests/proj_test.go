package tests

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestProj0(t *testing.T) {
	dir := runMg(t, "test_proj_0")
	defer os.RemoveAll(dir)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("Expected no output files, but got: %v", files)
	}
}

func TestProj1(t *testing.T) {
	dir := runMg(t, "test_proj_1")
	defer os.RemoveAll(dir)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Fatalf("Expected 3 output files, but got: %v", files)
	}

	assertFileContents(t, files, "a.txt", "")
	assertFileContents(t, files, "b.txt", "")
	assertFileContents(t, files, "main.txt", "10\n20")
}
