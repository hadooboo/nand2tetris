package codewriter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CodewriterInterface interface {
	SetFileName(string)
	WriteArithmetic(string)
	WritePushPop(string, string, int)
	Close()
}

var _ = CodewriterInterface(&Codewriter{})

type Codewriter struct {
	stream *bufio.Writer
	f      *os.File

	cnt      int
	fileName string
}

func NewCodewriter() *Codewriter {
	return &Codewriter{}
}

func (r *Codewriter) SetFileName(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	r.stream = bufio.NewWriter(f)
	r.f = f

	_, r.fileName = filepath.Split(fileName)
	r.fileName = strings.TrimSuffix(r.fileName, ".asm")
}

func (r *Codewriter) WriteArithmetic(command string) {
	var res string
	switch command {
	case "neg", "not":
		res = r.unaryArithmetic(command)
	case "add", "sub", "and", "or":
		res = r.binaryArithmetic(command)
	case "eq", "gt", "lt":
		res = r.comparingArithmetic(command)
	}
	fmt.Fprintln(r.stream, res)
}

func (r *Codewriter) unaryArithmetic(command string) string {
	base := strings.TrimSpace(`
@SP
AM=M-1
M=%vM
@SP
M=M+1
	`)
	return fmt.Sprintf(base, map[string]string{
		"neg": "-",
		"not": "!",
	}[command])
}

func (r *Codewriter) binaryArithmetic(command string) string {
	base := strings.TrimSpace(`
@SP
AM=M-1
D=M
@SP
AM=M-1
M=M%vD
@SP
M=M+1
	`)
	return fmt.Sprintf(base, map[string]string{
		"add": "+",
		"sub": "-",
		"and": "&",
		"or":  "|",
	}[command])
}

func (r *Codewriter) comparingArithmetic(command string) string {
	r.cnt++
	base := strings.TrimSpace(`
@SP
AM=M-1
D=M
@SP
AM=M-1
D=M-D
@%v
D;%v
@SP
A=M
M=0
@%v
0;JMP
(%v)
@SP
A=M
M=-1
(%v)
@SP
M=M+1
	`)
	trueLabel := fmt.Sprintf("%v.TRUE.%v", command, r.cnt)
	contLabel := fmt.Sprintf("%v.CONT.%v", command, r.cnt)
	return fmt.Sprintf(base, trueLabel, map[string]string{
		"eq": "JEQ",
		"gt": "JGT",
		"lt": "JLT",
	}[command], contLabel, trueLabel, contLabel)
}

func (r *Codewriter) WritePushPop(command string, segment string, index int) {
	var res string
	switch command {
	case "push":
		switch segment {
		case "constant":
			res = r.pushConstant(command, segment, index)
		case "local", "argument", "this", "that":
			res = r.pushBase(command, segment, index)
		case "pointer", "temp":
			res = r.pushDirect(command, segment, index)
		case "static":
			res = r.pushStatic(command, segment, index)
		}
	case "pop":
		switch segment {
		case "local", "argument", "this", "that":
			res = r.popBase(command, segment, index)
		case "pointer", "temp":
			res = r.popDirect(command, segment, index)
		case "static":
			res = r.popStatic(command, segment, index)
		}
	}
	fmt.Fprintln(r.stream, res)
}

func (r *Codewriter) pushConstant(command string, segment string, index int) string {
	base := strings.TrimSpace(`
@%v
D=A
@SP
A=M
M=D
@SP
M=M+1
	`)
	return fmt.Sprintf(base, index)
}

func (r *Codewriter) pushBase(command string, segment string, index int) string {
	base := strings.TrimSpace(`
@%v
D=A
@%v
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
	`)
	return fmt.Sprintf(base, index, map[string]string{
		"local":    "LCL",
		"argument": "ARG",
		"this":     "THIS",
		"that":     "THAT",
	}[segment])
}

func (r *Codewriter) pushDirect(command string, segment string, index int) string {
	base := strings.TrimSpace(`
@%v
D=A
@%v
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1
	`)
	return fmt.Sprintf(base, map[string]int{
		"pointer": 3,
		"temp":    5,
	}[segment], index)
}

func (r *Codewriter) pushStatic(command string, segment string, index int) string {
	base := strings.TrimSpace(`
@%v.%v
D=M
@SP
A=M
M=D
@SP
M=M+1
	`)
	return fmt.Sprintf(base, r.fileName, index)
}

func (r *Codewriter) popBase(command string, segment string, index int) string {
	base := strings.TrimSpace(`
@%v
D=A
@%v
D=M+D
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
	`)
	return fmt.Sprintf(base, index, map[string]string{
		"local":    "LCL",
		"argument": "ARG",
		"this":     "THIS",
		"that":     "THAT",
	}[segment])
}

func (r *Codewriter) popDirect(command string, segment string, index int) string {
	base := strings.TrimSpace(`
@%v
D=A
@%v
D=A+D
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
	`)
	return fmt.Sprintf(base, map[string]int{
		"pointer": 3,
		"temp":    5,
	}[segment], index)
}

func (r *Codewriter) popStatic(command string, segment string, index int) string {
	base := strings.TrimSpace(`
@SP
AM=M-1
D=M
@%v.%v
M=D
	`)
	return fmt.Sprintf(base, r.fileName, index)
}

func (r *Codewriter) Close() {
	r.stream.Flush()
	r.f.Close()
}
