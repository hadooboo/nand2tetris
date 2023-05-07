package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"jaehonam.com/nand2tetris/project/10/compilationegine"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: ./jackanalyzer [filepath to compile] [tokens|all]")
	}
	var fileSuffix string
	if os.Args[2] == "tokens" {
		fileSuffix = "T.out.xml"
	} else if os.Args[2] == "all" {
		fileSuffix = ".out.xml"
	} else {
		log.Fatalf("invalid operation: %v", os.Args[2])
	}

	for _, readFilePath := range flattenFilePaths(os.Args[1]) {
		writeFilePath := strings.TrimSuffix(readFilePath, ".jack") + fileSuffix
		ce := initCompilationEngine(readFilePath, writeFilePath)
		if os.Args[2] == "tokens" {
			ce.PrintTokensXML()
		} else if os.Args[2] == "all" {
			ce.PrintXML()
		}
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

func initCompilationEngine(readFilePath, writeFilePath string) compilationegine.CompilationEngineInterface {
	return compilationegine.NewCompilationEngine(readFilePath, writeFilePath)
}
