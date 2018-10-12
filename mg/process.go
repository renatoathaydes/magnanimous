package mg

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"log"
	"path/filepath"
	"strings"

	"os"
)

func isMd(file string) bool {
	return strings.ToLower(filepath.Ext(file)) == ".md"
}

func Process(files *[]string, basePath string, filesMap WebFilesMap) {
	for _, file := range *files {
		wf, err := ProcessFile(file, basePath)
		if err != nil {
			panic(err)
		}
		filesMap[file] = *wf
	}
}

func ProcessFile(file, basePath string) (*WebFile, *MagnanimousError) {
	f, err := os.Open(file)
	if err != nil {
		return nil, &MagnanimousError{message: err.Error(), Code: IOError}
	}
	reader := bufio.NewReader(f)
	s, err := f.Stat()
	if err != nil {
		return nil, &MagnanimousError{message: err.Error(), Code: IOError}
	}
	ctx, processed, magErr := ProcessReader(reader, file, int(s.Size()))
	if magErr != nil {
		return nil, magErr
	}
	nonWritable := strings.HasPrefix(filepath.Base(file), "_")
	return &WebFile{Context: ctx, BasePath: basePath, Processed: processed, NonWritable: nonWritable}, nil
}

func ProcessReader(reader *bufio.Reader, file string, size int) (WebFileContext, ProcessedFile, *MagnanimousError) {
	isMarkDown := isMd(file)
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
		if err != nil {
			return ctx, processed, &MagnanimousError{message: err.Error(), Code: IOError}
		}

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
					return ctx, processed, NewParseError(
						Location{Origin: file, Row: row, Col: col},
						"Unexpected characters: '}}'")
				}
				if builder.Len() > 0 {
					processed.AppendContent(instruction(builder.String(),
						isMarkDown,
						Location{Origin: file, Row: instrFirstRow, Col: instrFirstCol}))
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
					processed.AppendContent(&StringContent{Text: builder.String(), MarkDown: isMarkDown})
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
		return ctx, processed, NewParseError(
			Location{Origin: file, Row: row, Col: col},
			"instruction was not properly closed with: }}")
	} else if builder.Len() > 0 {
		processed.AppendContent(&StringContent{Text: builder.String(), MarkDown: isMarkDown})
	}
	if isMd(file) {
		processed = MarkdownToHtml(processed)
	}
	return ctx, processed, nil
}

func instruction(text string, isMarkDown bool, location Location) Content {
	parts := strings.SplitN(strings.TrimSpace(text), " ", 2)
	switch len(parts) {
	case 0:
		fallthrough
	case 1:
		return &StringContent{Text: fmt.Sprintf("{{ %s }}", text), MarkDown: isMarkDown}
	}
	return createInstruction(parts[0], parts[1], location)
}

func createInstruction(name, arg string, location Location) Content {
	switch name {
	case "include":
		path := ResolveFile(arg, "source", location.Origin)
		return &IncludeInstruction{Name: name, Path: path, Origin: location}
	}

	log.Printf("WARNING: (%s) Instruction not implemented yet: %s", location.String(), name)
	return &StringContent{Text: fmt.Sprintf("{{ %s %s }}", name, arg)}
}

func MarkdownToHtml(file ProcessedFile) ProcessedFile {
	convertedContent := make([]Content, 0, len(file.Contents))
	for _, c := range file.Contents {
		if c.IsMarkDown() {
			convertedContent = append(convertedContent, &HtmlFromMarkdownContent{MarkDownContent: c})
		} else {
			convertedContent = append(convertedContent, c)
		}
	}
	return ProcessedFile{Contents: convertedContent, NewExtension: ".html"}
}

func WriteTo(dir string, filesMap WebFilesMap) *MagnanimousError {
	err := os.MkdirAll(dir, 0770)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	for file, wf := range filesMap {
		if wf.NonWritable {
			continue
		}
		targetPath, err := filepath.Rel(wf.BasePath, file)
		if err != nil {
			log.Printf("Unable to relativize path %s", file)
			targetPath = file
		}
		targetFile := filepath.Join(dir, targetPath)
		if wf.Processed.NewExtension != "" {
			ext := filepath.Ext(targetFile)
			targetFile = targetFile[0:len(targetFile)-len(ext)] + wf.Processed.NewExtension
		}
		magErr := writeFile(file, targetFile, wf, filesMap)
		if magErr != nil {
			return magErr
		}
	}
	return nil
}

func writeFile(file, targetFile string, wf WebFile, filesMap WebFilesMap) *MagnanimousError {
	log.Printf("Creating file %s from %s", targetFile, file)
	err := os.MkdirAll(filepath.Dir(targetFile), 0770)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	f, err := os.Create(targetFile)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, c := range wf.Processed.Contents {
		err := c.Write(w, filesMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *StringContent) Write(writer io.Writer, files WebFilesMap) *MagnanimousError {
	_, err := writer.Write([]byte(c.Text))
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func (c *IncludeInstruction) Write(writer io.Writer, files WebFilesMap) *MagnanimousError {
	webFile, ok := files[c.Path]
	if !ok {
		log.Printf("WARNING: (%s) include non-existent resource: %s", c.Origin.String(), c.Path)
		_, err := writer.Write([]byte(fmt.Sprintf("{{ %s %s }}", c.Name, c.Path)))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		err := webFile.Write(writer, files)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *StringContent) IsMarkDown() bool {
	return c.MarkDown
}

func (c *IncludeInstruction) IsMarkDown() bool {
	return c.MarkDown
}

func (wf *WebFile) Write(writer io.Writer, files WebFilesMap) *MagnanimousError {
	for _, c := range wf.Processed.Contents {
		err := c.Write(writer, files)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *HtmlFromMarkdownContent) Write(writer io.Writer, files WebFilesMap) *MagnanimousError {
	content, magErr := readBytes(&f.MarkDownContent, files)
	if magErr != nil {
		return magErr
	}
	_, err := writer.Write(blackfriday.Run(content))
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func (_ *HtmlFromMarkdownContent) IsMarkDown() bool {
	return false
}

func readBytes(c *Content, files WebFilesMap) ([]byte, *MagnanimousError) {
	var b bytes.Buffer
	b.Grow(1024)
	err := (*c).Write(&b, files)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
