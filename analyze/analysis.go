package analyze

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/symbol/setsymbol"
)

// Analytics provides an analysis of a grammar.
type Analytics struct {
	parlex.Grammar
	first2nonterms [][]bool
	nonterm2firsts [][]int
	nilInFirst     []bool
	Terminals      []bool
	set            *setsymbol.Set
	done           []bool
}

// Analyze a grammar. The grammar is embeded so the return value can be used as
// parlex.Grammar.
func Analyze(grammar parlex.Grammar) *Analytics {
	if a, ok := grammar.(*Analytics); ok {
		return a
	}

	// todo: Move this into Set as LoadGrammar
	set := setsymbol.New()
	set.LoadGrammar(grammar)
	ln := set.Size()

	a := &Analytics{
		Grammar:        grammar,
		set:            set,
		first2nonterms: make([][]bool, ln),
		nonterm2firsts: make([][]int, ln),
		nilInFirst:     make([]bool, ln),
		Terminals:      make([]bool, ln),
		done:           make([]bool, ln),
	}

	// find firsts by populating nonterm2firsts
	for _, symbol := range a.NonTerminals() {
		a.firsts(a.set.Symbol(symbol))
	}

	// populate first2nonterms from nonterm2firsts
	for symbol, firsts := range a.nonterm2firsts {
		for _, first := range firsts {
			var nonterms []bool
			nonterms = a.first2nonterms[first]
			if nonterms == nil {
				nonterms = make([]bool, ln)
				a.first2nonterms[first] = nonterms
			}
			nonterms[symbol] = true
		}
	}

	return a
}

func (a *Analytics) firsts(s *setsymbol.Symbol) ([]int, bool) {
	sIdx := s.Idx()
	if a.done[sIdx] {
		return a.nonterm2firsts[sIdx], a.nilInFirst[sIdx]
	}
	a.done[sIdx] = true

	var fs []int
	nilInFirst := false
	for i := a.Productions(s).Iter(); i.Next(); {
		if i.Symbols() == 0 {
			nilInFirst = true
			continue
		}
		var firsts []int
		for j, doNext := i.Iter(), true; j.Next() && doNext; {
			symbol := a.set.Symbol(j.Symbol)
			if a.Productions(symbol) == nil {
				doNext = false
				fs = append(fs, symbol.Idx())
				a.Terminals[symbol.Idx()] = true
			} else {
				firsts, doNext = a.firsts(symbol)
				nilInFirst = nilInFirst || doNext
				fs = append(fs, firsts...)
			}
		}
	}
	a.nonterm2firsts[sIdx] = fs
	a.nilInFirst[sIdx] = nilInFirst
	return fs, nilInFirst
}

// Contains returns true if the grammar contains the symbol as either a terminal
// or non-terminal.
func (a *Analytics) Contains(symbol parlex.Symbol) bool {
	s := a.set.HasSymbol(symbol)
	if s == nil {
		return false
	}
	idx := s.Idx()
	return a.nonterm2firsts[idx] != nil || a.Terminals[idx]
}

// NonTerminal returns true if the symbol is a non-terminal in the grammar.
func (a *Analytics) NonTerminal(symbol parlex.Symbol) bool {
	s := a.set.HasSymbol(symbol)
	if s == nil {
		return false
	}
	idx := s.Idx()
	if idx >= len(a.nonterm2firsts) {
		return false
	}
	return a.nonterm2firsts[idx] != nil
}

// HasFirst returns true if first could be the first symbol in a tree with root
// symbol
func (a *Analytics) HasFirst(symbol parlex.Symbol, first parlex.Symbol) bool {
	f := a.set.HasSymbol(first)
	if f == nil {
		return false
	}
	s := a.set.HasSymbol(symbol)
	if s == nil {
		return false
	}
	nonterms := a.first2nonterms[f.Idx()]
	if nonterms == nil {
		return s.Idx() == f.Idx()
	}
	return nonterms[s.Idx()]
}

// HasNilInFirst will return true if it possible to reach an empty production
// through only left most children.
func (a *Analytics) HasNilInFirst(symbol parlex.Symbol) bool {
	s := a.set.HasSymbol(symbol)
	if s == nil {
		return false
	}
	return a.nilInFirst[s.Idx()]
}
