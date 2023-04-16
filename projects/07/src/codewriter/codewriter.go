package codewriter

import (
	"bufio"
	"fmt"
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

func (r *Codewriter) WriteArithmetic(command string) {
	var res string
	switch command {
	case "neg":
		res = `
@SP
AM=M-1
M=-M
@SP
M=M+1
		`
	case "not":
		res = `
@SP
AM=M-1
M=!M
@SP
M=M+1
		`
	case "add":
		res = `
@SP
AM=M-1
D=M
@SP
AM=M-1
M=M+D
@SP
M=M+1
		`
	case "sub":
		res = `
@SP
AM=M-1
D=M
@SP
AM=M-1
M=M-D
@SP
M=M+1
		`
	case "and":
		res = `
@SP
AM=M-1
D=M
@SP
AM=M-1
M=M&D
@SP
M=M+1
		`
	case "or":
		res = `
@SP
AM=M-1
D=M
@SP
AM=M-1
M=M|D
@SP
M=M+1
		`
	case "eq":
		r.cnt++
		res = fmt.Sprintf(`
@SP
AM=M-1
D=M
@SP
AM=M-1
D=M-D
@EQ.TRUE.%v
D;JEQ
@EQ.FALSE.%v
0;JMP
(EQ.TRUE.%v)
@SP
A=M
M=-1
@EQ.CONT.%v
0;JMP
(EQ.FALSE.%v)
@SP
A=M
M=0
@EQ.CONT.%v
0;JMP
(EQ.CONT.%v)
@SP
M=M+1
		`, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt)
	case "gt":
		r.cnt++
		res = fmt.Sprintf(`
@SP
AM=M-1
D=M
@GT.BGE.%v
D;JGE
@GT.BLT.%v
0;JMP
(GT.BGE.%v)
@SP
AM=M-1
M=M-D
D=D+M
@GT.BGE.AGE.%v
D;JGE
@GT.BGE.ALT.%v
0;JMP
(GT.BGE.AGE.%v)
@GT.CMPR.%v
0;JMP
(GT.BGE.ALT.%v)
@GT.FALSE.%v
0;JMP
(GT.BLT.%v)
@SP
AM=M-1
M=M-D
D=D+M
@GT.BLT.AGE.%v
D;JGE
@GT.BLT.ALT.%v
0;JMP
(GT.BLT.AGE.%v)
@GT.TRUE.%v
0;JMP
(GT.BLT.ALT.%v)
@GT.CMPR.%v
0;JMP
(GT.CMPR.%v)
@SP
A=M
D=M
@GT.TRUE.%v
D;JGT
@GT.FALSE.%v
0;JMP
(GT.TRUE.%v)
@SP
A=M
M=-1
@GT.CONT.%v
0;JMP
(GT.FALSE.%v)
@SP
A=M
M=0
@GT.CONT.%v
0;JMP
(GT.CONT.%v)
@SP
M=M+1
		`, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt)
	case "lt":
		r.cnt++
		res = fmt.Sprintf(`
@SP
AM=M-1
D=M
@LT.BGE.%v
D;JGE
@LT.BLT.%v
0;JMP
(LT.BGE.%v)
@SP
AM=M-1
M=M-D
D=D+M
@LT.BGE.AGE.%v
D;JGE
@LT.BGE.ALT.%v
0;JMP
(LT.BGE.AGE.%v)
@LT.CMPR.%v
0;JMP
(LT.BGE.ALT.%v)
@LT.TRUE.%v
0;JMP
(LT.BLT.%v)
@SP
AM=M-1
M=M-D
D=D+M
@LT.BLT.AGE.%v
D;JGE
@LT.BLT.ALT.%v
0;JMP
(LT.BLT.AGE.%v)
@LT.FALSE.%v
0;JMP
(LT.BLT.ALT.%v)
@LT.CMPR.%v
0;JMP
(LT.CMPR.%v)
@SP
A=M
D=M
@LT.TRUE.%v
D;JLT
@LT.FALSE.%v
0;JMP
(LT.TRUE.%v)
@SP
A=M
M=-1
@LT.CONT.%v
0;JMP
(LT.FALSE.%v)
@SP
A=M
M=0
@LT.CONT.%v
0;JMP
(LT.CONT.%v)
@SP
M=M+1
		`, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt, r.cnt)
	}
	fmt.Fprintln(r.stream, strings.TrimSpace(res))
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
