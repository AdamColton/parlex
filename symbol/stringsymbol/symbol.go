// Package stringsymbol provides a simple implementations of parlex.Symbol,
// parlex.Production and parlex.Productions.
package stringsymbol

import (
	"github.com/adamcolton/parlex"
	"strings"
)

// Symbol fulfills parlex.Symbol. It simply wraps a string.
type Symbol string

// String returns the Symbol as a string, fulfilling parlex.Symbol.
func (s Symbol) String() string { return string(s) }

// CastSymbol takes a parlex.Symbol and returns an instance of
// stringsymbol.Symbol.
func CastSymbol(symbol parlex.Symbol) Symbol {
	if s, ok := symbol.(Symbol); ok {
		return s
	}
	return Symbol(symbol.String())
}

// Production fulfills parlex.Production
type Production []Symbol

// Iter returns ProductionIterator for iterating over the symbols in the
// production.
func (p Production) Iter() *parlex.ProductionIterator {
	return &parlex.ProductionIterator{
		Production: p,
	}
}

// CastProduction will cast any parlex.Production to simplestring.Production,
// including casting the underlying symbols.
func CastProduction(production parlex.Production) Production {
	if p, ok := production.(Production); ok {
		return p
	}
	if production == nil {
		return nil
	}
	ln := production.Symbols()
	p := make(Production, ln, ln)
	for i := 0; i < ln; i++ {
		p[i] = CastSymbol(production.Symbol(i))
	}
	return p
}

// String joins the symbols of the prodction with spaces
func (p Production) String() string {
	strs := make([]string, len(p))
	for i, symbol := range p {
		strs[i] = string(symbol)
	}
	return strings.Join(strs, " ")
}

// Symbols returns how many symbols are in the production
func (p Production) Symbols() int {
	if p == nil {
		return 0
	}
	return len(p)
}

// Symbol returns a symbol at a given index.
func (p Production) Symbol(i int) parlex.Symbol {
	if i < len(p) {
		return p[i]
	}
	return nil
}

// Productions fulfills parlex.Productions.
type Productions []Production

// Iter returns a ProductionsIterator.
func (p Productions) Iter() *parlex.ProductionsIterator {
	return &parlex.ProductionsIterator{
		Productions: p,
	}
}

// CastProductions will cast any parlex.Productions to stringsymbol.Productions
// including casting the underlying productions and the symbols underlying them.
func CastProductions(productions parlex.Productions) Productions {
	if p, ok := productions.(Productions); ok {
		return p
	}
	if productions == nil {
		return nil
	}
	ln := productions.Productions()
	p := make(Productions, ln)
	for i := 0; i < ln; i++ {
		p[i] = CastProduction(productions.Production(i))
	}
	return p
}

// Productions returns the number of productions
func (p Productions) Productions() int {
	if p == nil {
		return 0
	}
	return len(p)
}

// Production returns the production at a given index.
func (p Productions) Production(i int) parlex.Production {
	if i < len(p) {
		return p[i]
	}
	return nil
}
