package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"jaehonam.com/nand2tetris/project/08/codewriter"
	"jaehonam.com/nand2tetris/project/08/parser"
)

func main() {
	fileInfo, err := os.Stat(os.Args[1])
	if err != nil {
		panic(err)
	}

	var filePath string
	if fileInfo.IsDir() {
		filePath = filepath.Join(filepath.Clean(os.Args[1]), fileInfo.Name()+".asm")
	} else {
		filePath = filepath.Join(filepath.Dir(os.Args[1]), strings.TrimSuffix(fileInfo.Name(), ".vm")+".asm")
	}

	c := initCodewriter(filePath)
	for _, fileName := range flattenFileNames(os.Args[1]) {
		c.SetFileName(fileName)
		p := initParser(fileName)
		for p.HasMoreCommands() {
			p.Advance()
			switch p.CommandType() {
			case parser.C_ARITHMETIC:
				c.WriteArithmetic(p.Arg1())
			case parser.C_PUSH:
				c.WritePushPop("push", p.Arg1(), p.Arg2())
			case parser.C_POP:
				c.WritePushPop("pop", p.Arg1(), p.Arg2())
			case parser.C_LABEL:
				c.WriteLabel(p.Arg1())
			case parser.C_GOTO:
				c.WriteGoto(p.Arg1())
			case parser.C_IF:
				c.WriteIf(p.Arg1())
			}
		}
		p.Close()
	}
	c.Close()
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

func initCodewriter(filePath string) codewriter.CodewriterInterface {
	return codewriter.NewCodewriter(filePath)
}
