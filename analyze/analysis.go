package analyze

import (
	"github.com/adamcolton/parlex"
)

// Analytics provides an analysis of a grammar.
type Analytics struct {
	parlex.Grammar
	first2nonterms map[parlex.Symbol]map[parlex.Symbol]bool
	nonterm2firsts map[parlex.Symbol][]parlex.Symbol
	nilInFirst     map[parlex.Symbol]bool
}

// Analyze a grammar. The grammar is embeded so the return value can be used as
// parlex.Grammar.
func Analyze(grammar parlex.Grammar) *Analytics {
	if a, ok := grammar.(*Analytics); ok {
		return a
	}
	a := &Analytics{
		Grammar:        grammar,
		first2nonterms: make(map[parlex.Symbol]map[parlex.Symbol]bool),
		nonterm2firsts: make(map[parlex.Symbol][]parlex.Symbol),
		nilInFirst:     make(map[parlex.Symbol]bool),
	}
	done := make(map[parlex.Symbol]bool)

	for _, symbol := range a.NonTerminals() {
		a.firsts(symbol, done)
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

func (a *Analytics) firsts(s parlex.Symbol, done map[parlex.Symbol]bool) ([]parlex.Symbol, bool) {
	if done[s] {
		return a.nonterm2firsts[s], a.nilInFirst[s]
	}
	done[s] = true
	var fs []parlex.Symbol
	nilInFirst := false
	for _, prod := range a.Productions(s) {
		if len(prod) == 0 {
			nilInFirst = true
			continue
		}
		var firsts []parlex.Symbol
		for i, doNext := 0, true; i < len(prod) && doNext; i++ {
			symbol := prod[i]
			if !a.NonTerminal(symbol) {
				doNext = false
				fs = append(fs, symbol)
			} else {
				firsts, doNext = a.firsts(symbol, done)
				nilInFirst = nilInFirst || doNext
				fs = append(fs, firsts...)
			}
		}
	}
	a.nonterm2firsts[s] = fs
	a.nilInFirst[s] = nilInFirst
	return fs, nilInFirst
}

func (a *Analytics) NonTerminal(symbol parlex.Symbol) bool {
	return !(a.Productions(symbol) == nil)
}

// HasFirst returns true if first could be the first symbol in a tree with root
// symbol
func (a *Analytics) HasFirst(symbol parlex.Symbol, first parlex.Symbol) bool {
	nonterms, ok := a.first2nonterms[first]
	if !ok {
		return symbol == first
	}
	return nonterms[symbol]
}

// HasNilInFirst will return true if it possible to reach an empty production
// through only left most children.
func (a *Analytics) HasNilInFirst(symbol parlex.Symbol) bool {
	return a.nilInFirst[symbol]
}
