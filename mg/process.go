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

func Process(files *[]string, basePath string, filesMap *WebFilesMap) {
	for _, file := range *files {
		wf := ProcessFile(file, basePath)
		(*filesMap)[file] = *wf
	}

	s := blackfriday.Run([]byte("# This is magnanimous!\n```go\nfunc main () { }\n```\n"))

	fmt.Println(string(s))
}

func ProcessFile(file, basePath string) *WebFile {
	f, err := os.Open(file)
	ExitIfError(&err, 5)
	reader := bufio.NewReader(f)
	s, err := f.Stat()
	ExitIfError(&err, 5)
	ctx, processed := ProcessReader(reader, file, int(s.Size()))
	return &WebFile{Context: ctx, BasePath: basePath, Processed: processed}
}

func ProcessReader(reader *bufio.Reader, file string, size int) (*WebFileContext, *ProcessedFile) {
	var builder strings.Builder
	builder.Grow(size)
	ctx := make(WebFileContext)
	processed := ProcessedFile{}
	previousWasOpenBracket := false
	previousWasCloseBracket := false
	parsingInstruction := false
	var row uint32 = 0
	var col uint32 = 0
	var instrFirstRow uint32
	var instrFirstCol uint32

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		ExitIfError(&err, 6)

		if r == '\n' {
			row++
			col = 1
			builder.WriteRune(r)
			continue
		}

		col++

		if r == '}' {
			if previousWasCloseBracket {
				if !parsingInstruction {
					log.Fatalf("Parsing Error at position %d:%d - Unexpected characters: }}", row, col)
				}
				if builder.Len() > 0 {
					processed.AppendContent(instruction(builder.String(), Location{Origin: file, Row: instrFirstRow, Col: instrFirstCol}))
				}
				builder.Reset()
				parsingInstruction = false
				previousWasCloseBracket = false
			} else {
				previousWasCloseBracket = true
			}
			continue
		} else {
			if previousWasCloseBracket {
				builder.WriteRune('}')
			}
			previousWasCloseBracket = false
		}

		if r == '{' {
			if previousWasOpenBracket {
				if builder.Len() > 0 {
					processed.AppendContent(&StringContent{Text: builder.String()})
					builder.Reset()
				}
				instrFirstRow = row + 1
				instrFirstCol = col + 1
				parsingInstruction = true
				previousWasOpenBracket = false
			} else {
				previousWasOpenBracket = true
			}
			continue
		} else {
			if previousWasOpenBracket {
				builder.WriteRune('{')
			}
			previousWasOpenBracket = false
		}

		builder.WriteRune(r)
	}

	if parsingInstruction {
		log.Fatalf("Parsing Error at position %d:%d - instruction was not properly closed with: }}", row, col)
	} else if builder.Len() > 0 {
		processed.AppendContent(&StringContent{Text: builder.String()})
	}

	return &ctx, &processed
}

func instruction(text string, location Location) Content {
	parts := strings.SplitN(strings.TrimSpace(text), " ", 2)
	switch len(parts) {
	case 0:
		fallthrough
	case 1:
		return &StringContent{Text: fmt.Sprintf("{{ %s }}", text)}
	}
	return validateInstruction(parts[0], parts[1], location)
}

func validateInstruction(name, arg string, location Location) Content {
	switch name {
	case "include":
		path := ResolveFile(arg, "source", location.Origin)
		return &IncludeInstruction{Name: name, Path: path, Origin: location}
	}

	log.Printf("WARNING: (%s) Instruction not implemented yet: %s", location.String(), name)
	return &StringContent{Text: fmt.Sprintf("{{ %s %s }}", name, arg)}
}

func WriteTo(dir string, filesMap WebFilesMap) {
	err := os.MkdirAll(dir, 0770)
	ExitIfError(&err, 9)
	for file, wf := range filesMap {
		targetPath, err := filepath.Rel(wf.BasePath, file)
		if err != nil {
			log.Printf("Unable to relativize path %s", file)
			targetPath = file
		}
		targetFile := filepath.Join(dir, targetPath)
		writeFile(file, targetFile, wf, filesMap)
	}
}

func writeFile(file, targetFile string, wf WebFile, filesMap WebFilesMap) {
	log.Printf("Creating file %s from %s", targetFile, file)
	err := os.MkdirAll(filepath.Dir(targetFile), 0770)
	ExitIfError(&err, 10)
	f, err := os.Create(targetFile)
	ExitIfError(&err, 10)
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, c := range wf.Processed.Contents {
		c.Write(w, filesMap)
	}
}

func ExitIfError(err *error, code int) {
	if *err != nil {
		log.Fatal(*err)
		os.Exit(code)
	}
}

func (c *StringContent) Write(writer io.Writer, files WebFilesMap) {
	_, err := writer.Write([]byte(c.Text))
	ExitIfError(&err, 11)
}

func (c *IncludeInstruction) Write(writer io.Writer, files WebFilesMap) {
	webFile, ok := files[c.Path]
	if !ok {
		log.Printf("WARNING: (%s) include non-existent resource: %s", c.Origin.String(), c.Path)
		_, err := writer.Write([]byte(fmt.Sprintf("{{ %s %s }}", c.Name, c.Path)))
		ExitIfError(&err, 11)
	} else {
		webFile.Write(writer, files)
	}
}

func (wf *WebFile) Write(writer io.Writer, files WebFilesMap) {
	for _, c := range wf.Processed.Contents {
		c.Write(writer, files)
	}
}
