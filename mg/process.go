package mg

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type parserState struct {
	file    string
	reader  *bufio.Reader
	row     uint32
	col     uint32
	pf      *ProcessedFile
	builder *strings.Builder
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
	processed, magErr := ProcessReader(reader, file, int(s.Size()))
	if magErr != nil {
		return nil, magErr
	}
	nonWritable := strings.HasPrefix(filepath.Base(file), "_")
	return &WebFile{BasePath: basePath, Processed: &processed, NonWritable: nonWritable}, nil
}

func ProcessReader(reader *bufio.Reader, file string, size int) (ProcessedFile, *MagnanimousError) {
	var builder strings.Builder
	builder.Grow(size)
	processed := ProcessedFile{}
	state := parserState{file: file, row: 1, col: 1, builder: &builder, reader: reader, pf: &processed}
	magErr := parseText(&state, isMd(file))
	if magErr != nil {
		return processed, magErr
	}
	if isMd(file) {
		processed = MarkdownToHtml(processed)
	}
	return processed, nil
}

func parseText(state *parserState, isMarkDown bool) *MagnanimousError {
	previousWasOpenBracket := false
	previousWasEscape := false
	reader := state.reader
	builder := state.builder

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &MagnanimousError{message: err.Error(), Code: IOError}
		}

		if r == '\n' {
			state.row++
			state.col = 1
			builder.WriteRune(r)
			continue
		}

		state.col++

		if r == '{' {
			if previousWasEscape {
				builder.WriteRune(r)
				previousWasEscape = false
				continue
			}
			if previousWasOpenBracket {
				if builder.Len() > 0 {
					state.pf.AppendContent(&StringContent{Text: builder.String(), MarkDown: isMarkDown})
					builder.Reset()
				}
				previousWasOpenBracket = false
				magErr := parseInstruction(state, isMarkDown)
				if magErr != nil {
					return magErr
				}
			} else {
				previousWasOpenBracket = true
			}
			continue
		}

		if previousWasOpenBracket {
			builder.WriteRune('{')
		}
		previousWasOpenBracket = false

		if r == '\\' {
			previousWasEscape = true
		} else {
			builder.WriteRune(r)
		}
	}

	// append any pending content to the builder
	if builder.Len() > 0 || previousWasOpenBracket || previousWasEscape {
		if previousWasOpenBracket {
			builder.WriteRune('{')
		} else if previousWasEscape {
			builder.WriteRune('\\')
		}
		state.pf.AppendContent(&StringContent{Text: builder.String(), MarkDown: isMarkDown})
		builder.Reset()
	}

	return nil
}

func parseInstruction(state *parserState, isMarkDown bool) *MagnanimousError {
	previousWasCloseBracket := false
	previousWasEscape := false
	instrFirstRow := state.row
	instrFirstCol := state.col - 2
	reader := state.reader
	builder := state.builder

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &MagnanimousError{message: err.Error(), Code: IOError}
		}

		if r == '\n' {
			state.row++
			state.col = 1
			builder.WriteRune(r)
			continue
		}

		state.col++

		if r == '}' {
			if previousWasEscape {
				builder.WriteRune(r)
				previousWasEscape = false
				continue
			}
			if previousWasCloseBracket {
				if builder.Len() > 0 {
					appendContent(state.pf, builder.String(), isMarkDown,
						Location{Origin: state.file, Row: instrFirstRow, Col: instrFirstCol})
				}
				builder.Reset()
				return nil
			} else {
				previousWasCloseBracket = true
			}
			continue
		}

		if previousWasCloseBracket {
			builder.WriteRune('}')
		}
		previousWasCloseBracket = false

		if r == '\\' {
			previousWasEscape = true
		} else {
			builder.WriteRune(r)
		}
	}

	return NewError(Location{Origin: state.file, Row: state.row, Col: state.col}, ParseError,
		fmt.Sprintf("instruction started at (%d:%d) was not properly closed with '}}'",
			instrFirstRow, instrFirstCol))
}

func appendContent(pf *ProcessedFile, text string, isMarkDown bool, location Location) {
	parts := strings.SplitN(strings.TrimSpace(text), " ", 2)
	switch len(parts) {
	case 0:
		// nothing to do
	case 1:
		if parts[0] == "end" {
			err := pf.EndScope()
			if err != nil {
				log.Printf("WARNING: (%s) %s", location.String(), err.Error())
				pf.AppendContent(unevaluatedExpression(text))
			}
		} else {
			log.Printf("WARNING: (%s) Instruction missing argument: %s", location.String(), text)
			pf.AppendContent(unevaluatedExpression(text))
		}
	case 2:
		content := createInstruction(parts[0], parts[1], isMarkDown, pf.currentScope(), location, text)
		if content != nil {
			pf.AppendContent(content)
		}
	}
}

func createInstruction(name, arg string, isMarkDown bool, scope Scope,
	location Location, original string) Content {
	switch name {
	case "include":
		return NewIncludeInstruction(arg, location)
	case "define":
		return NewVariable(arg, location, original, scope)
	case "eval":
		return NewExpression(arg, location, isMarkDown, original)
	case "for":
		return NewForInstruction(arg, location, isMarkDown, original)
	}

	log.Printf("WARNING: (%s) Unknown instruction: %s", location.String(), name)
	return unevaluatedExpression(original)
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
		err := c.Write(w, filesMap, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *StringContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError {
	_, err := writer.Write([]byte(c.Text))
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func (c *StringContent) IsMarkDown() bool {
	return c.MarkDown
}

func (c *StringContent) String() string {
	return fmt.Sprintf("StringContent{%s}", c.Text)
}

func (wf *WebFile) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError {
	for _, c := range wf.Processed.Contents {
		err := c.Write(writer, files, inclusionChain)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *HtmlFromMarkdownContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError {
	content, magErr := readBytes(&f.MarkDownContent, files, inclusionChain)
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

func readBytes(c *Content, files WebFilesMap, inclusionChain []Location) ([]byte, *MagnanimousError) {
	var b bytes.Buffer
	b.Grow(1024)
	err := (*c).Write(&b, files, inclusionChain)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func inclusionChainToString(locations []Location) string {
	var b strings.Builder
	b.WriteRune('[')
	last := len(locations) - 1
	for i, loc := range locations {
		b.WriteString(loc.String())
		if i != last {
			b.WriteString(" -> ")
		}
	}
	b.WriteRune(']')
	return b.String()
}
