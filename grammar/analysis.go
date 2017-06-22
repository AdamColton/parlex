package grammar

import (
	"github.com/adamcolton/parlex"
)

type Analytics struct {
	first2nonterms map[parlex.Symbol]map[parlex.Symbol]bool
	nonterm2firsts map[parlex.Symbol][]parlex.Symbol
}

func Analyse(g parlex.Grammar) *Analytics {
	a := &Analytics{
		first2nonterms: make(map[parlex.Symbol]map[parlex.Symbol]bool),
		nonterm2firsts: make(map[parlex.Symbol][]parlex.Symbol),
	}
	done := make(map[parlex.Symbol]bool)

	for _, symbol := range g.NonTerminals() {
		a.nonterm2firsts[symbol] = nil
	}

	for _, symbol := range g.NonTerminals() {
		firsts(symbol, g, a, done)
	}

	for symbol, firsts := range a.nonterm2firsts {
		for _, first := range firsts {
			nonterms, ok := a.first2nonterms[first]
			if !ok {
				nonterms = make(map[parlex.Symbol]bool)
				a.first2nonterms[first] = nonterms
			}
			nonterms[symbol] = true
		}
	}

	return a
}

func firsts(s parlex.Symbol, g parlex.Grammar, a *Analytics, done map[parlex.Symbol]bool) []parlex.Symbol {
	if done[s] {
		return a.nonterm2firsts[s]
	}
	done[s] = true
	var fs []parlex.Symbol
	for _, prod := range g.Productions(s) {
		symbol := prod[0]
		if a.Terminal(symbol) {
			fs = append(fs, symbol)
		} else {
			fs = append(fs, firsts(symbol, g, a, done)...)
		}
	}
	a.nonterm2firsts[s] = fs
	return fs
}

func (a *Analytics) Terminal(s parlex.Symbol) bool {
	_, nonterm := a.nonterm2firsts[s]
	return !nonterm
}

func (a *Analytics) HasFirst(nonterm parlex.Symbol, first parlex.Symbol) bool {
	nonterms, ok := a.first2nonterms[first]
	if !ok {
		return false
	}
	return nonterms[nonterm]
}
