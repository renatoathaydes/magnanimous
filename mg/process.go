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

func (mag *Magnanimous) ReadAll() (WebFilesMap, error) {
	processedDir := filepath.Join(mag.SourcesDir, "processed")
	staticDir := filepath.Join(mag.SourcesDir, "static")

	procFiles, staticFiles, otherFiles := collectFiles(mag.SourcesDir, processedDir, staticDir)
	webFiles := WebFilesMap{
		WebFiles: make(map[string]WebFile, len(procFiles)+len(staticFiles)+len(otherFiles)),
	}

	err := ProcessAll(procFiles, processedDir, mag.SourcesDir, &webFiles)
	if err != nil {
		return webFiles, err
	}
	CopyAll(&staticFiles, staticDir, webFiles)
	AddNonWritables(&otherFiles, mag.SourcesDir, webFiles)

	return webFiles, nil
}

func ProcessAll(files []string, basePath, sourcesDir string, webFiles *WebFilesMap) error {
	resolver := DefaultFileResolver{BasePath: sourcesDir, Files: webFiles}
	for _, file := range files {
		wf, err := ProcessFile(file, basePath, &resolver)
		if err != nil {
			return err
		}
		webFiles.WebFiles[file] = *wf
	}
	if globalCtx, ok := webFiles.WebFiles[filepath.Join(basePath, "_global_context")]; ok {
		globalCtx.runSideEffects(webFiles, nil)
		var globalContext RootScope = globalCtx.Processed.Context()
		webFiles.GlobalContext = globalContext
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

func appendInstructionContent(pf *ProcessedFile, text string, location Location, resolver FileResolver) {
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
	switch strings.TrimSpace(name) {
	case "include":
		return NewIncludeInstruction(arg, location, original, scope, resolver)
	case "define":
		return NewVariable(arg, location, original, scope)
	case "eval":
		return NewExpression(arg, location, original, scope)
	case "if":
		return NewIfInstruction(arg, location, original, scope)
	case "for":
		return NewForInstruction(arg, location, original, scope, resolver)
	case "doc":
		return nil
	case "component":
		return NewComponentInstruction(arg, location, original, scope, resolver)
	}

	log.Printf("WARNING: (%s) Unknown instruction: '%s'", location.String(), name)
	return unevaluatedExpression(original)
}

func WriteTo(dir string, filesMap WebFilesMap) error {
	err := os.MkdirAll(dir, 0770)
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	for file, wf := range filesMap.WebFiles {
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
	return writeContents(wf.Processed, writer, files, inclusionChain, false)
}

func (wf *WebFile) runSideEffects(files *WebFilesMap, inclusionChain []InclusionChainItem) {
	runSideEffects(wf.Processed, files, inclusionChain)
}

func writeContents(cc ContentContainer, writer io.Writer, files WebFilesMap,
	inclusionChain []InclusionChainItem, runSideEffectsFirst bool) error {
	if runSideEffectsFirst {
		runSideEffects(cc, &files, inclusionChain)
	}
	for _, c := range cc.GetContents() {
		if runSideEffectsFirst {
			// only skip define content because not all SideEffectContent does only side-effect
			if _, skip := c.(*DefineContent); skip {
				continue
			}
		}
		err := c.Write(writer, files, inclusionChain)
		if err != nil {
			return err
		}
	}
	return nil
}

func runSideEffects(container ContentContainer, files *WebFilesMap, inclusionChain []InclusionChainItem) {
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
