package code

type CodeInterface interface {
	Dest(string) [2]byte
	Comp(string) [2]byte
	Jump(string) [2]byte
}

var _ = CodeInterface(&Code{})

type Code struct {
}

func NewCode() *Code {
	return &Code{}
}

func (r *Code) Dest(s string) [2]byte {
	res := [2]byte{0, 0}
	switch s {
	case "null0":
		res[1] |= 0b00000000
	case "M":
		res[1] |= 0b00001000
	case "D":
		res[1] |= 0b00010000
	case "MD":
		res[1] |= 0b00011000
	case "A":
		res[1] |= 0b00100000
	case "AM":
		res[1] |= 0b00101000
	case "AD":
		res[1] |= 0b00110000
	case "AMD":
		res[1] |= 0b00111000
	}
	return res
}

func (r *Code) Comp(s string) [2]byte {
	res := [2]byte{0, 0}
	switch s {
	case "0":
		res[0] |= 0b00001010
		res[1] |= 0b10000000
	case "1":
		res[0] |= 0b00001111
		res[1] |= 0b11000000
	case "-1":
		res[0] |= 0b00001110
		res[1] |= 0b10000000
	case "D":
		res[0] |= 0b00000011
		res[1] |= 0b00000000
	case "A":
		res[0] |= 0b00001100
		res[1] |= 0b00000000
	case "!D":
		res[0] |= 0b00000011
		res[1] |= 0b01000000
	case "!A":
		res[0] |= 0b00001100
		res[1] |= 0b11000000
	case "-D":
		res[0] |= 0b00000011
		res[1] |= 0b11000000
	case "-A":
		res[0] |= 0b00001100
		res[1] |= 0b11000000
	case "D+1":
		res[0] |= 0b00000111
		res[1] |= 0b11000000
	case "A+1":
		res[0] |= 0b00001101
		res[1] |= 0b11000000
	case "D-1":
		res[0] |= 0b00000011
		res[1] |= 0b10000000
	case "A-1":
		res[0] |= 0b00001100
		res[1] |= 0b10000000
	case "D+A":
		res[0] |= 0b00000000
		res[1] |= 0b10000000
	case "D-A":
		res[0] |= 0b00000100
		res[1] |= 0b11000000
	case "A-D":
		res[0] |= 0b00000001
		res[1] |= 0b11000000
	case "D&A":
		res[0] |= 0b00000000
		res[1] |= 0b00000000
	case "D|A":
		res[0] |= 0b00000101
		res[1] |= 0b01000000
	case "M":
		res[0] |= 0b00011100
		res[1] |= 0b00000000
	case "!M":
		res[0] |= 0b00011100
		res[1] |= 0b01000000
	case "-M":
		res[0] |= 0b00011100
		res[1] |= 0b11000000
	case "M+1":
		res[0] |= 0b00011101
		res[1] |= 0b11000000
	case "M-1":
		res[0] |= 0b00011100
		res[1] |= 0b10000000
	case "D+M":
		res[0] |= 0b00010000
		res[1] |= 0b10000000
	case "D-M":
		res[0] |= 0b00010100
		res[1] |= 0b11000000
	case "M-D":
		res[0] |= 0b00010001
		res[1] |= 0b11000000
	case "D&M":
		res[0] |= 0b00010000
		res[1] |= 0b00000000
	case "D|M":
		res[0] |= 0b00010101
		res[1] |= 0b01000000
	}
	return res
}

func (r *Code) Jump(s string) [2]byte {
	res := [2]byte{0, 0}
	switch s {
	case "null":
		res[1] |= 0b00000000
	case "JGT":
		res[1] |= 0b00000001
	case "JEQ":
		res[1] |= 0b00000010
	case "JGE":
		res[1] |= 0b00000011
	case "JLT":
		res[1] |= 0b00000100
	case "JNE":
		res[1] |= 0b00000101
	case "JLE":
		res[1] |= 0b00000110
	case "JMP":
		res[1] |= 0b00000111
	}
	return res
}
