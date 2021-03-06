// Package bf parses Brainfuck source files.
//
package bf // import "github.com/andrewarchi/nebula/bf"

import (
	"fmt"
	"go/token"
	"strings"
)

// Token is a lexical token in Brainfuck.
type Token struct {
	Type Type
	Pos  token.Pos
}

// Type is the instruction type of a Brainfuck token.
type Type uint8

// Instruction types.
const (
	Illegal    Type = iota
	IncPtr          // >
	DecPtr          // <
	IncData         // +
	DecData         // -
	Print           // .
	Read            // ,
	Bracket         // [
	EndBracket      // ]
)

func (typ Type) String() string {
	switch typ {
	case IncPtr:
		return "incptr"
	case DecPtr:
		return "decptr"
	case IncData:
		return "incdata"
	case DecData:
		return "decdata"
	case Print:
		return "print"
	case Read:
		return "read"
	case Bracket:
		return "bracket"
	case EndBracket:
		return "endbracket"
	}
	return fmt.Sprintf("token(%d)", int(typ))
}

// StringBF formats the instruction type as Brainfuck syntax.
func (typ Type) StringBF() string {
	switch typ {
	case IncPtr:
		return ">"
	case DecPtr:
		return "<"
	case IncData:
		return "+"
	case DecData:
		return "-"
	case Print:
		return "."
	case Read:
		return ","
	case Bracket:
		return "["
	case EndBracket:
		return "]"
	}
	return fmt.Sprintf("token(%d)", int(typ))
}

// Program is a sequence of Brainfuck tokens with file information.
type Program struct {
	Tokens []*Token
	File   *token.File
}

func (p *Program) String() string {
	var b strings.Builder
	indent := "  "
	indentLevel := 0
	var prev Type
	for i, tok := range p.Tokens {
		if tok.Type == EndBracket && indentLevel > 0 {
			indentLevel--
		}
		if (prev != tok.Type && i != 0) || tok.Type == Bracket || tok.Type == EndBracket {
			b.WriteByte('\n')
			b.WriteString(strings.Repeat(indent, indentLevel))
		}
		if tok.Type == Bracket {
			indentLevel++
		}
		b.WriteString(tok.Type.StringBF())
		prev = tok.Type
	}
	return b.String()
}
