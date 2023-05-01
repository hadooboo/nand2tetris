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

	var writeFilePath string
	if fileInfo.IsDir() {
		writeFilePath = filepath.Join(filepath.Clean(os.Args[1]), fileInfo.Name()+".asm")
	} else {
		writeFilePath = filepath.Join(filepath.Dir(os.Args[1]), strings.TrimSuffix(fileInfo.Name(), ".vm")+".asm")
	}

	readFilePaths := flattenFilePaths(os.Args[1])
	c := initCodewriter(writeFilePath, isBootstrapNeeded(readFilePaths))
	for _, readFilePath := range readFilePaths {
		c.SetFileName(strings.TrimSuffix(filepath.Base(readFilePath), ".vm"))
		p := initParser(readFilePath)
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
			case parser.C_FUNCTION:
				c.WriteFunction(p.Arg1(), p.Arg2())
			case parser.C_CALL:
				c.WriteCall(p.Arg1(), p.Arg2())
			case parser.C_RETURN:
				c.WriteReturn()
			}
		}
		p.Close()
	}
	c.Close()
}

func flattenFilePaths(src string) []string {
	res := make([]string, 0)
	filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".vm") {
			res = append(res, path)
		}
		return nil
	})
	return res
}

func isBootstrapNeeded(readFilePaths []string) bool {
	for _, readFilePath := range readFilePaths {
		if filepath.Base(readFilePath) == "Sys.vm" {
			return true
		}
	}
	return false
}

func initParser(readFilePath string) parser.ParserInterface {
	return parser.NewParser(readFilePath)
}

func initCodewriter(writeFilePath string, isBootstrapNeeded bool) codewriter.CodewriterInterface {
	return codewriter.NewCodewriter(writeFilePath, isBootstrapNeeded)
}
