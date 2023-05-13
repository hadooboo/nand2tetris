package jacktokenizer

type TokenType int

const (
	KEYWORD TokenType = iota
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

const (
	CLASS       = "class"
	METHOD      = "method"
	FUNCTION    = "function"
	CONSTRUCTOR = "constructor"
	INT         = "int"
	BOOLEAN     = "boolean"
	CHAR        = "char"
	VOID        = "void"
	VAR         = "var"
	STATIC      = "static"
	FIELD       = "field"
	LET         = "let"
	DO          = "do"
	IF          = "if"
	ELSE        = "else"
	WHILE       = "while"
	RETURN      = "return"
	TRUE        = "true"
	FALSE       = "false"
	NULL        = "null"
	THIS        = "this"
)

var allKeywords = []string{
	CLASS, METHOD, FUNCTION, CONSTRUCTOR, INT, BOOLEAN, CHAR, VOID, VAR, STATIC,
	FIELD, LET, DO, IF, ELSE, WHILE, RETURN, TRUE, FALSE, NULL, THIS,
}
