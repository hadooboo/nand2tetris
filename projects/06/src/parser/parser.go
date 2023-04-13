package parser

import (
	"bufio"
	"strings"
)

type ParserInterface interface {
	HasMoreCommands() bool
	Advance()
	CommandType() CommandType
	Symbol() string
	Dest() string
	Comp() string
	Jump() string
}

type CommandType int

const (
	A_COMMAND CommandType = iota
	C_COMMAND
	L_COMMAND
)

var _ = ParserInterface(&Parser{})

type Parser struct {
	stream *bufio.Reader

	line        string
	nextLine    string
	commandType CommandType
	symbol      string
	dest        string
	comp        string
	jump        string
}

func NewParser(stream *bufio.Reader) *Parser {
	return &Parser{
		stream: stream,
	}
}

func (r *Parser) HasMoreCommands() bool {
	for {
		line, err := r.stream.ReadString('\n')
		if err != nil {
			return false
		}
		idx := strings.Index(line, "//")
		if idx >= 0 {
			line = line[:idx]
		}
		line = strings.Join(strings.Fields(line), "")
		if len(line) == 0 {
			continue
		}
		r.nextLine = line
		return true
	}
}

func (r *Parser) Advance() {
	r.line = r.nextLine
	if strings.HasPrefix(r.line, "@") {
		r.commandType = A_COMMAND
		r.symbol = r.line[1:]
	} else if strings.HasPrefix(r.line, "(") && strings.HasSuffix(r.line, ")") {
		r.commandType = L_COMMAND
		r.symbol = r.line[1 : len(r.line)-1]
	} else {
		r.commandType = C_COMMAND
		destIdx := strings.Index(r.line, "=")
		if destIdx >= 0 {
			r.dest = r.line[:destIdx]
		} else {
			r.dest = ""
		}
		jumpIdx := strings.Index(r.line, ";")
		if jumpIdx >= 0 {
			r.jump = r.line[jumpIdx+1:]
		} else {
			r.jump = ""
		}
		if destIdx >= 0 && jumpIdx >= 0 {
			r.comp = r.line[destIdx+1 : jumpIdx]
		} else if destIdx >= 0 {
			r.comp = r.line[destIdx+1:]
		} else if jumpIdx >= 0 {
			r.comp = r.line[:jumpIdx]
		} else {
			r.comp = ""
		}
	}
}

func (r *Parser) CommandType() CommandType {
	return r.commandType
}

func (r *Parser) Symbol() string {
	return r.symbol
}

func (r *Parser) Dest() string {
	return r.dest
}

func (r *Parser) Comp() string {
	return r.comp
}

func (r *Parser) Jump() string {
	return r.jump
}
