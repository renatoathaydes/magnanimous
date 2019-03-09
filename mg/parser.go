package mg

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

type parserState struct {
	file   string
	reader *bufio.Reader
	row    uint32
	col    uint32
	//pf           *ProcessedFile
	builder      *strings.Builder
	contentStack []ContentContainer
}

func (state *parserState) append(content Content) {
	container := state.contentStack[len(state.contentStack)-1]
	container.AppendContent(content)
	if c, ok := content.(ContentContainer); ok {
		state.contentStack = append(state.contentStack, c)
	}
}

func (state *parserState) dropStackItem() bool {
	if len(state.contentStack) > 1 {
		state.contentStack = state.contentStack[:len(state.contentStack)-1]
		return true
	}
	return false
}

func parseText(state *parserState, resolver FileResolver) error {
begin:
	eof, err := parseUntilDoubleRunes('{', state)
	if err != nil {
		return err
	}
	content := state.builder.String()
	state.builder.Reset()
	includeContent := len(content) > 0
	if eof {
		// last part, only include content if not only whitespaces
		includeContent = len(strings.TrimSpace(content)) > 0
	}
	if includeContent {
		state.append(&StringContent{Text: content})
	}
	if !eof {
		err = parseInstruction(state, resolver)
		if err != nil {
			return err
		} else {
			goto begin
		}
	}
	return nil
}

func parseInstruction(state *parserState, resolver FileResolver) error {
	instrFirstRow := state.row
	instrFirstCol := state.col - 2

	eof, err := parseUntilDoubleRunes('}', state)
	if err != nil {
		return err
	}
	if !eof {
		content := state.builder.String()
		state.builder.Reset()

		if len(content) > 0 {
			appendInstructionContent(state, content,
				&Location{Origin: state.file, Row: instrFirstRow, Col: instrFirstCol},
				resolver)
		}
		return nil
	}

	return NewError(Location{Origin: state.file, Row: state.row, Col: state.col}, ParseError,
		fmt.Sprintf("instruction started at (%d:%d) was not properly closed with '}}'",
			instrFirstRow, instrFirstCol))
}

func appendInstructionContent(state *parserState, text string, location *Location, resolver FileResolver) {
	parts := strings.SplitN(strings.TrimSpace(text), " ", 2)
	switch len(parts) {
	case 0:
		// nothing to do
	case 1:
		if parts[0] == "end" {
			wasDropped := state.dropStackItem()
			if !wasDropped {
				log.Printf("WARNING: (%s) %s", location, "end instruction does not match any open scope")
				state.append(unevaluatedExpression(text))
			}
		} else {
			log.Printf("WARNING: (%s) Instruction missing argument: %s", location.String(), text)
			state.append(unevaluatedExpression(text))
		}
	case 2:
		content := createInstruction(parts[0], parts[1], location, text, resolver)
		if content != nil {
			state.append(content)
		}
	}
}

func createInstruction(name, arg string, location *Location,
	original string, resolver FileResolver) Content {
	switch strings.TrimSpace(name) {
	case "include":
		return NewIncludeInstruction(arg, location, original, resolver)
	case "define":
		return NewVariable(arg, location, original)
	case "eval":
		return NewExpression(arg, location, original)
	case "if":
		return NewIfInstruction(arg, location, original)
	case "for":
		return NewForInstruction(arg, location, original, resolver)
	case "doc":
		return nil
	case "component":
		return NewComponentInstruction(arg, location, original, resolver)
	case "slot":
		return NewSlotInstruction(arg, location, original)
	}

	log.Printf("WARNING: (%s) Unknown instruction: '%s'", location.String(), name)
	return unevaluatedExpression(original)
}

func parseUntilDoubleRunes(specialRune rune, state *parserState) (bool, error) {
	reader := state.reader
	builder := state.builder

	//// All functions handling special characters return nil if they handle the next rune,
	//// or the next rune if it was not handled.
	//// If handling the rune, the func must increase the state.col counter.

	onEscapedReturn := func() (*rune, error) {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}
		if err != nil {
			return nil, &MagnanimousError{message: err.Error(), Code: IOError}
		}
		if r == '\n' {
			// forget both the return and the new-line
			state.row++
			state.col = 1
			return nil, nil
		}
		// forget just the return rune and continue with the next one
		return &r, nil
	}

	onSpecialRune := func() (*rune, bool, error) {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, true, nil
		}
		if err != nil {
			return nil, false, &MagnanimousError{message: err.Error(), Code: IOError}
		}
		if r == specialRune {
			state.col++
			return nil, false, nil
		}
		_, err = builder.WriteRune(specialRune)
		return &r, false, err
	}

	onEscape := func() (*rune, error) {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}
		if err != nil {
			return nil, &MagnanimousError{message: err.Error(), Code: IOError}
		}

		switch r {
		case specialRune, '\\':
			// don't treat specially, just let it be written
			_, err = builder.WriteRune(r)
			state.col++
			return nil, err
		case '\n':
			state.row++
			state.col = 1
			// forget the new line
			return nil, nil
		case '\r':
			return onEscapedReturn()
		}
		_, err = builder.WriteRune('\\')
		return &r, err
	}

	for {
		r, _, err := reader.ReadRune()

	parseRune:
		if err == io.EOF {
			break
		}
		state.col++
		if err != nil {
			return false, &MagnanimousError{message: err.Error(), Code: IOError}
		}
		var nextRune *rune
		switch r {
		case specialRune:
			var eof bool
			nextRune, eof, err = onSpecialRune()
			if err == nil && eof {
				break
			}
			if err == nil && nextRune == nil {
				return false, nil
			}
		case '\n':
			state.row++
			state.col = 1
			_, err = builder.WriteRune('\n')
		case '\\':
			nextRune, err = onEscape()
		default:
			_, err = builder.WriteRune(r)
		}

		if nextRune != nil || err != nil {
			if nextRune != nil {
				r = *nextRune
			}
			goto parseRune
		}
	}

	return true, nil
}
