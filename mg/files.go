package mg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type resolver struct{}

var DefaultFileResolver FileResolver = resolver{}

func (r resolver) FilesIn(dir string, from Location) (dirPath string, fi []os.FileInfo, e error) {
	dirPath = r.Resolve(dir, from)
	fi, e = ioutil.ReadDir(dirPath)
	return
}

func (r resolver) Resolve(path string, from Location) string {
	basePath := "source"
	if strings.HasPrefix(path, "/") {
		// absolute path
		return filepath.Join(basePath, path)
	}

	// relative path
	p := filepath.Join(filepath.Dir(from.Origin), path)

	// must not go up higher than basePath
	for strings.HasPrefix(p, "../") {
		p = p[3:]
	}
	if !strings.HasPrefix(p, basePath) {
		return filepath.Join(basePath, p)
	}
	return p
}

func isMd(file string) bool {
	return strings.ToLower(filepath.Ext(file)) == ".md"
}
