package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"jaehonam.com/nand2tetris/project/06/code"
	"jaehonam.com/nand2tetris/project/06/parser"
)

func main() {
	f, g := openAsmFileAndCreateHackFile(os.Args[1])
	defer f.Close()
	defer g.Close()

	r := bufio.NewReader(f)
	w := bufio.NewWriter(g)
	defer w.Flush()

	p := initParser(r)
	c := initCode()

	for p.HasMoreCommands() {
		p.Advance()
		switch p.CommandType() {
		case parser.A_COMMAND:
			v, err := strconv.ParseInt(p.Symbol(), 10, 64)
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "%016b\n", v)
		case parser.C_COMMAND:
			v := [2]byte{0b11100000, 0b00000000}
			dest := c.Dest(p.Dest())
			comp := c.Comp(p.Comp())
			jump := c.Jump(p.Jump())
			v[0] |= dest[0] | comp[0] | jump[0]
			v[1] |= dest[1] | comp[1] | jump[1]
			fmt.Fprintf(w, "%08b%08b\n", v[0], v[1])
		case parser.L_COMMAND:
		}
	}
}

func openAsmFileAndCreateHackFile(filepath string) (*os.File, *os.File) {
	if !strings.HasSuffix(filepath, ".asm") {
		panic("cannot handle none asm file")
	}

	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	filename := strings.TrimSuffix(filepath, ".asm") + ".hack"

	g, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	return f, g
}

func initParser(stream *bufio.Reader) parser.ParserInterface {
	return parser.NewParser(stream)
}

func initCode() code.CodeInterface {
	return code.NewCode()
}
