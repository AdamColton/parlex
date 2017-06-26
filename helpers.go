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
		if l := nt.Len(); l > longest {
			longest = l
		}
		totalCount += len(prods)
	}

	format := fmt.Sprintf("%%-%ds -> %%s", longest)
	segs := make([]string, 0, totalCount)
	for _, nt := range nonTerminals {
		prods := g.Productions(nt)
		segs = append(segs, fmt.Sprintf(format, string(nt), prods[0]))
		for _, prod := range prods[1:] {
			segs = append(segs, fmt.Sprintf(format, "", prod))
		}
	}
	return strings.Join(segs, "\n")
}

// LexemeString returns a Lexeme in the form "Kind : Value"
func LexemeString(ls ...Lexeme) string {
	if len(ls) == 0 {
		return ""
	}
	var strs = make([]string, len(ls))
	for i, l := range ls {
		k, v := string(l.Kind()), l.Value()
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
