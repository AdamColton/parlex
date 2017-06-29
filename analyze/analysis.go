package analyze

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
)

// Analytics provides an analysis of a grammar.
type Analytics struct {
	parlex.Grammar
	first2nonterms map[stringsymbol.Symbol]map[stringsymbol.Symbol]bool
	nonterm2firsts map[stringsymbol.Symbol][]stringsymbol.Symbol
	nilInFirst     map[stringsymbol.Symbol]bool
	Terminals      map[stringsymbol.Symbol]bool
}

// Analyze a grammar. The grammar is embeded so the return value can be used as
// parlex.Grammar.
func Analyze(grammar parlex.Grammar) *Analytics {
	if a, ok := grammar.(*Analytics); ok {
		return a
	}
	a := &Analytics{
		Grammar:        grammar,
		first2nonterms: make(map[stringsymbol.Symbol]map[stringsymbol.Symbol]bool),
		nonterm2firsts: make(map[stringsymbol.Symbol][]stringsymbol.Symbol),
		nilInFirst:     make(map[stringsymbol.Symbol]bool),
		Terminals:      make(map[stringsymbol.Symbol]bool),
	}
	done := make(map[stringsymbol.Symbol]bool)

	for _, symbol := range a.NonTerminals() {
		a.firsts(stringsymbol.CastSymbol(symbol), done)
	}

	for symbol, firsts := range a.nonterm2firsts {
		for _, first := range firsts {
			nonterms, ok := a.first2nonterms[first]
			if !ok {
				nonterms = make(map[stringsymbol.Symbol]bool)
				a.first2nonterms[first] = nonterms
			}
			nonterms[symbol] = true
		}
	}

	return a
}

func (a *Analytics) firsts(s stringsymbol.Symbol, done map[stringsymbol.Symbol]bool) ([]stringsymbol.Symbol, bool) {
	if done[s] {
		return a.nonterm2firsts[s], a.nilInFirst[s]
	}
	done[s] = true
	var fs []stringsymbol.Symbol
	nilInFirst := false
	prods := a.Productions(s)
	ln := prods.Productions()
	var prod parlex.Production
	for i := 0; i < ln; i++ {
		prod = prods.Production(i)
		ln := prod.Symbols()
		if ln == 0 {
			nilInFirst = true
			continue
		}
		var firsts []stringsymbol.Symbol
		for j, doNext := 0, true; j < ln && doNext; j++ {
			symbol := stringsymbol.CastSymbol(prod.Symbol(j))
			if !a.NonTerminal(symbol) {
				doNext = false
				fs = append(fs, symbol)
				a.Terminals[symbol] = true
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

// Contains returns true if the grammar contains the symbol as either a terminal
// or non-terminal.
func (a *Analytics) Contains(symbol parlex.Symbol) bool {
	s := stringsymbol.CastSymbol(symbol)
	return a.NonTerminal(s) || a.Terminals[s]
}

// NonTerminal returns true if the symbol is a non-terminal in the grammar.
func (a *Analytics) NonTerminal(symbol parlex.Symbol) bool {
	return a.Productions(symbol) != nil
}

// HasFirst returns true if first could be the first symbol in a tree with root
// symbol
func (a *Analytics) HasFirst(symbol parlex.Symbol, first parlex.Symbol) bool {
	f := stringsymbol.CastSymbol(first)
	s := stringsymbol.CastSymbol(symbol)
	nonterms, ok := a.first2nonterms[f]
	if !ok {
		return s == f
	}
	return nonterms[s]
}

// HasNilInFirst will return true if it possible to reach an empty production
// through only left most children.
func (a *Analytics) HasNilInFirst(symbol parlex.Symbol) bool {
	return a.nilInFirst[stringsymbol.CastSymbol(symbol)]
}
