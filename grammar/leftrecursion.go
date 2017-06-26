package grammar

import (
	"github.com/adamcolton/parlex"
)

// IsLeftRecursive will check if a grammar is left recursive.
func IsLeftRecursive(grammar parlex.Grammar) bool {
	checked := make(map[parlex.Symbol]bool)
	for _, nt := range grammar.NonTerminals() {
		if checkLeftRecursion(grammar, nt, checked, make(map[parlex.Symbol]bool)) {
			return true
		}
	}
	return false
}

func checkLeftRecursion(g parlex.Grammar, s parlex.Symbol, checked, stack map[parlex.Symbol]bool) bool {
	if checked[s] {
		return false
	}

	prods := g.Productions(s)
	if prods == nil {
		// s is a terminal
		checked[s] = true // add to checked to avoid the lookup in the future
		return false
	}

	if stack[s] {
		return true
	}
	stack[s] = true

	for _, prod := range prods {
		if len(prod) > 0 && checkLeftRecursion(g, prod[0], checked, stack) {
			return true
		}
	}

	checked[s] = true
	return false
}
