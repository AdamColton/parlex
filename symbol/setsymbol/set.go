package setsymbol

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
	"strings"
)

// Set is a collection of symbols and can map a symbol string to an int or an
// int to it's string.
type Set struct {
	str2sym  map[string]int
	symb2str []string
}

// New creates a new, empty set.
func New() *Set {
	return &Set{
		str2sym: make(map[string]int),
	}
}

// Size returns the number of symbols in the set
func (s *Set) Size() int { return len(s.symb2str) }

// ByIdx gets a symbol by index from the set
func (s *Set) ByIdx(idx int) *Symbol {
	if idx > len(s.symb2str) {
		return nil
	}
	return &Symbol{
		val: idx,
		set: s,
	}
}

// Idx takes a symbol and returns it's index value. If the symbol is not already
// in the set, it will be added.
func (s *Set) Idx(symbol parlex.Symbol) int {
	if cast, ok := symbol.(*Symbol); ok && cast.set == s {
		return cast.val
	}
	idx, found := s.str2sym[symbol.String()]
	if !found {
		return -1
	}
	return idx
}

// Symbol takes any parlex.Symbol and returns a setsymbol.Symbol
func (s *Set) Symbol(symbol parlex.Symbol) *Symbol {
	if symbol == nil {
		return nil
	}
	if cast, ok := symbol.(*Symbol); ok && cast.set == s {
		return cast
	}
	return s.Str(symbol.String())
}

// Str takes any string and returns it as symbol. If it is not already a part of
// the set, it is added.
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

// Has returns true if the set already contains the string
func (s *Set) Has(str string) bool {
	_, has := s.str2sym[str]
	return has
}

// HasSymbol returns true if the set already contains the symbol.
func (s *Set) HasSymbol(symbol parlex.Symbol) *Symbol {
	if cast, ok := symbol.(*Symbol); ok && cast.set == s {
		return cast
	}
	idx, has := s.str2sym[symbol.String()]
	if !has {
		return nil
	}
	return &Symbol{
		val: idx,
		set: s,
	}
}

// Symbol represents a symbol as int belonging to a set
type Symbol struct {
	val int
	set *Set
}

// String gets the string for symbol
func (s *Symbol) String() string { return s.set.symb2str[s.val] }

// Idx gets the symbol's index
func (s *Symbol) Idx() int { return s.val }

// Production implements parlex.Production using set symbols.
type Production struct {
	symbs []int
	set   *Set
}

// Production returns a production using the set
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

// Iter returns a ProductionIterator
func (p *Production) Iter() *parlex.ProductionIterator {
	return &parlex.ProductionIterator{
		Production: p,
	}
}

// AddSymbols appends symbols to the production
func (p *Production) AddSymbols(symbols ...parlex.Symbol) {
	for _, symbol := range symbols {
		p.symbs = append(p.symbs, p.set.Symbol(symbol).val)
	}
}

// CastProduction will take any parlex.Production and cast it to a
// setsymbol.Production including the underlying symbols.
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

// Symbols returns the number of symbols in the production
func (p *Production) Symbols() int {
	if p == nil {
		return 0
	}
	return len(p.symbs)
}

// Symbol returns the symbol at the given index
func (p *Production) Symbol(i int) parlex.Symbol {
	if i < len(p.symbs) {
		return &Symbol{
			val: p.symbs[i],
			set: p.set,
		}
	}
	return nil
}

// Productions fulfills parlex.Productions
type Productions struct {
	prods [][]int
	set   *Set
}

// Iter returns a ProductionsIterator
func (p *Productions) Iter() *parlex.ProductionsIterator {
	return &parlex.ProductionsIterator{
		Productions: p,
	}
}

// String will represent a production as a string
func (p *Productions) String() string {
	if p == nil {
		return "{nil}"
	}
	strsout := make([]string, len(p.prods))
	for i, prod := range p.prods {
		strsin := make([]string, len(prod))
		for j, idx := range prod {
			strsin[j] = p.set.symb2str[idx]
		}
		strsout[i] = "[" + strings.Join(strsin, ", ") + "]"
	}

	return "{" + strings.Join(strsout, " ") + "}"
}

// Productions will create a Productions from the given productions
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

// CastProductions will cast any parlex.Productions to a setsymbol.Productions.
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

// Productions returns the number of productions
func (p *Productions) Productions() int {
	if p == nil {
		return 0
	}
	return len(p.prods)
}

// Production returns the production at the given index
func (p *Productions) Production(i int) parlex.Production {
	if i < len(p.prods) {
		return &Production{
			symbs: p.prods[i],
			set:   p.set,
		}
	}
	return nil
}

// AddProductions appends a production
func (p *Productions) AddProductions(productions ...parlex.Production) {
	for _, prod := range productions {
		p.prods = append(p.prods, p.set.CastProduction(prod).symbs)
	}
}

// LoadGrammar takes a grammar and loads all the symbols into the set. This is
// useful to get the size of the set.
func (s *Set) LoadGrammar(grammar parlex.Grammar) {
	for _, nt := range grammar.NonTerminals() {
		s.Symbol(nt)
		for i := grammar.Productions(nt).Iter(); i.Next(); {
			for j := i.Iter(); j.Next(); {
				s.Symbol(j.Symbol)
			}
		}
	}
}

// LoadLexemes takes a slice of Lexemes and loads the symbols. This is
// useful to get the size of the set.
func (s *Set) LoadLexemes(lexemes []parlex.Lexeme) []*lexeme.Lexeme {
	out := make([]*lexeme.Lexeme, len(lexemes))
	for i, lx := range lexemes {
		out[i] = lexeme.New(s.Symbol(lx.Kind())).Set(lx.Value()).At(lx.Pos())
	}
	return out
}
