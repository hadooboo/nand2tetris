package vmwriter

import (
	"bufio"
	"fmt"
	"os"
)

type VMWriterInterface interface {
	WritePush(VMSegment, int)
	WritePop(VMSegment, int)
	WriteArithmetic(VMCommand)
	WriteLabel(string)
	WriteGoto(string)
	WriteIf(string)
	WriteCall(string, int)
	WriteFunction(string, int)
	WriteReturn()
	Close()
}

type VMSegment string

const (
	CONST   VMSegment = "constant"
	ARG     VMSegment = "argument"
	LOCAL   VMSegment = "local"
	STATIC  VMSegment = "static"
	THIS    VMSegment = "this"
	THAT    VMSegment = "that"
	POINTER VMSegment = "pointer"
	TEMP    VMSegment = "temp"
)

type VMCommand string

const (
	ADD VMCommand = "add"
	SUB VMCommand = "sub"
	NEG VMCommand = "neg"
	EQ  VMCommand = "eq"
	GT  VMCommand = "gt"
	LT  VMCommand = "lt"
	AND VMCommand = "and"
	OR  VMCommand = "or"
	NOT VMCommand = "not"
)

var _ = VMWriterInterface(&VMWriter{})

type VMWriter struct {
	stream *bufio.Writer
	f      *os.File
}

func NewVMWriter(writeFilePath string) *VMWriter {
	f, err := os.Create(writeFilePath)
	if err != nil {
		panic(err)
	}

	return &VMWriter{
		stream: bufio.NewWriter(f),
		f:      f,
	}
}

func (r *VMWriter) WritePush(segment VMSegment, index int) {
	fmt.Fprintf(r.stream, "push %v %v\n", segment, index)
}

func (r *VMWriter) WritePop(segment VMSegment, index int) {
	fmt.Fprintf(r.stream, "pop %v %v\n", segment, index)
}

func (r *VMWriter) WriteArithmetic(command VMCommand) {
	fmt.Fprintf(r.stream, "%v\n", command)
}

func (r *VMWriter) WriteLabel(label string) {
	fmt.Fprintf(r.stream, "label %v\n", label)
}

func (r *VMWriter) WriteGoto(label string) {
	fmt.Fprintf(r.stream, "goto %v\n", label)
}

func (r *VMWriter) WriteIf(label string) {
	fmt.Fprintf(r.stream, "if-goto %v\n", label)
}

func (r *VMWriter) WriteCall(name string, nArgs int) {
	fmt.Fprintf(r.stream, "call %v %v\n", name, nArgs)
}

func (r *VMWriter) WriteFunction(name string, nLocals int) {
	fmt.Fprintf(r.stream, "function %v %v\n", name, nLocals)
}

func (r *VMWriter) WriteReturn() {
	fmt.Fprintf(r.stream, "return\n")
}

func (r *VMWriter) Close() {
	r.stream.Flush()
	r.f.Close()
}
