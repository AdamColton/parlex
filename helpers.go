package parlex

import (
	"fmt"
	"strings"
)

// Stringer is a copy of the ubiquitous stringer interface
type Stringer interface {
	String() string
}

// LexemeString returns a Lexeme in the form "Kind : Value"
func LexemeString(ls ...Lexeme) string {
	if len(ls) == 0 {
		return ""
	}
	var strs = make([]string, len(ls))
	for i, l := range ls {
		k, v := l.Kind().String(), l.Value()
		if v == "" {
			strs[i] = k
		} else {
			strs[i] = k + ": " + v
		}
	}
	if len(strs) == 1 {
		return strs[0]
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

// LexemeList prints a slice of Lexemes, one per line. Useful when debugging a
// lexer.
func LexemeList(ls []Lexeme) string {
	if len(ls) == 0 {
		return ""
	}
	var strs = make([]string, len(ls))
	for i, l := range ls {
		k, v := l.Kind().String(), l.Value()
		pl, pc := l.Pos()
		if pl > 0 || pc > 0 {
			strs[i] = fmt.Sprintf("%s: %q (%d, %d)", k, v, pl, pc)
		} else {
			strs[i] = fmt.Sprintf("%s: %q", k, v)
		}
	}
	return strings.Join(strs, "\n")
}

// MustParser consumes the error from a parser constructor and panics if it is
// not nil.
func MustParser(p Parser, err error) Parser {
	if err != nil {
		panic(err)
	}
	return p
}

// MustLexer consumes the error from a lexer constructor and panics if it is
// not nil.
func MustLexer(l Lexer, err error) Lexer {
	if err != nil {
		panic(err)
	}
	return l
}

// MustGrammar consumes the error from a grammar constructor and panics if it is
// not nil.
func MustGrammar(g Grammar, err error) Grammar {
	if err != nil {
		panic(err)
	}
	return g
}

// LexError allows a Lexer to try to continue lexing when it encounters an error
// and return useful error values.
type LexError interface {
	Lexeme
	Error() string
}

type lexErr struct {
	Lexeme
}

func (l *lexErr) Error() string {
	if pl, pc := l.Pos(); pl > 0 {
		return fmt.Sprintf("%s %d:%d) %s", l.Kind().String(), pl, pc, l.Value())
	}
	return fmt.Sprintf("%s) %s", l.Kind().String(), l.Value())

}

// LexErrors takes a slice of Lexemes and returns any that are instances of
// LexError. This allows multiple errors to be caught.
func LexErrors(lexemes []Lexeme) (errs []LexError) {
	for _, lx := range lexemes {
		if err, ok := lx.(LexError); ok {
			errs = append(errs, err)
		}
	}
	return
}

// SymLen returns the rune length of a Symbol
func SymLen(s Symbol) int {
	return len([]rune(s.String()))
}

// ProductionIterator is used to iterate over a production.
type ProductionIterator struct {
	Symbol
	Production Production
	ln         int
	Idx        int
}

// Next moves Symbol to the next symbol in the production. It returns false if
// there are no more symbols.
func (p *ProductionIterator) Next() bool {
	if p.ln == 0 {
		if p.Production == nil {
			return false
		}
		p.ln = p.Production.Symbols()
		if p.ln == 0 {
			return false
		}
	} else {
		p.Idx++
	}
	if p.Idx >= p.ln {
		return false
	}
	p.Symbol = p.Production.Symbol(p.Idx)
	return p.Symbol != nil
}

// ProductionsIterator is used to iterate over a set of productions.
type ProductionsIterator struct {
	Production
	Productions Productions
	ln          int
	Idx         int
}

// Next moves Production to the next production in the productions. It returns
// false if there are no more productions.
func (p *ProductionsIterator) Next() bool {
	if p.ln == 0 {
		if p.Productions == nil {
			return false
		}
		p.ln = p.Productions.Productions()
		if p.ln == 0 {
			return false
		}
	} else {
		p.Idx++
	}
	if p.Idx >= p.ln {
		return false
	}
	p.Production = p.Productions.Production(p.Idx)
	return p.Symbol != nil
}
