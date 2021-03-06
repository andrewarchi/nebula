// Package ws parses Whitespace source files.
//
package ws // import "github.com/andrewarchi/nebula/ws"

import (
	"fmt"
	"go/token"
	"math/big"
	"strings"
)

// Token is a lexical token in Whitespace.
type Token struct {
	Type      Type
	Arg       *big.Int
	ArgString string    // Label string, if exists
	Pos       token.Pos // Start position in source
	End       token.Pos // End position in source (exclusive)
}

func (tok *Token) String() string {
	switch {
	case tok.Type == Label:
		return tok.formatArg()
	case tok.Type.HasArg():
		return fmt.Sprintf("%s %s", tok.Type, tok.formatArg())
	default:
		return tok.Type.String()
	}
}

func (tok *Token) formatArg() string {
	if !tok.Type.IsControl() {
		return tok.Arg.String()
	}
	if tok.ArgString != "" {
		return tok.ArgString
	}
	return fmt.Sprintf("label_%s", tok.Arg)
}

// StringWS formats a token as Whitespace.
func (tok *Token) StringWS() string {
	s := tok.Type.StringWS()
	if tok.Type.HasArg() {
		s += tok.formatArgWS()
	}
	return s
}

func (tok *Token) formatArgWS() string {
	var b strings.Builder
	num := tok.Arg
	if !tok.Type.IsControl() {
		if num.Sign() != -1 {
			b.WriteByte(' ')
		} else {
			b.WriteByte('\t')
		}
	}
	if num.Sign() == -1 {
		num = new(big.Int).Neg(num)
	}
	for i := num.BitLen() - 1; i >= 0; i-- {
		if num.Bit(i) == 0 {
			b.WriteByte(' ')
		} else {
			b.WriteByte('\t')
		}
	}
	b.WriteByte('\n')
	return b.String()
}

// Type is the instruction type of a Whitespace token.
type Type uint8

// Instruction types.
const (
	Illegal Type = iota

	// Stack manipulation instructions
	Push
	Dup
	Copy
	Swap
	Drop
	Slide
	Shuffle // non-standard; from Harold Lee's whitespace-0.4

	// Arithmetic instructions
	Add
	Sub
	Mul
	Div
	Mod

	// Heap access instructions
	Store
	Retrieve

	// Control flow instructions
	Label
	Call
	Jmp
	Jz
	Jn
	Ret
	End

	// I/O instructions
	Printc
	Printi
	Readc
	Readi

	// Debug instructions (non-standard)
	Trace     // from Phillip Bradbury's pywhitespace
	DumpStack // from Oliver Burghard's interpreter
	DumpHeap  // from Oliver Burghard's interpreter
)

// IsStack returns true for tokens corresponding to stack manipulation instructions.
func (typ Type) IsStack() bool { return Push <= typ && typ <= Shuffle }

// IsArith returns true for tokens corresponding to arithmetic instructions.
func (typ Type) IsArith() bool { return Add <= typ && typ <= Mod }

// IsHeap returns true for tokens corresponding to heap access instructions.
func (typ Type) IsHeap() bool { return typ == Store || typ == Retrieve }

// IsControl returns true for tokens corresponding to control flow instructions.
func (typ Type) IsControl() bool { return Label <= typ && typ <= End }

// IsIO returns true for tokens corresponding to I/O instructions.
func (typ Type) IsIO() bool { return Printc <= typ && typ <= Readi }

// IsDebug returns true for tokens corresponding to debug instructions.
func (typ Type) IsDebug() bool { return Trace <= typ && typ <= DumpHeap }

// HasArg returns true for instructions that require an argument.
func (typ Type) HasArg() bool {
	switch typ {
	case Push, Copy, Slide, Label, Call, Jmp, Jz, Jn:
		return true
	}
	return false
}

func (typ Type) String() string {
	switch typ {
	case Push:
		return "push"
	case Dup:
		return "dup"
	case Copy:
		return "copy"
	case Swap:
		return "swap"
	case Drop:
		return "drop"
	case Slide:
		return "slide"
	case Shuffle:
		return "shuffle"
	case Add:
		return "add"
	case Sub:
		return "sub"
	case Mul:
		return "mul"
	case Div:
		return "div"
	case Mod:
		return "mod"
	case Store:
		return "store"
	case Retrieve:
		return "retrieve"
	case Label:
		return "label"
	case Call:
		return "call"
	case Jmp:
		return "jmp"
	case Jz:
		return "jz"
	case Jn:
		return "jn"
	case Ret:
		return "ret"
	case End:
		return "end"
	case Printc:
		return "printc"
	case Printi:
		return "printi"
	case Readc:
		return "readc"
	case Readi:
		return "readi"
	case Trace:
		return "trace"
	case DumpStack:
		return "dumpstack"
	case DumpHeap:
		return "dumpheap"
	}
	return fmt.Sprintf("token(%d)", int(typ))
}

// StringWS formats the instruction type as Whitespace syntax.
func (typ Type) StringWS() string {
	switch typ {
	case Push:
		return "  "
	case Dup:
		return " \n "
	case Copy:
		return " \t "
	case Swap:
		return " \n\t"
	case Drop:
		return " \n\n"
	case Slide:
		return " \t\n"
	case Add:
		return "\t   "
	case Sub:
		return "\t  \t"
	case Mul:
		return "\t  \n"
	case Div:
		return "\t \t "
	case Mod:
		return "\t \t\t"
	case Store:
		return "\t\t "
	case Retrieve:
		return "\t\t\t"
	case Label:
		return "\n  "
	case Call:
		return "\n \t"
	case Jmp:
		return "\n \n"
	case Jz:
		return "\n\t "
	case Jn:
		return "\n\t\t"
	case Ret:
		return "\n\t\n"
	case End:
		return "\n\n\n"
	case Printc:
		return "\t\n  "
	case Printi:
		return "\t\n \t"
	case Readc:
		return "\t\n\t "
	case Readi:
		return "\t\n\t\t"
	}
	return fmt.Sprintf("token(%d)", int(typ))
}
