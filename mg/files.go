package mg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type resolver struct {
	basePath string
}

var DefaultFileResolver FileResolver = resolver{basePath: "source"}

func (r resolver) FilesIn(dir string, from Location) (dirPath string, fi []os.FileInfo, e error) {
	dirPath = r.Resolve(dir, from)
	fi, e = ioutil.ReadDir(dirPath)
	return
}

func (r resolver) Resolve(path string, from Location) string {
	if strings.HasPrefix(path, "/") {
		// absolute path
		return filepath.Join(r.basePath, path)
	}

	// relative path
	p := filepath.Join(filepath.Dir(from.Origin), path)

	// must not go up higher than basePath
	for strings.HasPrefix(p, "../") {
		p = p[3:]
	}
	if !strings.HasPrefix(p, r.basePath) {
		return filepath.Join(r.basePath, p)
	}
	return p
}

func isMd(file string) bool {
	return strings.ToLower(filepath.Ext(file)) == ".md"
}
