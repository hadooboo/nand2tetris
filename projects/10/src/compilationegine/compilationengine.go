package compilationegine

import (
	"bufio"
	"fmt"
	"os"

	"jaehonam.com/nand2tetris/project/10/jacktokenizer"
)

const TAB = "  "

type CompilationEngineInterface interface {
	PrintXML()
	PrintTokensXML()
	Close()
}

var _ = CompilationEngineInterface(&CompilationEngine{})

type CompilationEngine struct {
	stream *bufio.Writer
	f      *os.File

	jt jacktokenizer.JackTokenizerInterface
}

func NewCompilationEngine(readFilePath, writeFilePath string) *CompilationEngine {
	f, err := os.Create(writeFilePath)
	if err != nil {
		panic(err)
	}

	return &CompilationEngine{
		stream: bufio.NewWriter(f),
		f:      f,
		jt:     jacktokenizer.NewJackTokenizer(readFilePath),
	}
}

func (r *CompilationEngine) PrintXML() {
	r.compileClass("")
}

func (r *CompilationEngine) compileClass(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<class>")
	r.mustPrintKeyword(inner)
	r.mustPrintIdentifier(inner)
	r.mustPrintSymbol(inner)
	for r.isClassVarDec() {
		r.compileClassVarDec(inner)
	}
	for r.isSubroutineDec() {
		r.compileSubroutineDec(inner)
	}
	r.mustPrintSymbol(inner)
	fmt.Fprintln(r.stream, depth+"</class>")
}

func (r *CompilationEngine) compileClassVarDec(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<classVarDec>")
	r.mustPrintKeyword(inner)
	r.mustPrintType(inner)
	r.mustPrintIdentifier(inner)
	for r.isSymbol(",") {
		r.mustPrintSymbol(inner)
		r.mustPrintIdentifier(inner)
	}
	r.mustPrintSymbol(inner)
	fmt.Fprintln(r.stream, depth+"</classVarDec>")
}

func (r *CompilationEngine) compileSubroutineDec(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<subroutineDec>")
	r.mustPrintKeyword(inner)
	r.mustPrintType(inner)
	r.mustPrintIdentifier(inner)
	r.mustPrintSymbol(inner)
	r.compileParameterList(depth + TAB)
	r.mustPrintSymbol(inner)
	r.compileSubroutineBody(depth + TAB)
	fmt.Fprintln(r.stream, depth+"</subroutineDec>")
}

func (r *CompilationEngine) compileParameterList(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<parameterList>")
	if r.isType() {
		r.mustPrintType(inner)
		r.mustPrintIdentifier(inner)
		for r.isSymbol(",") {
			r.mustPrintSymbol(inner)
			r.mustPrintType(inner)
			r.mustPrintIdentifier(inner)
		}
	}
	fmt.Fprintln(r.stream, depth+"</parameterList>")
}

func (r *CompilationEngine) compileSubroutineBody(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<subroutineBody>")
	r.mustPrintSymbol(inner)
	for r.isVarDec() {
		r.compileVarDec(inner)
	}
	r.compileStatements(inner)
	r.mustPrintSymbol(inner)
	fmt.Fprintln(r.stream, depth+"</subroutineBody>")
}

func (r *CompilationEngine) compileVarDec(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<varDec>")
	r.mustPrintKeyword(inner)
	r.mustPrintType(inner)
	r.mustPrintIdentifier(inner)
	for r.isSymbol(",") {
		r.mustPrintSymbol(inner)
		r.mustPrintIdentifier(inner)
	}
	r.mustPrintSymbol(inner)
	fmt.Fprintln(r.stream, depth+"</varDec>")
}

func (r *CompilationEngine) compileStatements(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<statements>")
	for r.isStatement() {
		switch {
		case r.isKeyword(jacktokenizer.LET):
			r.compileLetStatement(inner)
		case r.isKeyword(jacktokenizer.IF):
			r.compileIfStatement(inner)
		case r.isKeyword(jacktokenizer.WHILE):
			r.compileWhileStatement(inner)
		case r.isKeyword(jacktokenizer.DO):
			r.compileDoStatement(inner)
		case r.isKeyword(jacktokenizer.RETURN):
			r.compileReturnStatement(inner)
		}
	}
	fmt.Fprintln(r.stream, depth+"</statements>")
}

func (r *CompilationEngine) compileLetStatement(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<letStatement>")
	r.mustPrintKeyword(inner)
	r.mustPrintIdentifier(inner)
	if r.isSymbol("[") {
		r.mustPrintSymbol(inner)
		r.compileExpression(inner)
		r.mustPrintSymbol(inner)
	}
	r.mustPrintSymbol(inner)
	r.compileExpression(inner)
	r.mustPrintSymbol(inner)
	fmt.Fprintln(r.stream, depth+"</letStatement>")
}

func (r *CompilationEngine) compileIfStatement(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<ifStatement>")
	r.mustPrintKeyword(inner)
	r.mustPrintSymbol(inner)
	r.compileExpression(inner)
	r.mustPrintSymbol(inner)
	r.mustPrintSymbol(inner)
	r.compileStatements(inner)
	r.mustPrintSymbol(inner)
	if r.isKeyword(jacktokenizer.ELSE) {
		r.mustPrintKeyword(inner)
		r.mustPrintSymbol(inner)
		r.compileStatements(inner)
		r.mustPrintSymbol(inner)
	}
	fmt.Fprintln(r.stream, depth+"</ifStatement>")
}

func (r *CompilationEngine) compileWhileStatement(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<whileStatement>")
	r.mustPrintKeyword(inner)
	r.mustPrintSymbol(inner)
	r.compileExpression(inner)
	r.mustPrintSymbol(inner)
	r.mustPrintSymbol(inner)
	r.compileStatements(inner)
	r.mustPrintSymbol(inner)
	fmt.Fprintln(r.stream, depth+"</whileStatement>")
}

func (r *CompilationEngine) compileDoStatement(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<doStatement>")
	r.mustPrintKeyword(inner)
	r.mustPrintIdentifier(inner)
	switch {
	case r.isSymbol("("):
		r.mustPrintSymbol(inner)
		r.compileExpressionList(inner)
		r.mustPrintSymbol(inner)
	case r.isSymbol("."):
		r.mustPrintSymbol(inner)
		r.mustPrintIdentifier(inner)
		r.mustPrintSymbol(inner)
		r.compileExpressionList(inner)
		r.mustPrintSymbol(inner)
	}
	r.mustPrintSymbol(inner)
	fmt.Fprintln(r.stream, depth+"</doStatement>")
}

func (r *CompilationEngine) compileReturnStatement(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<returnStatement>")
	r.mustPrintKeyword(inner)
	if r.isExpression() {
		r.compileExpression(inner)
	}
	r.mustPrintSymbol(inner)
	fmt.Fprintln(r.stream, depth+"</returnStatement>")
}

func (r *CompilationEngine) compileExpression(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<expression>")
	r.compileTerm(inner)
	for r.isOp() {
		r.mustPrintSymbol(inner)
		r.compileTerm(inner)
	}
	fmt.Fprintln(r.stream, depth+"</expression>")
}

func (r *CompilationEngine) compileTerm(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<term>")
	switch {
	case r.isIntegerConstant():
		r.mustPrintIntegerConstant(inner)
	case r.isStringConstant():
		r.mustPrintStringConstant(inner)
	case r.isKeywordConstant():
		r.mustPrintKeyword(inner)
	case r.isIdentifier():
		r.mustPrintIdentifier(inner)
		switch {
		case r.isSymbol("["):
			r.mustPrintSymbol(inner)
			r.compileExpression(inner)
			r.mustPrintSymbol(inner)
		case r.isSymbol("("):
			r.mustPrintSymbol(inner)
			r.compileExpressionList(inner)
			r.mustPrintSymbol(inner)
		case r.isSymbol("."):
			r.mustPrintSymbol(inner)
			r.mustPrintIdentifier(inner)
			r.mustPrintSymbol(inner)
			r.compileExpressionList(inner)
			r.mustPrintSymbol(inner)
		}
	case r.isSymbol("("):
		r.mustPrintSymbol(inner)
		r.compileExpression(inner)
		r.mustPrintSymbol(inner)
	case r.isUnaryOp():
		r.mustPrintSymbol(inner)
		r.compileTerm(inner)
	}
	fmt.Fprintln(r.stream, depth+"</term>")
}

func (r *CompilationEngine) compileExpressionList(depth string) {
	inner := depth + TAB
	fmt.Fprintln(r.stream, depth+"<expressionList>")
	if r.isExpression() {
		r.compileExpression(inner)
		for r.isSymbol(",") {
			r.mustPrintSymbol(inner)
			r.compileExpression(inner)
		}
	}
	fmt.Fprintln(r.stream, depth+"</expressionList>")
}

func (r *CompilationEngine) isClassVarDec() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return (r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.STATIC) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.FIELD)
	}

	return false
}

func (r *CompilationEngine) isType() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return (r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.INT) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.CHAR) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.BOOLEAN) ||
			r.jt.TokenType() == jacktokenizer.IDENTIFIER
	}

	return false
}

func (r *CompilationEngine) isSymbol(symbol string) bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == symbol
	}

	return false
}

func (r *CompilationEngine) isKeyword(keyword string) bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == keyword
	}

	return false
}

func (r *CompilationEngine) isOp() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return (r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "+") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "-") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "*") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "/") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "&") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "|") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "<") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == ">") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "=")
	}

	return false
}

func (r *CompilationEngine) isSubroutineDec() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return (r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.CONSTRUCTOR) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.FUNCTION) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.METHOD)
	}

	return false
}

func (r *CompilationEngine) isVarDec() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.VAR
	}

	return false
}

func (r *CompilationEngine) isStatement() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return (r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.LET) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.IF) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.WHILE) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.DO) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.RETURN)
	}

	return false
}

func (r *CompilationEngine) isExpression() bool {
	return r.isTerm()
}

func (r *CompilationEngine) isTerm() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return r.jt.TokenType() == jacktokenizer.INT_CONST ||
			r.jt.TokenType() == jacktokenizer.STRING_CONST ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.TRUE) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.FALSE) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.NULL) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.THIS) ||
			r.jt.TokenType() == jacktokenizer.IDENTIFIER ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "(") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "-") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "~")
	}

	return false
}

func (r *CompilationEngine) isIntegerConstant() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return r.jt.TokenType() == jacktokenizer.INT_CONST
	}

	return false
}

func (r *CompilationEngine) isStringConstant() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return r.jt.TokenType() == jacktokenizer.STRING_CONST
	}

	return false
}

func (r *CompilationEngine) isKeywordConstant() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return (r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.TRUE) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.FALSE) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.NULL) ||
			(r.jt.TokenType() == jacktokenizer.KEYWORD && r.jt.Keyword() == jacktokenizer.THIS)
	}

	return false
}

func (r *CompilationEngine) isUnaryOp() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return (r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "-") ||
			(r.jt.TokenType() == jacktokenizer.SYMBOL && r.jt.Symbol() == "~")
	}

	return false
}

func (r *CompilationEngine) isIdentifier() bool {
	if r.jt.HasMoreTokens() {
		r.jt.Advance()
		defer r.jt.Retreat()
		return r.jt.TokenType() == jacktokenizer.IDENTIFIER
	}

	return false
}

func (r *CompilationEngine) printKeyword(depth string) {
	fmt.Fprintln(r.stream, depth+fmt.Sprintf("<keyword> %v </keyword>", r.jt.Keyword()))
}

func (r *CompilationEngine) mustPrintKeyword(depth string) {
	r.mustAdvance()
	if r.jt.TokenType() != jacktokenizer.KEYWORD {
		panic("invalid token type")
	}
	r.printKeyword(depth)
}

func (r *CompilationEngine) printSymbol(depth string) {
	symbol := r.jt.Symbol()
	if symbol == "<" {
		symbol = "&lt;"
	}
	if symbol == ">" {
		symbol = "&gt;"
	}
	if symbol == "&" {
		symbol = "&amp;"
	}
	fmt.Fprintln(r.stream, depth+fmt.Sprintf("<symbol> %v </symbol>", symbol))
}

func (r *CompilationEngine) mustPrintSymbol(depth string) {
	r.mustAdvance()
	if r.jt.TokenType() != jacktokenizer.SYMBOL {
		panic("invalid token type")
	}
	r.printSymbol(depth)
}

func (r *CompilationEngine) printIntegerConstant(depth string) {
	fmt.Fprintln(r.stream, depth+fmt.Sprintf("<integerConstant> %v </integerConstant>", r.jt.IntVal()))
}

func (r *CompilationEngine) mustPrintIntegerConstant(depth string) {
	r.mustAdvance()
	if r.jt.TokenType() != jacktokenizer.INT_CONST {
		panic("invalid token type")
	}
	r.printIntegerConstant(depth)
}

func (r *CompilationEngine) printStringConstant(depth string) {
	fmt.Fprintln(r.stream, depth+fmt.Sprintf("<stringConstant> %v </stringConstant>", r.jt.StringVal()))
}

func (r *CompilationEngine) mustPrintStringConstant(depth string) {
	r.mustAdvance()
	if r.jt.TokenType() != jacktokenizer.STRING_CONST {
		panic("invalid token type")
	}
	r.printStringConstant(depth)
}

func (r *CompilationEngine) printIdentifier(depth string) {
	fmt.Fprintln(r.stream, depth+fmt.Sprintf("<identifier> %v </identifier>", r.jt.Identifier()))
}

func (r *CompilationEngine) mustPrintIdentifier(depth string) {
	r.mustAdvance()
	if r.jt.TokenType() != jacktokenizer.IDENTIFIER {
		panic("invalid token type")
	}
	r.printIdentifier(depth)
}

func (r *CompilationEngine) printType(depth string) {
	if r.jt.TokenType() == jacktokenizer.KEYWORD {
		r.printKeyword(depth)
	}
	if r.jt.TokenType() == jacktokenizer.IDENTIFIER {
		r.printIdentifier(depth)
	}
}

func (r *CompilationEngine) mustPrintType(depth string) {
	r.mustAdvance()
	if r.jt.TokenType() != jacktokenizer.KEYWORD && r.jt.TokenType() != jacktokenizer.IDENTIFIER {
		panic("invalid token type")
	}
	r.printType(depth)
}

func (r *CompilationEngine) mustAdvance() {
	if !r.jt.HasMoreTokens() {
		panic("no more tokens")
	}
	r.jt.Advance()
}

func (r *CompilationEngine) PrintTokensXML() {
	fmt.Fprintln(r.stream, "<tokens>")
	for r.jt.HasMoreTokens() {
		r.jt.Advance()
		switch r.jt.TokenType() {
		case jacktokenizer.KEYWORD:
			r.printKeyword("")
		case jacktokenizer.SYMBOL:
			r.printSymbol("")
		case jacktokenizer.INT_CONST:
			r.printIntegerConstant("")
		case jacktokenizer.STRING_CONST:
			r.printStringConstant("")
		case jacktokenizer.IDENTIFIER:
			r.printIdentifier("")
		}
	}
	fmt.Fprintln(r.stream, "</tokens>")
}

func (r *CompilationEngine) Close() {
	r.stream.Flush()
	r.f.Close()
}
