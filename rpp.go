package rpp

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
)

type ParameterValue struct {
	Type  int
	Value interface{}
}

type NamedParameter struct {
	Name string
	ParameterValue
}

type RPP struct {
	Name                 string
	PositionalParameters []ParameterValue
	NamedParameters      []NamedParameter
	Children             []*RPP
}

type rppParser struct {
	reader *bufio.Reader

	line            string
	isEOF           bool
	trim_line       string
	expected_indent int
	actual_indent   int
}

func (parser *rppParser) grabLine() (err error) {
	var isPrefix bool
	var lineBytes []byte
	lineBytes, isPrefix, err = parser.reader.ReadLine()
	if err != nil {
		return err
	}
	if isPrefix {
		return fmt.Errorf("isPrefix=true")
	}
	if lineBytes == nil {
		parser.isEOF = true
		parser.line = ""
		parser.trim_line = ""
		parser.actual_indent = 0
		return nil
	}

	parser.line = string(lineBytes)
	parser.trim_line = strings.TrimLeftFunc(parser.line, unicode.IsSpace)
	parser.actual_indent = len(parser.line) - len(parser.trim_line)

	fmt.Printf("line[%d]=%v\n", parser.actual_indent, parser.trim_line)

	return nil
}

func (parser *rppParser) parseChild() (node *RPP, err error) {
	if parser.trim_line[0] != '<' {
		return nil, fmt.Errorf("Expected '<'")
	}

	node, err = &RPP{}, nil
	parser.expected_indent += 2

	// Extract name:
	wrk := parser.trim_line[1:]
	var c int32
	i := 0
	for i, c = range wrk {
		if unicode.IsSpace(c) {
			i--
			break
		}
	}

	node.Name = wrk[0:i+1]
	fmt.Printf("name: %s\n", node.Name)

	wrk = wrk[i+1:]

	// TODO: parse remaining positional arguments

	for {
		if err = parser.grabLine(); err != nil {
			return
		}
		if parser.isEOF {
			return
		}

		// Close current:
		if parser.trim_line[0] == '>' {
			if parser.actual_indent != parser.expected_indent-2 {
				err = fmt.Errorf("Invalid indentation for closing '>' char")
				return
			}
			parser.expected_indent -= 2
			return
		}

		if parser.trim_line[0] == '<' {
			var child *RPP
			child, err = parser.parseChild()
			node.Children = append(node.Children, child)
			if err != nil {
				return
			}
		}
	}

	return
}

func ParseRPP(reader *bufio.Reader) (project *RPP, err error) {
	parser := &rppParser{
		reader: reader,
	}

	// Read first line:
	if err = parser.grabLine(); err != nil {
		return nil, err
	}
	if parser.isEOF {
		return nil, nil
	}

	return parser.parseChild()
}
