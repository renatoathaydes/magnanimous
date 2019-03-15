package mg

import (
	"os"
	"path/filepath"
	"strings"
)

// DefaultFileResolver is the default implementation of [FileResolver].
//
// It can resolve absolute and relative paths given its BasePath and the Files map.
//
// It can also resolve up-paths, i.e. paths starting with '.../', which resolve to a existing file in
// the current directory or a parent directory, up until the file is found, or the BasePath is reached.
type DefaultFileResolver struct {
	BasePath string
	Files    *WebFilesMap
}

var _ FileResolver = (*DefaultFileResolver)(nil)

type filesCollector func() ([]string, error)

func collectFiles(sourcesDir, processedDir, staticDir string) (procFiles, staticFiles, otherFiles []string) {
	async := func(fc filesCollector, c chan []string) {
		s, err := fc()
		if err != nil {
			panic(err)
		}
		c <- s
	}

	procC, statC, othersC := make(chan []string), make(chan []string), make(chan []string)
	go async(func() ([]string, error) { return getFilesAt(processedDir) }, procC)
	go async(func() ([]string, error) { return getFilesAt(staticDir) }, statC)
	go async(func() ([]string, error) {
		return getFilesAt(sourcesDir, processedDir, staticDir)
	}, othersC)

	procFiles, staticFiles, otherFiles = <-procC, <-statC, <-othersC
	return
}

func getFilesAt(root string, exclusions ...string) ([]string, error) {
	var files []string
	notExcluded := func(path string) bool {
		for _, e := range exclusions {
			if strings.HasPrefix(path, e) {
				return false
			}
		}
		return true
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && notExcluded(path) {
			files = append(files, path)
		}
		return err
	})
	if os.IsNotExist(err) {
		return []string{}, nil
	}
	return files, err
}

func (r *DefaultFileResolver) FilesIn(dir string, from *Location) (dirPath string, webFiles []WebFile, e error) {
	dirPath = r.Resolve(dir, from)
	for path, wf := range r.Files.WebFiles {
		if !wf.NonWritable && filepath.Dir(path) == dirPath {
			webFiles = append(webFiles, wf)
		}
	}
	return
}

func (r *DefaultFileResolver) Resolve(path string, from *Location) string {
	if strings.HasPrefix(path, "/") {
		// absolute path
		return filepath.Join(r.BasePath, path)
	}

	if strings.HasPrefix(path, ".../") {
		return r.searchUp(path[4:], from)
	}

	// relative path
	p := filepath.Join(filepath.Dir(from.Origin), path)

	// must not go up higher than basePath
	for strings.HasPrefix(p, "../") {
		p = p[3:]
	}
	if !strings.HasPrefix(p, r.BasePath) {
		return filepath.Join(r.BasePath, p)
	}
	return p
}

func (r *DefaultFileResolver) searchUp(path string, from *Location) string {
	dir := filepath.Dir(from.Origin)
	for dir != "." {
		name := filepath.Join(dir, path)
		if _, ok := r.Files.WebFiles[name]; ok {
			return name
		}
		dir = filepath.Dir(dir)
	}
	// can't find it
	return path
}

func isMd(file string) bool {
	return strings.ToLower(filepath.Ext(file)) == ".md"
}

func changeFileExt(path, extension string) string {
	ext := filepath.Ext(path)
	return path[0:len(path)-len(ext)] + "." + extension
}
