package tests

import (
	"os"
	"testing"
)

func TestProj0(t *testing.T) {
	dir := runMg(t, "test_proj_0")
	defer os.RemoveAll(dir)

	files, err := readAll(dir)
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

	files, err := readAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("Expected 2 output files, but got: %v", files)
	}

	assertFileContents(t, files, dir, "a.txt", "")
	assertFileContents(t, files, dir, "main.txt", "A and B:\n\n10\n20")
}

func TestProj2(t *testing.T) {
	dir := runMg(t, "test_proj_2")
	defer os.RemoveAll(dir)

	files, err := readAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 4 {
		t.Fatalf("Expected 4 output files, but got: %v", files)
	}

	assertFileContents(t, files, dir, "sections/section_a.txt", "This is section A.")
	assertFileContents(t, files, dir, "sections/section_b.txt", "This is section B.")
	assertFileContents(t, files, dir, "main.txt", "Main.\n\nSections:\n\n"+
		"  * Section A\n\n"+
		"  * Section B\n\n"+
		"  * Section C\n\n"+
		"End.")
}
