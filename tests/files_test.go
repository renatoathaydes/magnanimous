package tests

import (
	"github.com/renatoathaydes/magnanimous/mg"
	"testing"
)

func TestResolveFile(t *testing.T) {
	verifyEqual(1, t, mg.ResolveFile("a", "target", "target"), "target/a")
	verifyEqual(2, t, mg.ResolveFile("/a", "target", "target"), "target/a")
	verifyEqual(3, t, mg.ResolveFile("/a", "target/", "target"), "target/a")
	verifyEqual(4, t, mg.ResolveFile("/a", "target/", "target/other"), "target/a")
	verifyEqual(5, t, mg.ResolveFile("a", "target/", "target/other"), "target/a")
	verifyEqual(6, t, mg.ResolveFile("a", "target/", "target/abc/file.html"), "target/abc/a")
	verifyEqual(7, t, mg.ResolveFile("../a", "target/", "target/other"), "target/a")
	verifyEqual(8, t, mg.ResolveFile("../../a", "target/", "target/other"), "target/a")
	verifyEqual(9, t, mg.ResolveFile("../../../a", "target/", "target/other"), "target/a")
}

func TestResolveRelativePath(t *testing.T) {
	verifyEqual(1, t, mg.ResolveFile("example.html", "source", "source/processed/hello.html"),
		"source/processed/example.html")
}

func TestResolveAbsolutePath(t *testing.T) {
	verifyEqual(1, t, mg.ResolveFile("/site/example.html", "source", "source/processed/hello.html"),
		"source/site/example.html")
}

func verifyEqual(i uint16, t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("[%d] Expected '%s' but was '%s'", i, expected, actual)
	}
}
