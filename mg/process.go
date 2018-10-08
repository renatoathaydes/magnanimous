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

type strContent struct {
	text string
}

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
	ctx, processed := ProcessReader(reader, int(s.Size()))
	return &WebFile{Context: ctx, BasePath: basePath, Processed: processed}
}

func ProcessReader(reader *bufio.Reader, size int) (*WebFileContext, *ProcessedFile) {
	var builder strings.Builder
	builder.Grow(size)
	ctx := make(WebFileContext)
	processed := ProcessedFile{}
	previousWasOpenBracket := false
	previousWasCloseBracket := false
	parsingInstruction := false
	row := 1
	col := 0

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
					processed.appendContent(instr(builder.String(), row, col))
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
					processed.appendContent(&strContent{text: builder.String()})
					builder.Reset()
				}
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
		processed.appendContent(&strContent{text: builder.String()})
	}

	log.Printf("Parsed contents: %s", processed.Contents)

	return &ctx, &processed
}

func instr(text string, row, col int) Content {
	parts := strings.Fields(text)
	switch len(parts) {
	case 0:
		return &strContent{text: ""}
	case 1:
		log.Fatalf("Instruction Error at position %d:%d - instruction missing argument: %s", row, col, parts[0])
	default:
		return parseInstr(parts[0], parts[1:], row, col)
	}
	panic("Unreachable")
}

func parseInstr(name string, args []string, row, col int) Content {
	switch name {
	case "include":
		if len(args) == 1 {
			return &strContent{text: "example"}
		} else {
			log.Fatalf("Instruction Error at position %d:%d - wrong number of arguments for "+
				"include instruction, expected 1, got %d", row, col, len(args))
		}
	default:
		log.Fatalf("Unknown instruction: %s", name)
	}
	panic("Unreachable")
}

func WriteTo(dir string, filesMap *WebFilesMap) {
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
		defer w.Flush()
		for _, c := range wf.Processed.Contents {
			c.Write(w)
		}
	}
}

func ExitIfError(err *error, code int) {
	if *err != nil {
		log.Fatal(*err)
		os.Exit(code)
	}
}

func (c *strContent) Write(writer io.Writer) {
	_, err := writer.Write([]byte(c.text))
	ExitIfError(&err, 11)
}
