package mg

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

func CopyAll(files *[]string, basePath string, filesMap WebFilesMap) error {
	for _, file := range *files {
		wf, err := Copy(file, basePath, true)
		if err != nil {
			return err
		}
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
	f, err := os.Open(file)
	if err != nil {
		return nil, &MagnanimousError{Code: IOError, message: err.Error()}
	}
	defer f.Close()
	stats, err := f.Stat()
	if err != nil {
		return nil, &MagnanimousError{Code: IOError, message: err.Error()}
	}

	var proc = ProcessedFile{Path: file, LastUpdated: stats.ModTime()}
	proc.AppendContent(&copiedContent{file: file})
	return &WebFile{BasePath: basePath, Name: filepath.Base(file), Processed: &proc, NonWritable: !writable}, nil
}

type copiedContent struct {
	file string
}

func (c *copiedContent) Write(writer io.Writer, stack ContextStack) error {
	f, err := os.Open(c.file)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	defer f.Close()
	_, err = io.Copy(writer, bufio.NewReader(f))
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}
