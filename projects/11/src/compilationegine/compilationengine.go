package compilationegine

import (
	"jaehonam.com/nand2tetris/project/11/jacktokenizer"
	"jaehonam.com/nand2tetris/project/11/symboltable"
	"jaehonam.com/nand2tetris/project/11/vmwriter"
)

type CompilationEngineInterface interface {
	Compile()
	Close()
}

var _ = CompilationEngineInterface(&CompilationEngine{})

type CompilationEngine struct {
	jt jacktokenizer.JackTokenizerInterface
	st symboltable.SymboltableInterface
	vw vmwriter.VMWriterInterface
}

func NewCompilationEngine(
	ji jacktokenizer.JackTokenizerInterface,
	si symboltable.SymboltableInterface,
	vi vmwriter.VMWriterInterface,
) *CompilationEngine {
	return &CompilationEngine{
		jt: ji,
		st: si,
		vw: vi,
	}
}

func (r *CompilationEngine) Compile() {
	r.compileClass()
}

func (r *CompilationEngine) compileClass() {
	_ = r.mustReturnKeyword()
	name := r.mustReturnIdentifier()
	r.st.StartClass(name)
	_ = r.mustReturnSymbol()
	for r.isClassVarDec() {
		r.compileClassVarDec()
	}
	for r.isSubroutineDec() {
		r.compileSubroutineDec()
	}
	_ = r.mustReturnSymbol()
}

func (r *CompilationEngine) compileClassVarDec() {
	var kind symboltable.SymbolKind
	switch r.mustReturnKeyword() {
	case "static":
		kind = symboltable.STATIC
	case "field":
		kind = symboltable.FIELD
	}
	t := r.mustReturnType()
	name := r.mustReturnIdentifier()
	r.st.Define(name, t, kind)
	for r.isSymbol(",") {
		_ = r.mustReturnSymbol()
		name := r.mustReturnIdentifier()
		r.st.Define(name, t, kind)
	}
	_ = r.mustReturnSymbol()
}

func (r *CompilationEngine) compileSubroutineDec() {
	_ = r.mustReturnKeyword()
	_ = r.mustReturnType()
	name := r.mustReturnIdentifier()
	r.st.StartSubroutine(name)
	_ = r.mustReturnSymbol()
	r.compileParameterList()
	_ = r.mustReturnSymbol()
	r.compileSubroutineBody()
}

func (r *CompilationEngine) compileParameterList() {
	if r.isType() {
		t := r.mustReturnType()
		name := r.mustReturnIdentifier()
		r.st.Define(name, t, symboltable.ARG)
		for r.isSymbol(",") {
			_ = r.mustReturnSymbol()
			t := r.mustReturnType()
			name := r.mustReturnIdentifier()
			r.st.Define(name, t, symboltable.ARG)
		}
	}
}

func (r *CompilationEngine) compileSubroutineBody() {
	_ = r.mustReturnSymbol()
	for r.isVarDec() {
		r.compileVarDec()
	}
	r.compileStatements()
	_ = r.mustReturnSymbol()
}

func (r *CompilationEngine) compileVarDec() {
	_ = r.mustReturnKeyword()
	t := r.mustReturnType()
	name := r.mustReturnIdentifier()
	r.st.Define(name, t, symboltable.VAR)
	for r.isSymbol(",") {
		_ = r.mustReturnSymbol()
		name := r.mustReturnIdentifier()
		r.st.Define(name, t, symboltable.VAR)
	}
	_ = r.mustReturnSymbol()
}

func (r *CompilationEngine) compileStatements() {
	for r.isStatement() {
		switch {
		case r.isKeyword(jacktokenizer.LET):
			r.compileLetStatement()
		case r.isKeyword(jacktokenizer.IF):
			r.compileIfStatement()
		case r.isKeyword(jacktokenizer.WHILE):
			r.compileWhileStatement()
		case r.isKeyword(jacktokenizer.DO):
			r.compileDoStatement()
		case r.isKeyword(jacktokenizer.RETURN):
			r.compileReturnStatement()
		}
	}
}

func (r *CompilationEngine) compileLetStatement() {
	_ = r.mustReturnKeyword()
	_ = r.mustReturnIdentifier()
	if r.isSymbol("[") {
		_ = r.mustReturnSymbol()
		r.compileExpression()
		_ = r.mustReturnSymbol()
	}
	_ = r.mustReturnSymbol()
	r.compileExpression()
	_ = r.mustReturnSymbol()
}

func (r *CompilationEngine) compileIfStatement() {
	_ = r.mustReturnKeyword()
	_ = r.mustReturnSymbol()
	r.compileExpression()
	_ = r.mustReturnSymbol()
	_ = r.mustReturnSymbol()
	r.compileStatements()
	_ = r.mustReturnSymbol()
	if r.isKeyword(jacktokenizer.ELSE) {
		_ = r.mustReturnKeyword()
		_ = r.mustReturnSymbol()
		r.compileStatements()
		_ = r.mustReturnSymbol()
	}
}

func (r *CompilationEngine) compileWhileStatement() {
	_ = r.mustReturnKeyword()
	_ = r.mustReturnSymbol()
	r.compileExpression()
	_ = r.mustReturnSymbol()
	_ = r.mustReturnSymbol()
	r.compileStatements()
	_ = r.mustReturnSymbol()
}

func (r *CompilationEngine) compileDoStatement() {
	_ = r.mustReturnKeyword()
	_ = r.mustReturnIdentifier()
	switch {
	case r.isSymbol("("):
		_ = r.mustReturnSymbol()
		r.compileExpressionList()
		_ = r.mustReturnSymbol()
	case r.isSymbol("."):
		_ = r.mustReturnSymbol()
		_ = r.mustReturnIdentifier()
		_ = r.mustReturnSymbol()
		r.compileExpressionList()
		_ = r.mustReturnSymbol()
	}
	_ = r.mustReturnSymbol()
}

func (r *CompilationEngine) compileReturnStatement() {
	_ = r.mustReturnKeyword()
	if r.isExpression() {
		r.compileExpression()
	}
	_ = r.mustReturnSymbol()
}

func (r *CompilationEngine) compileExpression() {
	r.compileTerm()
	for r.isOp() {
		_ = r.mustReturnSymbol()
		r.compileTerm()
	}
}

func (r *CompilationEngine) compileTerm() {
	switch {
	case r.isIntegerConstant():
		_ = r.mustReturnIntegerConstant()
	case r.isStringConstant():
		_ = r.mustReturnStringConstant()
	case r.isKeywordConstant():
		_ = r.mustReturnKeyword()
	case r.isIdentifier():
		_ = r.mustReturnIdentifier()
		switch {
		case r.isSymbol("["):
			_ = r.mustReturnSymbol()
			r.compileExpression()
			_ = r.mustReturnSymbol()
		case r.isSymbol("("):
			_ = r.mustReturnSymbol()
			r.compileExpressionList()
			_ = r.mustReturnSymbol()
		case r.isSymbol("."):
			_ = r.mustReturnSymbol()
			_ = r.mustReturnIdentifier()
			_ = r.mustReturnSymbol()
			r.compileExpressionList()
			_ = r.mustReturnSymbol()
		}
	case r.isSymbol("("):
		_ = r.mustReturnSymbol()
		r.compileExpression()
		_ = r.mustReturnSymbol()
	case r.isUnaryOp():
		_ = r.mustReturnSymbol()
		r.compileTerm()
	}
}

func (r *CompilationEngine) compileExpressionList() {
	if r.isExpression() {
		r.compileExpression()
		for r.isSymbol(",") {
			_ = r.mustReturnSymbol()
			r.compileExpression()
		}
	}
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

func (r *CompilationEngine) mustReturnKeyword() string {
	r.mustAdvance()
	if r.jt.TokenType() == jacktokenizer.KEYWORD {
		return r.jt.Keyword()
	}
	panic("invalid token type")
}

func (r *CompilationEngine) mustReturnSymbol() string {
	r.mustAdvance()
	if r.jt.TokenType() == jacktokenizer.SYMBOL {
		return r.jt.Symbol()
	}
	panic("invalid token type")
}

func (r *CompilationEngine) mustReturnIdentifier() string {
	r.mustAdvance()
	if r.jt.TokenType() == jacktokenizer.IDENTIFIER {
		return r.jt.Identifier()
	}
	panic("invalid token type")
}

func (r *CompilationEngine) mustReturnIntegerConstant() int {
	r.mustAdvance()
	if r.jt.TokenType() == jacktokenizer.INT_CONST {
		return r.jt.IntVal()
	}
	panic("invalid token type")
}

func (r *CompilationEngine) mustReturnStringConstant() string {
	r.mustAdvance()
	if r.jt.TokenType() == jacktokenizer.STRING_CONST {
		return r.jt.StringVal()
	}
	panic("invalid token type")
}

func (r *CompilationEngine) mustReturnType() string {
	r.mustAdvance()
	if r.jt.TokenType() == jacktokenizer.KEYWORD {
		return r.jt.Keyword()
	}
	if r.jt.TokenType() == jacktokenizer.IDENTIFIER {
		return r.jt.Identifier()
	}
	panic("invalid token type")
}

func (r *CompilationEngine) mustAdvance() {
	if !r.jt.HasMoreTokens() {
		panic("no more tokens")
	}
	r.jt.Advance()
}

func (r *CompilationEngine) Close() {
	r.vw.Close()
}
