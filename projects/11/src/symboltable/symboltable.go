package symboltable

type SymboltableInterface interface {
	StartClass(string)
	StartSubroutine(string)
	Class() string
	Subroutine() string
	Define(string, string, SymbolKind)
	VarCount(SymbolKind) int
	KindOf(string) SymbolKind
	TypeOf(string) string
	IndexOf(string) int
}

type SymbolKind int

const (
	NONE SymbolKind = iota
	STATIC
	FIELD
	ARG
	VAR
)

type symbolInfo struct {
	Kind  SymbolKind
	Type  string
	Index int
}

var _ = SymboltableInterface(&Symboltable{})

type Symboltable struct {
	classSymbols      map[string]symbolInfo
	class             string
	staticCount       int
	fieldCount        int
	subroutineSymbols map[string]symbolInfo
	subroutine        string
	argCount          int
	varCount          int
}

func NewSymboltable() *Symboltable {
	return &Symboltable{}
}

func (r *Symboltable) StartClass(class string) {
	r.classSymbols = make(map[string]symbolInfo)
	r.class = class
	r.staticCount = 0
	r.fieldCount = 0
}

func (r *Symboltable) StartSubroutine(subroutine string) {
	r.subroutineSymbols = make(map[string]symbolInfo)
	r.subroutine = subroutine
	r.argCount = 0
	r.varCount = 0
}

func (r *Symboltable) Class() string {
	return r.class
}

func (r *Symboltable) Subroutine() string {
	return r.subroutine
}

func (r *Symboltable) Define(n string, t string, k SymbolKind) {
	switch k {
	case STATIC:
		r.classSymbols[n] = symbolInfo{STATIC, t, r.staticCount}
		r.staticCount++
	case FIELD:
		r.classSymbols[n] = symbolInfo{FIELD, t, r.fieldCount}
		r.fieldCount++
	case ARG:
		r.subroutineSymbols[n] = symbolInfo{ARG, t, r.argCount}
		r.argCount++
	case VAR:
		r.subroutineSymbols[n] = symbolInfo{VAR, t, r.varCount}
		r.varCount++
	}
}

func (r *Symboltable) VarCount(k SymbolKind) int {
	switch k {
	case STATIC:
		return r.staticCount
	case FIELD:
		return r.fieldCount
	case ARG:
		return r.argCount
	case VAR:
		return r.varCount
	}
	return 0
}

func (r *Symboltable) KindOf(name string) SymbolKind {
	if v, ok := r.subroutineSymbols[name]; ok {
		return v.Kind
	}
	if v, ok := r.classSymbols[name]; ok {
		return v.Kind
	}
	return NONE
}

func (r *Symboltable) TypeOf(name string) string {
	if v, ok := r.subroutineSymbols[name]; ok {
		return v.Type
	}
	if v, ok := r.classSymbols[name]; ok {
		return v.Type
	}
	return ""
}

func (r *Symboltable) IndexOf(name string) int {
	if v, ok := r.subroutineSymbols[name]; ok {
		return v.Index
	}
	if v, ok := r.classSymbols[name]; ok {
		return v.Index
	}
	return 0
}
