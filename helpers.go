package parlex

import (
	"fmt"
	"strings"
)

// GrammarString formats a grammar into a string
func GrammarString(g Grammar) string {
	longest := -1
	totalCount := 0
	nonTerminals := g.NonTerminals()
	for _, nt := range nonTerminals {
		prods := g.Productions(nt)
		if kl := nt.Len(); kl > longest {
			longest = kl
		}
		totalCount += len(prods)
	}

	format := fmt.Sprintf("%%-%ds -> %%s", longest)
	segs := make([]string, totalCount)
	s := 0
	for _, nt := range nonTerminals {
		prods := g.Productions(nt)
		ks := string(nt)
		for i, prod := range prods {
			segs[s] = fmt.Sprintf(format, ks, prod)
			s++
			if i == 0 {
				ks = "" //don't print non terminal after first production
			}
		}
	}
	return strings.Join(segs, "\n")
}

// LexemeString returns a Lexeme in the form "Kind : Value"
func LexemeString(l Lexeme) string {
	k, v := string(l.Kind()), l.Value()
	if v == "" {
		return k
	}
	return k + " : " + v
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
