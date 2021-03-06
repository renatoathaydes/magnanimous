package mg

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// CopyAll copies all given files under the basePath by putting the files into the provided filesMap,
// then writing them out when the resulting ProcessedFiles are written.
func CopyAll(files *[]string, basePath string, filesMap WebFilesMap) error {
	for _, file := range *files {
		wf, err := Copy(file, basePath, true)
		if err != nil {
			return err
		}
		wf.SkipIfUpToDate = true // all static files should not be written when older than destination
		filesMap.WebFiles[file] = *wf
	}
	return nil
}

func AddNonWritables(files *[]string, basePath string, filesMap WebFilesMap) error {
	for _, file := range *files {
		wf, err := Copy(file, basePath, false)
		if err != nil {
			return err
		}
		filesMap.WebFiles[file] = *wf
	}
	return nil
}

func Copy(file, basePath string, writable bool) (*WebFile, error) {
	stats, err := os.Stat(file)
	if err != nil {
		return nil, &MagnanimousError{Code: IOError, message: err.Error()}
	}

	var proc = ProcessedFile{Path: file, LastUpdated: stats.ModTime()}
	proc.AppendContent(&copiedContent{file: file})
	return &WebFile{BasePath: basePath, Name: filepath.Base(file), Processed: &proc, NonWritable: !writable}, nil
}

type copiedContent struct {
	UnscopedContent
	file string
}

var _ Content = (*copiedContent)(nil)

func (c *copiedContent) GetLocation() *Location {
	loc := Location{Origin: c.file, Row: 0, Col: 0}
	return &loc
}

func (c *copiedContent) Write(writer io.Writer, context Context) ([]Content, error) {
	source, err := os.Open(c.file)
	if err != nil {
		return nil, err
	}
	defer source.Close()
	_, err = io.Copy(writer, bufio.NewReader(source))
	return nil, err
}
