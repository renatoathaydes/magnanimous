package mg

import (
	"path/filepath"
	"strings"
)

func ResolveFile(file, basePath, currentLocation string) string {
	if strings.HasPrefix(file, "/") {
		// absolute path
		return filepath.Join(basePath, file)
	}

	// relative path
	p := filepath.Join(currentLocation, file)

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
