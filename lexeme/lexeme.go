// Package lexeme provides a concrete implementation of parlex.Lexeme.
package lexeme

import (
	"fmt"
	"github.com/adamcolton/parlex"
)

type symbol string

func (s symbol) String() string { return string(s) }

// Lexeme is a concrete implementation of parlex.Lexeme
type Lexeme struct {
	K    parlex.Symbol
	V    string
	L, C int
}

// New returns a new Lexeme. Line is initially set to -1 to indicate the the
// position has not been set.
func New(kind parlex.Symbol) *Lexeme {
	return &Lexeme{
		K: kind,
		L: -1,
	}
}

// String is a helper that takes a string and returns a Lexeme with whose kind
// will be the given string.
func String(str string) *Lexeme {
	return &Lexeme{
		K: symbol(str),
		L: -1,
	}
}

// Set the value and returns the Lexeme, it's intended to be used right after a
// call to New
//   noVal := lexeme.New("E")
//   valOf := lexeme.New("int").Val("12")
func (l *Lexeme) Set(val string) *Lexeme {
	l.V = val
	return l
}

// At sets the line and column and returns the Lexeme, it's intended to be used
// right after a call to New
//   noPos := lexeme.New("E")
//   atPos := lexeme.New("E").At(1,2)
func (l *Lexeme) At(line, col int) *Lexeme {
	l.L, l.C = line, col
	return l
}

// Copy a parlex.Lexeme to *Lexeme
func Copy(l parlex.Lexeme) *Lexeme {
	return New(l.Kind()).Set(l.Value()).At(l.Pos())
}

// Kind returns the token indicating what kind of lexeme this is
func (l *Lexeme) Kind() parlex.Symbol { return l.K }

// Value returns the string segment
func (l *Lexeme) Value() string { return l.V }

// Pos returns the position as (line, column) of where the lexeme started in
// the original string.
func (l *Lexeme) Pos() (int, int) { return l.L, l.C }

// String returns a formatted representation of the lexeme.
func (l *Lexeme) String() string {
	pos := ""
	if l.L != -1 {
		pos = fmt.Sprintf(" (%d, %d)", l.L, l.C)
	}
	if l.V == "" {
		return fmt.Sprintf("Lexeme{%s%s}", l.K, pos)
	}
	return fmt.Sprintf("Lexeme{%s:%q%s}", l.K, l.V, pos)
}
