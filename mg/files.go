package mg

import (
	"os"
	"path/filepath"
	"strings"
)

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

func isMd(file string) bool {
	return strings.ToLower(filepath.Ext(file)) == ".md"
}

func changeFileExt(path, extension string) string {
	ext := filepath.Ext(path)
	return path[0:len(path)-len(ext)] + "." + extension
}
