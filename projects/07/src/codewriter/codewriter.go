package codewriter

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	cnt int
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
}

func (r *Codewriter) unaryArithmetic(command string) string {
	base := strings.TrimSpace(`
@SP
AM=M-1
M=%vM
@SP
M=M+1
	`)
	var op string
	switch command {
	case "neg":
		op = "-"
	case "not":
		op = "!"
	default:
		log.Panicf("invalid unary arithmetic command: %v", command)
	}
	return fmt.Sprintf(base, op)
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
	var op string
	switch command {
	case "add":
		op = "+"
	case "sub":
		op = "-"
	case "and":
		op = "&"
	case "or":
		op = "|"
	default:
		log.Panicf("invalid binary arithmetic command: %v", command)
	}
	return fmt.Sprintf(base, op)
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
	var jump string
	switch command {
	case "eq":
		jump = "JEQ"
	case "gt":
		jump = "JGT"
	case "lt":
		jump = "JLT"
	default:
		log.Panicf("invalid comparing arithmetic command: %v", command)
	}
	trueLabel := fmt.Sprintf("%v.TRUE.%v", command, r.cnt)
	contLabel := fmt.Sprintf("%v.CONT.%v", command, r.cnt)
	return fmt.Sprintf(base, trueLabel, jump, contLabel, trueLabel, contLabel)
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
	default:
		log.Panicf("invalid arithmetic command: %v", command)
	}
	fmt.Fprintln(r.stream, res)
}

func (r *Codewriter) WritePushPop(command string, segment string, index int) {
	var res string
	switch command {
	case "push":
		switch segment {
		case "constant":
			res = fmt.Sprintf(`
@%v
D=A
@SP
A=M
M=D
@SP
M=M+1
			`, index)
		}
	case "pop":
	}
	fmt.Fprintln(r.stream, strings.TrimSpace(res))
}

func (r *Codewriter) Close() {
	r.stream.Flush()
	r.f.Close()
}
