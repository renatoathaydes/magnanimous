package mg

import (
	"bufio"
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"log"
	"path/filepath"
	"strings"

	"os"
)

type StrContent struct {
	text string
}

func Process(files *[]string, basePath string, filesMap *WebFilesMap) {
	for _, file := range *files {
		wf := ParseFile(file, basePath)
		(*filesMap)[file] = *wf
	}

	s := blackfriday.Run([]byte("# This is magnanimous!\n```go\nfunc main () { }\n```\n"))

	fmt.Println(string(s))
}

func ParseFile(file, basePath string) *WebFile {
	f, err := os.Open(file)
	ExitIfError(&err, 5)
	reader := bufio.NewReader(f)
	s, err := f.Stat()
	ExitIfError(&err, 5)
	ctx, processed := ParseReader(reader, int(s.Size()))
	return &WebFile{Context: ctx, BasePath: basePath, Processed: processed}
}

func ParseReader(reader *bufio.Reader, size int) (*WebFileContext, *ProcessedFile) {
	var builder strings.Builder
	builder.Grow(size)
	ctx := make(WebFileContext)
	processed := ProcessedFile{}
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		ExitIfError(&err, 6)
		builder.WriteRune(r)
	}
	processed.appendContent(&StrContent{text: builder.String()})
	return &ctx, &processed
}

func WriteAt(dir string, filesMap *WebFilesMap) {
	err := os.MkdirAll(dir, 0770)
	ExitIfError(&err, 9)
	for file, wf := range *filesMap {
		targetPath, err := filepath.Rel(wf.BasePath, file)
		if err != nil {
			log.Printf("Unable to relativize path %s", file)
			targetPath = file
		}
		targetFile := filepath.Join(dir, targetPath)
		log.Printf("Creating file %s from %s", targetFile, file)
		err = os.MkdirAll(filepath.Dir(targetFile), 0770)
		ExitIfError(&err, 10)
		f, err := os.Create(targetFile)
		ExitIfError(&err, 10)
		defer f.Close()
		w := bufio.NewWriter(f)
		for _, c := range wf.Processed.Contents {
			c.Write(w)
		}
		w.Flush()
	}
}

func ExitIfError(err *error, code int) {
	if *err != nil {
		log.Fatal(*err)
		os.Exit(code)
	}
}

func (c *StrContent) Write(writer io.Writer) {
	_, err := writer.Write([]byte(c.text))
	ExitIfError(&err, 11)
}
