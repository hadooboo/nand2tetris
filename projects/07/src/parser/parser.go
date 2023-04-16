package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type ParserInterface interface {
	HasMoreCommands() bool
	Advance()
	CommandType() CommandType
	Arg1() string
	Arg2() int
	Close()
}

type CommandType int

const (
	C_ARITHMETIC CommandType = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

var _ = ParserInterface(&Parser{})

type Parser struct {
	stream *bufio.Reader
	f      *os.File

	line        string
	nextLine    string
	commandType CommandType
	arg1        string
	arg2        int
}

func NewParser(fileName string) *Parser {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	return &Parser{
		stream: bufio.NewReader(f),
		f:      f,
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
		line = strings.Join(strings.Fields(line), " ")
		if len(line) == 0 {
			continue
		}
		r.nextLine = line
		return true
	}
}

func (r *Parser) Advance() {
	r.line = r.nextLine
	tokens := strings.Fields(r.line)
	switch tokens[0] {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		r.commandType = C_ARITHMETIC
	case "push":
		r.commandType = C_PUSH
	case "pop":
		r.commandType = C_POP
	}
	switch r.commandType {
	case C_ARITHMETIC:
		r.arg1 = tokens[0]
	case C_PUSH, C_POP:
		r.arg1 = tokens[1]
		arg2, err := strconv.ParseInt(tokens[2], 10, 64)
		if err != nil {
			panic(err)
		}
		r.arg2 = int(arg2)
	}
}

func (r *Parser) CommandType() CommandType {
	return r.commandType
}

func (r *Parser) Arg1() string {
	return r.arg1
}

func (r *Parser) Arg2() int {
	return r.arg2
}

func (r *Parser) Close() {
	r.f.Close()
}
