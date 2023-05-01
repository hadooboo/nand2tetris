package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"jaehonam.com/nand2tetris/project/10/compilationegine"
	"jaehonam.com/nand2tetris/project/10/jacktokenizer"
)

func main() {
	for _, readFilePath := range flattenFilePaths(os.Args[1]) {
		jt := initJackTokenizer(readFilePath)
		writeFilePath := filepath.Join(filepath.Dir(readFilePath), strings.ReplaceAll(filepath.Base(readFilePath), ".jack", "T.out.xml"))
		ce := initCompilationEngine(jt, writeFilePath)
		ce.PrintTokensXML()
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

func initJackTokenizer(readFilePath string) jacktokenizer.JackTokenizerInterface {
	return jacktokenizer.NewJackTokenizer(readFilePath)
}

func initCompilationEngine(jt jacktokenizer.JackTokenizerInterface, writeFilePath string) compilationegine.CompilationEngineInterface {
	return compilationegine.NewCompilationEngine(jt, writeFilePath)
}
