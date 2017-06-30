package setsymbol

import (
	"github.com/adamcolton/parlex"
	"strings"
)

type Set struct {
	str2sym  map[string]int
	symb2str []string
}

func New() *Set {
	return &Set{
		str2sym: make(map[string]int),
	}
}

type Symbol struct {
	val int
	set *Set
}

func (s *Symbol) String() string { return s.set.symb2str[s.val] }

func (s *Symbol) Idx() int { return s.val }

func (s *Set) Size() int { return len(s.symb2str) }

func (s *Set) Symbol(symbol parlex.Symbol) *Symbol {
	if cast, ok := symbol.(*Symbol); ok && cast.set == s {
		return cast
	}
	return s.Str(symbol.String())
}

func (s *Set) Str(str string) *Symbol {
	if val, ok := s.str2sym[str]; ok {
		return &Symbol{
			val: val,
			set: s,
		}
	}
	sym := &Symbol{
		val: len(s.symb2str),
		set: s,
	}
	s.str2sym[str] = sym.val
	s.symb2str = append(s.symb2str, str)
	return sym
}

func (s *Set) Has(str string) bool {
	_, has := s.str2sym[str]
	return has
}

type Production struct {
	symbs []int
	set   *Set
}

func (s *Set) Production(symbols ...parlex.Symbol) *Production {
	p := &Production{
		symbs: make([]int, len(symbols)),
		set:   s,
	}
	for i, symbol := range symbols {
		p.symbs[i] = s.Symbol(symbol).val
	}
	return p
}

func (p *Production) AddSymbols(symbols ...parlex.Symbol) {
	for _, symbol := range symbols {
		p.symbs = append(p.symbs, p.set.Symbol(symbol).val)
	}
}

func (s *Set) CastProduction(production parlex.Production) *Production {
	if production == nil {
		return nil
	}
	if p, ok := production.(*Production); ok && p.set == s {
		return p
	}
	ln := production.Symbols()
	p := &Production{
		symbs: make([]int, ln, ln),
		set:   s,
	}
	for i := 0; i < ln; i++ {
		p.symbs[i] = s.Symbol(production.Symbol(i)).val
	}
	return p
}

// String joins the symbols of the prodction with spaces
func (p *Production) String() string {
	strs := make([]string, len(p.symbs))
	for i, idx := range p.symbs {
		strs[i] = p.set.symb2str[idx]
	}
	return strings.Join(strs, " ")
}

func (p *Production) Symbols() int {
	if p == nil {
		return 0
	}
	return len(p.symbs)
}

func (p *Production) Symbol(i int) parlex.Symbol {
	if i < len(p.symbs) {
		return &Symbol{
			val: p.symbs[i],
			set: p.set,
		}
	}
	return nil
}

type Productions struct {
	prods [][]int
	set   *Set
}

func (s *Set) Productions(productions ...parlex.Production) *Productions {
	p := &Productions{
		prods: make([][]int, len(productions)),
		set:   s,
	}
	for i, prod := range productions {
		p.prods[i] = s.CastProduction(prod).symbs
	}
	return p
}

func (s *Set) CastProductions(productions parlex.Productions) *Productions {
	if productions == nil {
		return nil
	}
	if p, ok := productions.(*Productions); ok && p.set == s {
		return p
	}
	ln := productions.Productions()
	p := &Productions{
		prods: make([][]int, ln, ln),
		set:   s,
	}
	for i := 0; i < ln; i++ {
		p.prods[i] = s.CastProduction(productions.Production(i)).symbs
	}
	return p
}

func (p *Productions) Productions() int {
	if p == nil {
		return 0
	}
	return len(p.prods)
}

func (p *Productions) Production(i int) parlex.Production {
	if i < len(p.prods) {
		return &Production{
			symbs: p.prods[i],
			set:   p.set,
		}
	}
	return nil
}
