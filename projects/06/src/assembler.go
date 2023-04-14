package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"jaehonam.com/nand2tetris/project/06/code"
	"jaehonam.com/nand2tetris/project/06/parser"
	"jaehonam.com/nand2tetris/project/06/symboltable"
)

func main() {
	f1, f2, g := openTwoAsmFileAndCreateOneHackFile(os.Args[1])
	defer f1.Close()
	defer f2.Close()
	defer g.Close()

	r1 := bufio.NewReader(f1)
	r2 := bufio.NewReader(f2)
	w := bufio.NewWriter(g)
	defer w.Flush()

	p1 := initParser(r1)
	p2 := initParser(r2)
	c := initCode()
	s := initSymboltable()

	address := 0
	for p1.HasMoreCommands() {
		p1.Advance()
		switch p1.CommandType() {
		case parser.A_COMMAND, parser.C_COMMAND:
			address++
		case parser.L_COMMAND:
			s.AddEntry(p1.Symbol(), address)
		}
	}

	address = 16
	for p2.HasMoreCommands() {
		p2.Advance()
		switch p2.CommandType() {
		case parser.A_COMMAND:
			if p2.Symbol()[0] >= '0' && p2.Symbol()[0] <= '9' {
				v, err := strconv.ParseInt(p2.Symbol(), 10, 64)
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(w, "%016b\n", v)
			} else {
				if !s.Contains(p2.Symbol()) {
					s.AddEntry(p2.Symbol(), address)
					address++
				}
				fmt.Fprintf(w, "%016b\n", s.GetAddress(p2.Symbol()))
			}
		case parser.C_COMMAND:
			v := [2]byte{0b11100000, 0b00000000}
			dest := c.Dest(p2.Dest())
			comp := c.Comp(p2.Comp())
			jump := c.Jump(p2.Jump())
			v[0] |= dest[0] | comp[0] | jump[0]
			v[1] |= dest[1] | comp[1] | jump[1]
			fmt.Fprintf(w, "%08b%08b\n", v[0], v[1])
		case parser.L_COMMAND:
		}
	}
}

func openTwoAsmFileAndCreateOneHackFile(filepath string) (*os.File, *os.File, *os.File) {
	if !strings.HasSuffix(filepath, ".asm") {
		panic("cannot handle none asm file")
	}

	f1, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	f2, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	filename := strings.TrimSuffix(filepath, ".asm") + ".hack"

	g, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	return f1, f2, g
}

func initParser(stream *bufio.Reader) parser.ParserInterface {
	return parser.NewParser(stream)
}

func initCode() code.CodeInterface {
	return code.NewCode()
}

func initSymboltable() symboltable.SymboltableInterface {
	return symboltable.NewSymboltable()
}
