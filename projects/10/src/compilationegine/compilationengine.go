package compilationegine

import (
	"bufio"
	"fmt"
	"os"

	"jaehonam.com/nand2tetris/project/10/jacktokenizer"
)

type CompilationEngineInterface interface {
	PrintTokensXML()
	Close()
}

var _ = CompilationEngineInterface(&CompilationEngine{})

type CompilationEngine struct {
	stream *bufio.Writer
	f      *os.File

	jt jacktokenizer.JackTokenizerInterface
}

func NewCompilationEngine(jt jacktokenizer.JackTokenizerInterface, writeFilePath string) *CompilationEngine {
	f, err := os.Create(writeFilePath)
	if err != nil {
		panic(err)
	}

	return &CompilationEngine{
		stream: bufio.NewWriter(f),
		f:      f,
		jt:     jt,
	}
}

func (r *CompilationEngine) PrintTokensXML() {
	fmt.Fprintln(r.stream, "<tokens>")
	for r.jt.HasMoreTokens() {
		r.jt.Advance()
		switch r.jt.TokenType() {
		case jacktokenizer.KEYWORD:
			fmt.Fprintf(r.stream, "<keyword> %s </keyword>\n", r.jt.Keyword())
		case jacktokenizer.SYMBOL:
			fmt.Fprintf(r.stream, "<symbol> %s </symbol>\n", r.jt.Symbol())
		case jacktokenizer.INT_CONST:
			fmt.Fprintf(r.stream, "<integerConstant> %d </integerConstant>\n", r.jt.IntVal())
		case jacktokenizer.STRING_CONST:
			fmt.Fprintf(r.stream, "<stringConstant> %s </stringConstant>\n", r.jt.StringVal())
		case jacktokenizer.IDENTIFIER:
			fmt.Fprintf(r.stream, "<identifier> %s </identifier>\n", r.jt.Identifier())
		}
	}
	fmt.Fprintln(r.stream, "</tokens>")
}

func (r *CompilationEngine) Close() {
	r.stream.Flush()
	r.f.Close()
}
