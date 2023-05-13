package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"jaehonam.com/nand2tetris/project/11/compilationegine"
	"jaehonam.com/nand2tetris/project/11/jacktokenizer"
	"jaehonam.com/nand2tetris/project/11/symboltable"
	"jaehonam.com/nand2tetris/project/11/vmwriter"
)

func main() {
	if len(os.Args) < 1 {
		log.Fatalf("usage: ./jackcompiler [filepath to compile]")
	}

	for _, readFilePath := range flattenFilePaths(os.Args[1]) {
		writeFilePath := strings.TrimSuffix(readFilePath, ".jack") + ".vm"
		jt := initJacktokenizer(readFilePath)
		st := initSymboltable()
		vw := initVMWriter(writeFilePath)
		ce := initCompilationEngine(jt, st, vw)
		ce.Compile()
		ce.Close()
	}
}

func flattenFilePaths(src string) []string {
	res := make([]string, 0)
	filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".jack") {
			res = append(res, path)
		}
		return nil
	})
	return res
}

func initJacktokenizer(readFilePath string) jacktokenizer.JackTokenizerInterface {
	return jacktokenizer.NewJackTokenizer(readFilePath)
}

func initSymboltable() symboltable.SymboltableInterface {
	return symboltable.NewSymboltable()
}

func initVMWriter(writeFilePath string) vmwriter.VMWriterInterface {
	return vmwriter.NewVMWriter(writeFilePath)
}

func initCompilationEngine(
	ji jacktokenizer.JackTokenizerInterface,
	si symboltable.SymboltableInterface,
	vi vmwriter.VMWriterInterface,
) compilationegine.CompilationEngineInterface {
	return compilationegine.NewCompilationEngine(ji, si, vi)
}
