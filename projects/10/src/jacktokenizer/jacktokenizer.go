package jacktokenizer

import (
	"os"
	"strconv"
	"strings"
)

type JackTokenizerInterface interface {
	HasMoreTokens() bool
	Advance()
	TokenType() TokenType
	Keyword() Keyword
	Symbol() string
	Identifier() string
	IntVal() int
	StringVal() string
}

type TokenType int

const (
	KEYWORD TokenType = iota
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

type Keyword string

const (
	CLASS       Keyword = "class"
	METHOD      Keyword = "method"
	FUNCTION    Keyword = "function"
	CONSTRUCTOR Keyword = "constructor"
	INT         Keyword = "int"
	BOOLEAN     Keyword = "boolean"
	CHAR        Keyword = "char"
	VOID        Keyword = "void"
	VAR         Keyword = "var"
	STATIC      Keyword = "static"
	FIELD       Keyword = "field"
	LET         Keyword = "let"
	DO          Keyword = "do"
	IF          Keyword = "if"
	ELSE        Keyword = "else"
	WHILE       Keyword = "while"
	RETURN      Keyword = "return"
	TRUE        Keyword = "true"
	FALSE       Keyword = "false"
	NULL        Keyword = "null"
	THIS        Keyword = "this"
)

var allKeywords = []Keyword{
	CLASS, METHOD, FUNCTION, CONSTRUCTOR, INT, BOOLEAN, CHAR, VOID, VAR, STATIC,
	FIELD, LET, DO, IF, ELSE, WHILE, RETURN, TRUE, FALSE, NULL, THIS,
}

var allSymbols = []byte{
	'{', '}', '(', ')', '[', ']', '.',
	',', ';', '+', '-', '*', '/', '&',
	'|', '<', '>', '=', '~',
}

var _ = JackTokenizerInterface(&JackTokenizer{})

type JackTokenizer struct {
	line       string
	tokenType  TokenType
	keyword    Keyword
	symbol     string
	identifier string
	intVal     int
	stringVal  string
}

func NewJackTokenizer(readFilePath string) *JackTokenizer {
	b, err := os.ReadFile(readFilePath)
	if err != nil {
		panic(err)
	}

	return &JackTokenizer{
		line: string(b),
	}
}

func (r *JackTokenizer) HasMoreTokens() bool {
	if len(r.line) == 0 {
		return false
	}
	r.line = strings.TrimSpace(r.line)
	for {
		if strings.HasPrefix(r.line, "//") {
			idx := strings.IndexByte(r.line, '\n')
			if idx < 0 {
				r.line = ""
				return false
			}
			r.line = strings.TrimSpace(r.line[idx+1:])
			continue
		}
		break
	}
	for {
		if strings.HasPrefix(r.line, "/*") {
			idx := strings.Index(r.line, "*/")
			if idx < 0 {
				r.line = ""
				return false
			}
			r.line = strings.TrimSpace(r.line[idx+2:])
			continue
		}
		break
	}
	return true
}

func (r *JackTokenizer) Advance() {
	if l, v, ok := startsWithKeyword(r.line); ok {
		r.tokenType = KEYWORD
		r.keyword = v
		r.line = l
		return
	}
	if l, v, ok := startsWithSymbol(r.line); ok {
		r.tokenType = SYMBOL
		r.symbol = v
		r.line = l
		return
	}
	if l, v, ok := startsWithIntVal(r.line); ok {
		r.tokenType = INT_CONST
		r.intVal = v
		r.line = l
		return
	}
	if l, v, ok := startsWithStringVal(r.line); ok {
		r.tokenType = STRING_CONST
		r.stringVal = v
		r.line = l
		return
	}
	if l, v, ok := startsWithIdentifier(r.line); ok {
		r.tokenType = IDENTIFIER
		r.identifier = v
		r.line = l
		return
	}
	panic("failed to parse next token")
}

func startsWithKeyword(line string) (string, Keyword, bool) {
	for _, k := range allKeywords {
		if strings.HasPrefix(line, string(k)) {
			return strings.TrimPrefix(line, string(k)), k, true
		}
	}
	return line, "", false
}

func startsWithSymbol(line string) (string, string, bool) {
	for _, s := range allSymbols {
		if line[0] == s {
			symbol := line[0:1]
			if s == '<' {
				symbol = "&lt;"
			}
			if s == '>' {
				symbol = "&gt;"
			}
			if s == '&' {
				symbol = "&amp;"
			}
			return line[1:], symbol, true
		}
	}
	return line, "", false
}

func startsWithIntVal(line string) (string, int, bool) {
	idx := strings.IndexFunc(line, func(r rune) bool {
		return r < '0' || r > '9'
	})
	if idx == 0 {
		return line, 0, false
	}
	if idx < 0 {
		idx = len(line)
	}
	v, err := strconv.ParseInt(line[:idx], 10, 64)
	if err != nil {
		return line, 0, false
	}
	if v < 0 || v > 32767 {
		return line, 0, false
	}
	return line[idx:], int(v), true
}

func startsWithStringVal(line string) (string, string, bool) {
	if line[0] != '"' {
		return line, "", false
	}
	idx := strings.IndexByte(line[1:], '"')
	if idx < 0 {
		return line, "", false
	}
	return line[idx+2:], line[1 : idx+1], true
}

func startsWithIdentifier(line string) (string, string, bool) {
	if line[0] >= '0' && line[0] <= '9' {
		return line, "", false
	}
	idx := strings.IndexFunc(line, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_')
	})
	if idx == 0 {
		return line, "", false
	}
	if idx < 0 {
		idx = len(line)
	}
	return line[idx:], line[:idx], true
}

func (r *JackTokenizer) TokenType() TokenType {
	return r.tokenType
}

func (r *JackTokenizer) Keyword() Keyword {
	return r.keyword
}

func (r *JackTokenizer) Symbol() string {
	return r.symbol
}

func (r *JackTokenizer) Identifier() string {
	return r.identifier
}

func (r *JackTokenizer) IntVal() int {
	return r.intVal
}

func (r *JackTokenizer) StringVal() string {
	return r.stringVal
}
