package jacktokenizer

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

type JackTokenizerInterface interface {
	HasMoreTokens() bool
	Advance()
	Retreat()
	TokenType() TokenType
	Keyword() string
	Symbol() string
	IntVal() int
	StringVal() string
	Identifier() string
}

var _ = JackTokenizerInterface(&JackTokenizer{})

type Token struct {
	tokenType  TokenType
	keyword    string
	symbol     string
	intVal     int
	stringVal  string
	identifier string
}

type JackTokenizer struct {
	tokens []Token
	p      int
}

func NewJackTokenizer(readFilePath string) *JackTokenizer {
	b, err := os.ReadFile(readFilePath)
	if err != nil {
		panic(err)
	}

	jackTokenizer := &JackTokenizer{
		tokens: make([]Token, 0),
		p:      -1,
	}

	jackTokenizer.tokenize(string(b))

	return jackTokenizer
}

var (
	slashSlashCommentRegexStr = `//(.*?)\n`
	slashStarCommentRegexStr  = `/\*((.|\n)*?)\*/`
	commentRegex              = regexp.MustCompile(strings.Join([]string{
		slashSlashCommentRegexStr,
		slashStarCommentRegexStr,
	}, "|"))

	keywordRegexStr    = `\b` + strings.Join(allKeywords, `\b|\b`) + `\b`
	keywordRegex       = regexp.MustCompile(keywordRegexStr)
	symbolRegexStr     = `{|}|\(|\)|\[|\]|\.|,|;|\+|\-|\*|/|&|\||<|>|=|~`
	symbolRegex        = regexp.MustCompile(symbolRegexStr)
	intRegexStr        = `\b\d+\b`
	intRegex           = regexp.MustCompile(intRegexStr)
	stringRegexStr     = `"(.*?)"`
	stringRegex        = regexp.MustCompile(stringRegexStr)
	identifierRegexStr = `\b[a-zA-Z_]\w*\b`
	identifierRegex    = regexp.MustCompile(identifierRegexStr)
	tokenRegex         = regexp.MustCompile(strings.Join([]string{
		keywordRegexStr,
		symbolRegexStr,
		intRegexStr,
		stringRegexStr,
		identifierRegexStr,
	}, "|"))
)

func (r *JackTokenizer) tokenize(data string) {
	data = commentRegex.ReplaceAllString(data, "")
	tokens := tokenRegex.FindAllString(data, -1)
	for _, t := range tokens {
		if keywordRegex.MatchString(t) {
			r.tokens = append(r.tokens, Token{
				tokenType: KEYWORD,
				keyword:   t,
			})
		} else if intRegex.MatchString(t) {
			v, _ := strconv.ParseInt(t, 10, 64)
			r.tokens = append(r.tokens, Token{
				tokenType: INT_CONST,
				intVal:    int(v),
			})
		} else if stringRegex.MatchString(t) {
			r.tokens = append(r.tokens, Token{
				tokenType: STRING_CONST,
				stringVal: t[1 : len(t)-1],
			})
		} else if symbolRegex.MatchString(t) {
			r.tokens = append(r.tokens, Token{
				tokenType: SYMBOL,
				symbol:    t,
			})
		} else if identifierRegex.MatchString(t) {
			r.tokens = append(r.tokens, Token{
				tokenType:  IDENTIFIER,
				identifier: t,
			})
		}
	}
}

func (r *JackTokenizer) HasMoreTokens() bool {
	return r.p+1 < len(r.tokens)
}

func (r *JackTokenizer) Advance() {
	r.p++
}

func (r *JackTokenizer) Retreat() {
	r.p--
}

func (r *JackTokenizer) TokenType() TokenType {
	return r.tokens[r.p].tokenType
}

func (r *JackTokenizer) Keyword() string {
	return r.tokens[r.p].keyword
}

func (r *JackTokenizer) Symbol() string {
	return r.tokens[r.p].symbol
}

func (r *JackTokenizer) IntVal() int {
	return r.tokens[r.p].intVal
}

func (r *JackTokenizer) StringVal() string {
	return r.tokens[r.p].stringVal
}

func (r *JackTokenizer) Identifier() string {
	return r.tokens[r.p].identifier
}
