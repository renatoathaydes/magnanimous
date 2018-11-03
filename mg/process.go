package mg

import (
	"bufio"
	"fmt"
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

func (mag *Magnanimous) ReadAll() (WebFilesMap, error) {
	processedDir := filepath.Join(mag.SourcesDir, "processed")
	staticDir := filepath.Join(mag.SourcesDir, "static")

	procFiles, staticFiles, otherFiles := collectFiles(mag.SourcesDir, processedDir, staticDir)
	webFiles := make(WebFilesMap, len(procFiles)+len(staticFiles)+len(otherFiles))

	err := ProcessAll(procFiles, processedDir, mag.SourcesDir, webFiles)
	if err != nil {
		return nil, err
	}
	CopyAll(&staticFiles, staticDir, webFiles)
	AddNonWritables(&otherFiles, mag.SourcesDir, webFiles)

	return webFiles, nil
}

func ProcessAll(files []string, basePath, sourcesDir string, webFiles WebFilesMap) error {
	resolver := DefaultFileResolver{BasePath: sourcesDir, Files: webFiles}
	for _, file := range files {
		wf, err := ProcessFile(file, basePath, &resolver)
		if err != nil {
			return err
		}
		webFiles[file] = *wf
	}
	if globalCtx, ok := webFiles[filepath.Join(basePath, "_global_context")]; ok {
		globalCtx.runSideEffects(webFiles, nil)
		var globalContext RootScope = globalCtx.Processed.Context()
		if len(globalContext) > 0 {
			for _, webFile := range webFiles {
				webFile.Processed.setParent(globalContext)
			}
		}
	}
	return nil
}

func ProcessFile(file, basePath string, resolver FileResolver) (*WebFile, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, &MagnanimousError{message: err.Error(), Code: IOError}
	}
	reader := bufio.NewReader(f)
	s, err := f.Stat()
	if err != nil {
		return nil, &MagnanimousError{message: err.Error(), Code: IOError}
	}
	processed, magErr := ProcessReader(reader, file, int(s.Size()), resolver)
	if magErr != nil {
		return nil, magErr
	}
	processed.scopeStack = nil // the stack is no longer required
	nonWritable := strings.HasPrefix(filepath.Base(file), "_")
	return &WebFile{BasePath: basePath, Name: filepath.Base(file), Processed: processed, NonWritable: nonWritable}, nil
}

func ProcessReader(reader *bufio.Reader, file string, sizeHint int, resolver FileResolver) (*ProcessedFile, error) {
	var builder strings.Builder
	builder.Grow(sizeHint)
	isMarkDown := isMd(file)
	processed := ProcessedFile{context: make(map[string]interface{}, 4)}
	state := parserState{file: file, row: 1, col: 1, builder: &builder, reader: reader, pf: &processed}
	magErr := parseText(&state, resolver)
	if magErr != nil {
		return &processed, magErr
	}
	if isMarkDown {
		processed = MarkdownToHtml(processed)
	}
	return &processed, nil
}

func parseText(state *parserState, resolver FileResolver) error {
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
					state.pf.AppendContent(&StringContent{Text: builder.String()})
					builder.Reset()
				}
				previousWasOpenBracket = false
				magErr := parseInstruction(state, resolver)
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
		state.pf.AppendContent(&StringContent{Text: builder.String()})
		builder.Reset()
	}

	return nil
}

func parseInstruction(state *parserState, resolver FileResolver) error {
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
					appendContent(state.pf, builder.String(),
						Location{Origin: state.file, Row: instrFirstRow, Col: instrFirstCol},
						resolver)
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

func appendContent(pf *ProcessedFile, text string, location Location, resolver FileResolver) {
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
		content := createInstruction(parts[0], parts[1], pf.currentScope(), location, text, resolver)
		if content != nil {
			pf.AppendContent(content)
		}
	}
}

func createInstruction(name, arg string, scope Scope, location Location,
	original string, resolver FileResolver) Content {
	switch name {
	case "include":
		return NewIncludeInstruction(arg, location, scope, resolver)
	case "define":
		return NewVariable(arg, location, original, scope)
	case "eval":
		return NewExpression(arg, location, original, scope)
	case "if":
		return NewIfInstruction(arg, location, original)
	case "for":
		return NewForInstruction(arg, location, original, resolver)
	}

	log.Printf("WARNING: (%s) Unknown instruction: %s", location.String(), name)
	return unevaluatedExpression(original)
}

func WriteTo(dir string, filesMap WebFilesMap) error {
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

func writeFile(file, targetFile string, wf WebFile, filesMap WebFilesMap) error {
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
	for _, c := range wf.Processed.contents {
		err := c.Write(w, filesMap, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *StringContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	_, err := writer.Write([]byte(c.Text))
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func (c *StringContent) String() string {
	return fmt.Sprintf("StringContent{%s}", c.Text)
}

func (wf *WebFile) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	return writeContents(wf.Processed, writer, files, inclusionChain)
}

func (wf *WebFile) runSideEffects(files WebFilesMap, inclusionChain []InclusionChainItem) {
	runSideEffects(wf.Processed, files, inclusionChain)
}

func writeContents(cc ContentContainer, writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	for _, c := range cc.GetContents() {
		err := c.Write(writer, files, inclusionChain)
		if err != nil {
			return err
		}
	}
	return nil
}

func runSideEffects(container ContentContainer, files WebFilesMap, inclusionChain []InclusionChainItem) {
	for _, c := range container.GetContents() {
		switch sf := c.(type) {
		case SideEffectContent:
			sf.Run(files, inclusionChain)
		}
	}
}

func inclusionChainToString(locations []InclusionChainItem) string {
	var b strings.Builder
	b.WriteRune('[')
	last := len(locations) - 1
	for i, loc := range locations {
		b.WriteString(loc.Location.String())
		if i != last {
			b.WriteString(" -> ")
		}
	}
	b.WriteRune(']')
	return b.String()
}
