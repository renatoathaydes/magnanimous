package mg

import (
	"bufio"
	"fmt"
	"io"
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
		state.pf.AppendContent(&StringContent{Text: content})
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
			appendInstructionContent(state.pf, content,
				Location{Origin: state.file, Row: instrFirstRow, Col: instrFirstCol},
				resolver)
		}
		return nil
	}

	return NewError(Location{Origin: state.file, Row: state.row, Col: state.col}, ParseError,
		fmt.Sprintf("instruction started at (%d:%d) was not properly closed with '}}'",
			instrFirstRow, instrFirstCol))
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
