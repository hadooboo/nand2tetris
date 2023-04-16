package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"jaehonam.com/nand2tetris/project/07/codewriter"
	"jaehonam.com/nand2tetris/project/07/parser"
)

func main() {
	c := initCodewriter()
	for _, fileName := range flattenFileNames(os.Args[1]) {
		p := initParser(fileName)
		c.SetFileName(strings.TrimSuffix(fileName, ".vm") + ".asm")
		for p.HasMoreCommands() {
			p.Advance()
			switch p.CommandType() {
			case parser.C_ARITHMETIC:
				c.WriteArithmetic(p.Arg1())
			case parser.C_PUSH:
				c.WritePushPop("push", p.Arg1(), p.Arg2())
			case parser.C_POP:
				c.WritePushPop("pop", p.Arg1(), p.Arg2())
			}
		}
		p.Close()
		c.Close()
	}
}

func flattenFileNames(src string) []string {
	res := make([]string, 0)
	filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".vm") {
			res = append(res, path)
		}
		return nil
	})
	return res
}

func initParser(fileName string) parser.ParserInterface {
	return parser.NewParser(fileName)
}

func initCodewriter() codewriter.CodewriterInterface {
	return codewriter.NewCodewriter()
}
