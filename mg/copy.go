package mg

import (
	"bufio"
	"io"
	"os"
)

func CopyAll(files *[]string, basePath string, filesMap *WebFilesMap) {
	for _, file := range *files {
		wf := Copy(file, basePath)
		(*filesMap)[file] = *wf
	}
}

func Copy(file, basePath string) *WebFile {
	var proc = ProcessedFile{}
	proc.appendContent(&copiedContent{file: file})
	return &WebFile{BasePath: basePath, Processed: &proc}
}

type copiedContent struct {
	file string
}

func (c *copiedContent) Write(writer io.Writer) {
	f, err := os.Open(c.file)
	ExitIfError(&err, 30)
	defer f.Close()
	io.Copy(writer, bufio.NewReader(f))
}
