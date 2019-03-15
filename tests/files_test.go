package tests

import (
	"github.com/renatoathaydes/magnanimous/mg"
	"testing"
)

var resolver mg.DefaultFileResolver

func init() {
	files := make(map[string]mg.WebFile)
	resolver = mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	// files for testing methods that need real files, not just paths
	files["source/abc/def/ghi/file.txt"] = mg.WebFile{Name: "file.txt"}
	files["source/abc/def/ghi/other.md"] = mg.WebFile{Name: "other.md"}
	files["source/abc/def/123.n"] = mg.WebFile{Name: "123.n"}
	files["source/abc/123.n"] = mg.WebFile{Name: "123.n"}
	files["source/123.n"] = mg.WebFile{Name: "123.n"}
	files["source/xxx/yyy.zzz"] = mg.WebFile{Name: "yyy.zzz"}
}

func ResolveFile(file, origin string) string {
	return resolver.Resolve(file, &mg.Location{Origin: origin})
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
	verifyEqual(2, t, ResolveFile("example.html", "source/hello.html"),
		"source/example.html")
	verifyEqual(3, t, ResolveFile("example.html", "hello.html"),
		"source/example.html")
}

func TestResolveAbsolutePath(t *testing.T) {
	verifyEqual(1, t, ResolveFile("/site/example.html", "source/processed/hello.html"),
		"source/site/example.html")
	verifyEqual(2, t, ResolveFile("/site/example.html", "source/hello.html"),
		"source/site/example.html")
	verifyEqual(3, t, ResolveFile("/site/example.html", "hello.html"),
		"source/site/example.html")
}

func TestResolveUpPath(t *testing.T) {
	verifyEqual(1, t, ResolveFile(".../other.md", "source/abc/def/ghi/file.txt"),
		"source/abc/def/ghi/other.md")
	verifyEqual(2, t, ResolveFile(".../123.n", "source/abc/def/ghi/file.txt"),
		"source/abc/def/123.n")
	verifyEqual(3, t, ResolveFile(".../123.n", "source/xxx/yyy.zzz"),
		"source/123.n")
	verifyEqual(4, t, ResolveFile(".../123.n", "source/123.n"),
		"source/123.n")
	verifyEqual(5, t, ResolveFile(".../yyy.zzz", "source/xxx/yyy.zzz"),
		"source/xxx/yyy.zzz")

	// cannot resolve non-existing file
	verifyEqual(6, t, ResolveFile(".../not_exists", "source/xxx/yyy.zzz"),
		"not_exists")

	// cannot find empty filename
	verifyEqual(6, t, ResolveFile(".../", "source/xxx/yyy.zzz"),
		"")
}
