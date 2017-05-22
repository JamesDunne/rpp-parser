package rpp

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type rppParser struct {
	reader *bufio.Reader

	line            string
	isEOF           bool
	trim_line       string
	expected_indent int
	actual_indent   int
	trim_i          int
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
	parser.trim_i = 0
	parser.actual_indent = len(parser.line) - len(parser.trim_line)

	fmt.Printf("line[%d]=%v\n", parser.actual_indent, parser.trim_line)

	return nil
}

func (parser *rppParser) parseWord() string {
	i := parser.trim_i
	w := 0
	start := i

	for ; i < len(parser.trim_line); i += w {
		c, width := utf8.DecodeRuneInString(parser.trim_line[i:])
		w = width

		if unicode.IsSpace(c) {
			break
		}
	}
	parser.trim_i = i

	return parser.trim_line[start:parser.trim_i]
}

func (parser *rppParser) parseNumber() string {
	return ""
}

func (parser *rppParser) parseProject() (project *Project, err error) {
	if parser.trim_line[0] != '<' {
		return nil, fmt.Errorf("Expected '<'")
	}

	project, err = &Project{}, nil
	parser.expected_indent += 2

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

		directive := parser.parseWord()
		if len(directive) == 0 {
			continue
		}

		if directive[0] == '<' {
			name := directive[1:]
			if name == "TRACK" {
				var track *Track
				track, err = parser.parseTrack()
				if err != nil {
					return
				}
				project.Tracks = append(project.Tracks, track)
			} else {
				fmt.Printf("%s\n", name)
				parser.skipUnknownBlock()
			}
		}
	}

	return
}

func (parser *rppParser) skipUnknownBlock() (err error) {
	parser.expected_indent += 2

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

		directive := parser.parseWord()
		if len(directive) == 0 {
			continue
		}

		if directive[0] == '<' {
			parser.skipUnknownBlock()
			continue
		}
	}
}

func (parser *rppParser) parseTrack() (track *Track, err error) {
	if parser.trim_line[0] != '<' {
		return nil, fmt.Errorf("Expected '<'")
	}

	track, err = &Track{}, nil
	parser.expected_indent += 2

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

		directive := parser.parseWord()
		if len(directive) == 0 {
			continue
		}

		if directive[0] == '<' {
			name := directive[1:]
			if name == "FXCHAIN" {
				var fxChain *FXChain
				fxChain, err = parser.parseFXChain()
				if err != nil {
					return
				}
				track.FXChain = fxChain
			} else {
				parser.skipUnknownBlock()
			}
			continue
		}
	}

	return
}

func (parser *rppParser) parseFXChain() (chain *FXChain, err error) {
	if parser.trim_line[0] != '<' {
		return nil, fmt.Errorf("Expected '<'")
	}

	chain, err = &FXChain{}, nil
	parser.expected_indent += 2

	//var fx *FX = nil

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

		name := parser.parseWord()

		if name[0] == '<' {
			//if name == "VST" {
			//	var vst *VST
			//	vst, err = parser.parseVST()
			//	if err != nil {
			//		return
			//	}
			//	//fx.VST = vst
			//}
			parser.skipUnknownBlock()
			continue
		}
	}

	return
}

func ParseRPP(reader *bufio.Reader) (project *Project, err error) {
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

	return parser.parseProject()
}
