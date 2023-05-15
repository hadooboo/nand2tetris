package compilationegine

import (
	"fmt"

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

	cnt int
}

func NewCompilationEngine(
	ji jacktokenizer.JackTokenizerInterface,
	si symboltable.SymboltableInterface,
	vi vmwriter.VMWriterInterface,
) *CompilationEngine {
	return &CompilationEngine{
		jt:  ji,
		st:  si,
		vw:  vi,
		cnt: 0,
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
	subroutineType := r.mustReturnKeyword()
	_ = r.mustReturnType()
	name := r.mustReturnIdentifier()
	r.st.StartSubroutine(name)
	if subroutineType == "method" {
		r.st.Define("this", r.st.Class(), symboltable.ARG)
	}
	_ = r.mustReturnSymbol()
	r.compileParameterList()
	_ = r.mustReturnSymbol()
	r.compileSubroutineBody(subroutineType)
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

func (r *CompilationEngine) compileSubroutineBody(subroutineType string) {
	_ = r.mustReturnSymbol()
	for r.isVarDec() {
		r.compileVarDec()
	}
	r.vw.WriteFunction(r.st.Class()+"."+r.st.Subroutine(), r.st.VarCount(symboltable.VAR))
	switch subroutineType {
	case "constructor":
		r.vw.WritePush(vmwriter.CONST, r.st.VarCount(symboltable.FIELD))
		r.vw.WriteCall("Memory.alloc", 1)
		r.vw.WritePop(vmwriter.POINTER, 0)
	case "method":
		r.vw.WritePush(vmwriter.ARG, 0)
		r.vw.WritePop(vmwriter.POINTER, 0)
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
	name := r.mustReturnIdentifier()
	if r.isSymbol("[") {
		switch r.st.KindOf(name) {
		case symboltable.STATIC:
			r.vw.WritePush(vmwriter.STATIC, r.st.IndexOf(name))
		case symboltable.FIELD:
			r.vw.WritePush(vmwriter.THIS, r.st.IndexOf(name))
		case symboltable.VAR:
			r.vw.WritePush(vmwriter.LOCAL, r.st.IndexOf(name))
		case symboltable.ARG:
			r.vw.WritePush(vmwriter.ARG, r.st.IndexOf(name))
		}
		_ = r.mustReturnSymbol()
		r.compileExpression()
		_ = r.mustReturnSymbol()
		_ = r.mustReturnSymbol()
		r.compileExpression()
		_ = r.mustReturnSymbol()
		r.vw.WritePop(vmwriter.TEMP, 0)
		r.vw.WriteArithmetic(vmwriter.ADD)
		r.vw.WritePop(vmwriter.POINTER, 1)
		r.vw.WritePush(vmwriter.TEMP, 0)
		r.vw.WritePop(vmwriter.THAT, 0)
	} else {
		_ = r.mustReturnSymbol()
		r.compileExpression()
		_ = r.mustReturnSymbol()
		switch r.st.KindOf(name) {
		case symboltable.STATIC:
			r.vw.WritePop(vmwriter.STATIC, r.st.IndexOf(name))
		case symboltable.FIELD:
			r.vw.WritePop(vmwriter.THIS, r.st.IndexOf(name))
		case symboltable.VAR:
			r.vw.WritePop(vmwriter.LOCAL, r.st.IndexOf(name))
		case symboltable.ARG:
			r.vw.WritePop(vmwriter.ARG, r.st.IndexOf(name))
		}
	}
}

func (r *CompilationEngine) compileIfStatement() {
	r.cnt++
	l1 := fmt.Sprintf("%v.IF.%v.START", r.st.Class(), r.cnt)
	l2 := fmt.Sprintf("%v.IF.%v.END", r.st.Class(), r.cnt)
	_ = r.mustReturnKeyword()
	_ = r.mustReturnSymbol()
	r.compileExpression()
	r.vw.WriteArithmetic(vmwriter.NOT)
	r.vw.WriteIf(l1)
	_ = r.mustReturnSymbol()
	_ = r.mustReturnSymbol()
	r.compileStatements()
	_ = r.mustReturnSymbol()
	r.vw.WriteGoto(l2)
	r.vw.WriteLabel(l1)
	if r.isKeyword(jacktokenizer.ELSE) {
		_ = r.mustReturnKeyword()
		_ = r.mustReturnSymbol()
		r.compileStatements()
		_ = r.mustReturnSymbol()
	}
	r.vw.WriteLabel(l2)
}

func (r *CompilationEngine) compileWhileStatement() {
	r.cnt++
	l1 := fmt.Sprintf("%v.WHILE.%v.START", r.st.Class(), r.cnt)
	l2 := fmt.Sprintf("%v.WHILE.%v.END", r.st.Class(), r.cnt)
	r.vw.WriteLabel(l1)
	_ = r.mustReturnKeyword()
	_ = r.mustReturnSymbol()
	r.compileExpression()
	r.vw.WriteArithmetic(vmwriter.NOT)
	r.vw.WriteIf(l2)
	_ = r.mustReturnSymbol()
	_ = r.mustReturnSymbol()
	r.compileStatements()
	_ = r.mustReturnSymbol()
	r.vw.WriteGoto(l1)
	r.vw.WriteLabel(l2)
}

func (r *CompilationEngine) compileDoStatement() {
	_ = r.mustReturnKeyword()
	name := r.mustReturnIdentifier()
	switch {
	case r.isSymbol("("):
		r.vw.WritePush(vmwriter.POINTER, 0)
		_ = r.mustReturnSymbol()
		n := r.compileExpressionList()
		_ = r.mustReturnSymbol()
		r.vw.WriteCall(r.st.Class()+"."+name, n+1)
	case r.isSymbol("."):
		switch r.st.KindOf(name) {
		case symboltable.STATIC:
			r.vw.WritePush(vmwriter.STATIC, r.st.IndexOf(name))
		case symboltable.FIELD:
			r.vw.WritePush(vmwriter.THIS, r.st.IndexOf(name))
		case symboltable.VAR:
			r.vw.WritePush(vmwriter.LOCAL, r.st.IndexOf(name))
		case symboltable.ARG:
			r.vw.WritePush(vmwriter.ARG, r.st.IndexOf(name))
		}
		_ = r.mustReturnSymbol()
		subroutineName := r.mustReturnIdentifier()
		_ = r.mustReturnSymbol()
		n := r.compileExpressionList()
		_ = r.mustReturnSymbol()
		switch r.st.KindOf(name) {
		case symboltable.NONE:
			r.vw.WriteCall(name+"."+subroutineName, n)
		default:
			r.vw.WriteCall(r.st.TypeOf(name)+"."+subroutineName, n+1)
		}
	}
	_ = r.mustReturnSymbol()
	r.vw.WritePop(vmwriter.TEMP, 0)
}

func (r *CompilationEngine) compileReturnStatement() {
	_ = r.mustReturnKeyword()
	if r.isExpression() {
		r.compileExpression()
	} else {
		r.vw.WritePush(vmwriter.CONST, 0)
	}
	_ = r.mustReturnSymbol()
	r.vw.WriteReturn()
}

func (r *CompilationEngine) compileExpression() {
	r.compileTerm()
	for r.isOp() {
		v := r.mustReturnSymbol()
		r.compileTerm()
		switch v {
		case "+":
			r.vw.WriteArithmetic(vmwriter.ADD)
		case "-":
			r.vw.WriteArithmetic(vmwriter.SUB)
		case "*":
			r.vw.WriteCall("Math.multiply", 2)
		case "/":
			r.vw.WriteCall("Math.divide", 2)
		case "&":
			r.vw.WriteArithmetic(vmwriter.AND)
		case "|":
			r.vw.WriteArithmetic(vmwriter.OR)
		case "<":
			r.vw.WriteArithmetic(vmwriter.LT)
		case ">":
			r.vw.WriteArithmetic(vmwriter.GT)
		case "=":
			r.vw.WriteArithmetic(vmwriter.EQ)
		}
	}
}

func (r *CompilationEngine) compileTerm() {
	switch {
	case r.isIntegerConstant():
		v := r.mustReturnIntegerConstant()
		r.vw.WritePush(vmwriter.CONST, v)
	case r.isStringConstant():
		v := r.mustReturnStringConstant()
		r.vw.WritePush(vmwriter.CONST, len(v))
		r.vw.WriteCall("String.new", 1)
		for _, c := range []byte(v) {
			r.vw.WritePush(vmwriter.CONST, int(c))
			r.vw.WriteCall("String.appendChar", 2)
		}
	case r.isKeywordConstant():
		switch r.mustReturnKeyword() {
		case "true":
			r.vw.WritePush(vmwriter.CONST, 1)
			r.vw.WriteArithmetic(vmwriter.NEG)
		case "false", "null":
			r.vw.WritePush(vmwriter.CONST, 0)
		case "this":
			r.vw.WritePush(vmwriter.POINTER, 0)
		}
	case r.isIdentifier():
		name := r.mustReturnIdentifier()
		switch {
		case r.isSymbol("["):
			switch r.st.KindOf(name) {
			case symboltable.STATIC:
				r.vw.WritePush(vmwriter.STATIC, r.st.IndexOf(name))
			case symboltable.FIELD:
				r.vw.WritePush(vmwriter.THIS, r.st.IndexOf(name))
			case symboltable.VAR:
				r.vw.WritePush(vmwriter.LOCAL, r.st.IndexOf(name))
			case symboltable.ARG:
				r.vw.WritePush(vmwriter.ARG, r.st.IndexOf(name))
			}
			_ = r.mustReturnSymbol()
			r.compileExpression()
			_ = r.mustReturnSymbol()
			r.vw.WriteArithmetic(vmwriter.ADD)
			r.vw.WritePop(vmwriter.POINTER, 1)
			r.vw.WritePush(vmwriter.THAT, 0)
		case r.isSymbol("("):
			r.vw.WritePush(vmwriter.POINTER, 0)
			_ = r.mustReturnSymbol()
			n := r.compileExpressionList()
			_ = r.mustReturnSymbol()
			r.vw.WriteCall(r.st.Class()+"."+name, n+1)
		case r.isSymbol("."):
			switch r.st.KindOf(name) {
			case symboltable.STATIC:
				r.vw.WritePush(vmwriter.STATIC, r.st.IndexOf(name))
			case symboltable.FIELD:
				r.vw.WritePush(vmwriter.THIS, r.st.IndexOf(name))
			case symboltable.VAR:
				r.vw.WritePush(vmwriter.LOCAL, r.st.IndexOf(name))
			case symboltable.ARG:
				r.vw.WritePush(vmwriter.ARG, r.st.IndexOf(name))
			}
			_ = r.mustReturnSymbol()
			subroutineName := r.mustReturnIdentifier()
			_ = r.mustReturnSymbol()
			n := r.compileExpressionList()
			_ = r.mustReturnSymbol()
			switch r.st.KindOf(name) {
			case symboltable.NONE:
				r.vw.WriteCall(name+"."+subroutineName, n)
			default:
				r.vw.WriteCall(r.st.TypeOf(name)+"."+subroutineName, n+1)
			}
		default:
			switch r.st.KindOf(name) {
			case symboltable.STATIC:
				r.vw.WritePush(vmwriter.STATIC, r.st.IndexOf(name))
			case symboltable.FIELD:
				r.vw.WritePush(vmwriter.THIS, r.st.IndexOf(name))
			case symboltable.VAR:
				r.vw.WritePush(vmwriter.LOCAL, r.st.IndexOf(name))
			case symboltable.ARG:
				r.vw.WritePush(vmwriter.ARG, r.st.IndexOf(name))
			}
		}
	case r.isSymbol("("):
		_ = r.mustReturnSymbol()
		r.compileExpression()
		_ = r.mustReturnSymbol()
	case r.isUnaryOp():
		v := r.mustReturnSymbol()
		r.compileTerm()
		switch v {
		case "-":
			r.vw.WriteArithmetic(vmwriter.NEG)
		case "~":
			r.vw.WriteArithmetic(vmwriter.NOT)
		}
	}
}

func (r *CompilationEngine) compileExpressionList() (n int) {
	n = 0
	if r.isExpression() {
		r.compileExpression()
		n++
		for r.isSymbol(",") {
			_ = r.mustReturnSymbol()
			r.compileExpression()
			n++
		}
	}
	return n
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
