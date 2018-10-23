package tests

import (
	"github.com/renatoathaydes/magnanimous/mg"
	"testing"
)

func ResolveFile(file, origin string) string {
	return mg.Resolve(file, "source", mg.Location{Origin: origin})
}

func TestResolveFile(t *testing.T) {
	verifyEqual(1, t, ResolveFile("a", "source"), "source/a")
	verifyEqual(2, t, ResolveFile("/a", "source"), "source/a")
	verifyEqual(3, t, ResolveFile("/a", "source"), "source/a")
	verifyEqual(4, t, ResolveFile("/a", "source/other"), "source/a")
	verifyEqual(5, t, ResolveFile("a", "source/other"), "source/a")
	verifyEqual(6, t, ResolveFile("a", "source/abc/file.html"), "source/abc/a")
	verifyEqual(7, t, ResolveFile("../a", "source/other"), "source/a")
	verifyEqual(8, t, ResolveFile("../../a", "source/other"), "source/a")
	verifyEqual(9, t, ResolveFile("../../../a", "source/other"), "source/a")
}

func TestResolveRelativePath(t *testing.T) {
	verifyEqual(1, t, ResolveFile("example.html", "source/processed/hello.html"),
		"source/processed/example.html")
}

func TestResolveAbsolutePath(t *testing.T) {
	verifyEqual(1, t, ResolveFile("/site/example.html", "source/processed/hello.html"),
		"source/site/example.html")
}
