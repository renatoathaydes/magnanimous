package mg

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/russross/blackfriday"
	"io"
	"log"
	"path/filepath"
	"strings"

	"os"
)

type parserState struct {
	file    string
	reader  *bufio.Reader
	row     uint32
	col     uint32
	ctx     *WebFileContext
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
	ctx, processed, magErr := ProcessReader(reader, file, int(s.Size()))
	if magErr != nil {
		return nil, magErr
	}
	nonWritable := strings.HasPrefix(filepath.Base(file), "_")
	return &WebFile{Context: ctx, BasePath: basePath, Processed: processed, NonWritable: nonWritable}, nil
}

func ProcessReader(reader *bufio.Reader, file string, size int) (WebFileContext, ProcessedFile, *MagnanimousError) {
	var builder strings.Builder
	builder.Grow(size)
	ctx := make(WebFileContext)
	processed := ProcessedFile{}
	state := parserState{file: file, row: 1, col: 1, builder: &builder, reader: reader, ctx: &ctx, pf: &processed}
	magErr := parseText(&state, isMd(file))
	if magErr != nil {
		return ctx, processed, magErr
	}
	if isMd(file) {
		processed = MarkdownToHtml(processed)
	}
	return ctx, processed, nil
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
					content := instructionContent(builder.String(),
						isMarkDown,
						state.ctx,
						Location{Origin: state.file, Row: instrFirstRow, Col: instrFirstCol})
					if content != nil {
						state.pf.AppendContent(content)
					}
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

	return NewParseError(Location{Origin: state.file, Row: state.row, Col: state.col},
		fmt.Sprintf("instruction started at (%d:%d) was not properly closed with '}}'",
			instrFirstRow, instrFirstCol))
}

func instructionContent(text string, isMarkDown bool, ctx *WebFileContext, location Location) Content {
	parts := strings.SplitN(strings.TrimSpace(text), " ", 2)
	switch len(parts) {
	case 0:
		fallthrough
	case 1:
		return &StringContent{Text: fmt.Sprintf("{{%s}}", text), MarkDown: isMarkDown}
	}
	return createInstruction(parts[0], parts[1], isMarkDown, ctx, location, text)
}

func createInstruction(name, arg string, isMarkDown bool, ctx *WebFileContext,
	location Location, original string) Content {
	switch name {
	case "include":
		path := ResolveFile(arg, "source", location.Origin)
		return &IncludeInstruction{Name: name, Path: path, Origin: location}
	case "define":
		parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
		if len(parts) == 2 {
			variable, rawExpr := parts[0], parts[1]
			expr, err := govaluate.NewEvaluableExpression(rawExpr)
			if err != nil {
				log.Printf("WARNING: (%s) Unable to eval (defining %s): %s (%s)",
					location.String(), variable, rawExpr, err.Error())
				goto returnUnevaluated
			}
			v, err := expr.Evaluate(*ctx)
			if err != nil {
				log.Printf("WARNING: (%s) eval failure: %s", location.String(), err.Error())
				goto returnUnevaluated
			}
			(*ctx)[variable] = v
			return nil
		}
		log.Printf("WARNING: (%s) malformed define expression: %s", location.String(), arg)
		goto returnUnevaluated
	case "eval":
		expr, err := govaluate.NewEvaluableExpression(arg)
		if err != nil {
			log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
			goto returnUnevaluated
		}
		r, err := expr.Evaluate(*ctx)
		if err != nil {
			log.Printf("WARNING: (%s) eval failure: %s", location.String(), err.Error())
			goto returnUnevaluated
		}
		return &StringContent{Text: fmt.Sprintf("%v", r), MarkDown: isMarkDown}
	}

	log.Printf("WARNING: (%s) Instruction not implemented yet: %s", location.String(), name)

returnUnevaluated:
	// do not set MarkDown flag as unevaluated content cannot be markdown
	return &StringContent{Text: fmt.Sprintf("{{%s}}", original)}
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
