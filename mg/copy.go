package mg

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

func CopyAll(files *[]string, basePath string, filesMap WebFilesMap) {
	for _, file := range *files {
		wf := Copy(file, basePath, true)
		filesMap.WebFiles[file] = *wf
	}
}

func AddNonWritables(files *[]string, basePath string, filesMap WebFilesMap) {
	for _, file := range *files {
		wf := Copy(file, basePath, false)
		filesMap.WebFiles[file] = *wf
	}
}

func Copy(file, basePath string, writable bool) *WebFile {
	var proc = ProcessedFile{}
	proc.AppendContent(&copiedContent{file: file})
	return &WebFile{BasePath: basePath, Name: filepath.Base(file), Processed: &proc, NonWritable: !writable}
}

type copiedContent struct {
	file string
}

func (c *copiedContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []ContextStackItem) error {
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
