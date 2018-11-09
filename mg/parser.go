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
	reader := state.reader
	builder := state.builder

	//// all functions handling special characters return nil if they handle the next rune,
	//// or the next rune if it was not handled

	onEscapedReturn := func() (*rune, error) {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}
		state.col++
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

	onOpenBracket := func() (*rune, error) {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}
		state.col++
		if err != nil {
			return nil, &MagnanimousError{message: err.Error(), Code: IOError}
		}
		if r == '{' {
			if builder.Len() > 0 {
				state.pf.AppendContent(&StringContent{Text: builder.String()})
				builder.Reset()
			}
			return nil, parseInstruction(state, resolver)
		}
		_, err = builder.WriteRune('{')
		return &r, err
	}

	onEscape := func() (*rune, error) {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}
		state.col++
		if err != nil {
			return nil, &MagnanimousError{message: err.Error(), Code: IOError}
		}
		switch r {
		case '{':
			// don't treat specially, just let it be written
			_, err = builder.WriteRune('{')
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
		state.col++

	parseRune:
		if err == io.EOF {
			break
		}
		if err != nil {
			return &MagnanimousError{message: err.Error(), Code: IOError}
		}
		var nextRune *rune
		switch r {
		case '\n':
			state.row++
			state.col = 1
			_, err = builder.WriteRune('\n')
		case '{':
			nextRune, err = onOpenBracket()
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

	// append any pending content to the builder
	if builder.Len() > 0 {
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
