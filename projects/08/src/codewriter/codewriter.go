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
	WriteLabel(string)
	WriteGoto(string)
	WriteIf(string)
	WriteFunction(string, int)
	WriteCall(string, int)
	WriteReturn()
	Close()
}

var _ = CodewriterInterface(&Codewriter{})

type Codewriter struct {
	stream *bufio.Writer
	f      *os.File

	cnt      int
	fileName string
	function string
}

func NewCodewriter(writeFilePath string, isBootstrapNeeded bool) *Codewriter {
	f, err := os.Create(writeFilePath)
	if err != nil {
		panic(err)
	}

	codewriter := &Codewriter{
		stream: bufio.NewWriter(f),
		f:      f,
	}

	if isBootstrapNeeded {
		codewriter.writeBootstrap()
	}

	return codewriter
}

func (r *Codewriter) writeBootstrap() {
	res := strings.TrimSpace(`
@256
D=A
@SP
M=D
	`)
	fmt.Fprintln(r.stream, res)
	r.WriteCall("Sys.init", 0)
}

func (r *Codewriter) SetFileName(fileName string) {
	r.fileName = fileName
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

func (r *Codewriter) WriteLabel(label string) {
	if len(r.function) > 0 {
		label = fmt.Sprintf("%v$%v", r.function, label)
	}
	res := fmt.Sprintf("(%v)", label)
	fmt.Fprintln(r.stream, res)
}

func (r *Codewriter) WriteGoto(label string) {
	if len(r.function) > 0 {
		label = fmt.Sprintf("%v$%v", r.function, label)
	}
	res := fmt.Sprintf("@%v\n0;JMP", label)
	fmt.Fprintln(r.stream, res)
}

func (r *Codewriter) WriteIf(label string) {
	if len(r.function) > 0 {
		label = fmt.Sprintf("%v$%v", r.function, label)
	}
	base := strings.TrimSpace(`
@SP
AM=M-1
D=M
@%v
D;JNE
	`)
	res := fmt.Sprintf(base, label)
	fmt.Fprintln(r.stream, res)
}

func (r *Codewriter) WriteFunction(name string, nVars int) {
	r.function = name
	pushZero := strings.TrimSpace(`
@SP
A=M
M=0
@SP
M=M+1
	`)
	var base, res string
	if nVars == 0 {
		base = strings.TrimSpace(`
(%v)
		`)
		res = fmt.Sprintf(base, r.function)
	} else {
		base = strings.TrimSpace(`
(%v)
%v
		`)
		res = fmt.Sprintf(base, r.function, strings.TrimSpace(strings.Repeat(pushZero+"\n", nVars)))
	}
	fmt.Fprintln(r.stream, res)
}

func (r *Codewriter) WriteCall(name string, nArgs int) {
	r.cnt++
	label := fmt.Sprintf("%v$ret.%v", name, r.cnt)
	base := strings.TrimSpace(`
@%v
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
D=M
@%v
D=D-A
@ARG
M=D
@SP
D=M
@LCL
M=D
@%v
0;JMP
(%v)
	`)
	res := fmt.Sprintf(base, label, 5+nArgs, name, label)
	fmt.Fprintln(r.stream, res)
}

func (r *Codewriter) WriteReturn() {
	res := strings.TrimSpace(`
@LCL
D=M
@R13
M=D
@5
A=D-A
D=M
@R14
M=D
@SP
AM=M-1
D=M
@ARG
A=M
M=D
@ARG
D=M
@SP
M=D+1
@R13
AM=M-1
D=M
@THAT
M=D
@R13
AM=M-1
D=M
@THIS
M=D
@R13
AM=M-1
D=M
@ARG
M=D
@R13
AM=M-1
D=M
@LCL
M=D
@R14
A=M
0;JMP
	`)
	fmt.Fprintln(r.stream, res)
}

func (r *Codewriter) Close() {
	r.stream.Flush()
	r.f.Close()
}
