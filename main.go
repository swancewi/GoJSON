package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Input struct {
	input    string
	position int
}

func ParseJSON(input string) (any, error) {
	newInput := Input{input: input, position: 0}
	newInput.skipWhitespace()
	return newInput.parseObject()
}

func (input *Input) parseObject() (any, error) {
	var theValue any
	var err error
	switch nextValue := input.peek(); {
	case nextValue == '"':
		theValue, err = input.parseString()
	case nextValue >= '1' && nextValue <= '9':
		theValue, err = input.parseNumber()
	case nextValue == '{':
		theValue, err = input.parseMap()
	case nextValue == '[':
		theValue, err = input.parseArray()
	case nextValue == 't' || nextValue == 'f' || nextValue == 'n':
		theValue, err = input.parseBooleanOrNil()
		fallthrough
	default:
		if input.position == 0 && nextValue >= 'A' && nextValue <= 'z' {
			return input.remaining(), nil
		}
		if err != nil {
			err = fmt.Errorf("couldn't figure out how to parse JSON for %s", input.remaining())
		}
	}
	if err != nil {
		return nil, err
	}
	return theValue, nil
}

func (input *Input) parseArray() (any, error) {
	input.position++
	returnArray := []any{}
	canContinue := true

	for {
		input.skipWhitespace()
		if input.peek() == ']' {
			input.position++
			break
		} else if input.peek() == ',' {
			canContinue = true
			input.position++
			input.skipWhitespace()
		}
		if !canContinue {
			return nil, errors.New("unable to continue in array parse")
		}
		canContinue = false
		if val, err := input.parseObject(); err != nil {
			return nil, err
		} else {
			returnArray = append(returnArray, val)
		}
	}

	return returnArray, nil
}

func (input *Input) parseMap() (any, error) {
	res := make(map[string]any)
	input.position++ // could assert on that it should be '{'
	for {
		nextValue := true
		for nextValue {
			nextValue = false
			input.skipWhitespace()
			key, err := input.parseString()
			if err != nil {
				return nil, err
			}

			input.skipWhitespace()
			if input.peek() != ':' {
				return nil, errors.New("missing semicolon")
			}
			input.position++

			input.skipWhitespace()
			theValue, err := input.parseObject()
			if err != nil {
				return nil, err
			}
			res[key] = theValue
			input.skipWhitespace()
			if input.peek() == ',' {
				nextValue = true
				input.position++
			}
		}
		input.skipWhitespace()
		if input.peek() == '}' {
			input.position++
			break
		} else {
			return nil, errors.New("missing end bracket")
		}
	}

	// fmt.Println(res)

	return res, nil
}

func (input *Input) parseBooleanOrNil() (any, error) {
	if len(input.remaining()) < 5 {
		return nil, fmt.Errorf("not true, false, null %s", input.remaining())
	} // TODO - also need check of false length

	if input.peek() == 't' && input.input[input.position:input.position+4] == "true" {
		input.position += 4
		return true, nil
	}
	if input.peek() == 'n' && input.input[input.position:input.position+4] == "null" {
		input.position += 4
		return nil, nil
	}
	if input.peek() == 'f' && input.input[input.position:input.position+5] == "false" {
		input.position += 5
		return false, nil
	}

	return nil, fmt.Errorf("non-nil or non-boolean found starting at %c", input.peek())
}

func (input *Input) remaining() string {
	return input.input[input.position:]
}

func (input *Input) parseNumber() (any, error) {
	startValue := input.position
	for unicode.IsDigit(input.peek()) {
		input.position++
	}
	// fmt.Printf("attemping to parse %d:%d %s", startValue, input.position, input.input[startValue:input.position])
	return strconv.ParseFloat(input.input[startValue:input.position], 64)
}

func (input *Input) parseString() (string, error) {
	if input.peek() != '"' {
		return "", fmt.Errorf("expected \" found %c", input.peek())
	}
	input.position++
	startChar := input.position
	for input.peek() != '"' {
		input.position++
	}
	val := input.input[startChar:input.position]
	input.position++
	return val, nil
}

func (input *Input) skipWhitespace() {
	for {
		nextRune, size := utf8.DecodeRuneInString(input.input[input.position:])
		if nextRune != ' ' {
			return
		}
		input.position += size
		if input.position > len(input.input) {
			return
		}
	}
}

func (input Input) peek() rune {
	nextRune, _ := utf8.DecodeRuneInString(input.input[input.position:])
	return nextRune
}
