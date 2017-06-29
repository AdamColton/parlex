package stringsymbol

import (
	"github.com/adamcolton/parlex"
	"strings"
)

type Symbol string

func (s Symbol) String() string { return string(s) }

func CastSymbol(symbol parlex.Symbol) Symbol {
	if s, ok := symbol.(Symbol); ok {
		return s
	}
	return Symbol(symbol.String())
}

type Production []Symbol

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

func (p Production) Symbols() int {
	if p == nil {
		return 0
	}
	return len(p)
}

func (p Production) Symbol(i int) parlex.Symbol {
	if i < len(p) {
		return p[i]
	}
	return nil
}

type Productions []Production

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

func (p Productions) Productions() int {
	if p == nil {
		return 0
	}
	return len(p)
}

func (p Productions) Production(i int) parlex.Production {
	if i < len(p) {
		return p[i]
	}
	return nil
}
