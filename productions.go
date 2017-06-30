package parlex

// Symbol is base of a grammar.
type Symbol interface {
	String() string
}

// Production is a slice of symbols. It is not actually the full grammatic
// production because it does not contain the left side of the production.
type Production interface {
	Symbols() int
	Symbol(int) Symbol
	Iter() *ProductionIterator
}

// Productions are used to represent the set of productions available from a
// non-terminal
type Productions interface {
	Productions() int
	Production(int) Production
	Iter() *ProductionsIterator
}

func SymLen(s Symbol) int {
	return len([]rune(s.String()))
}

type ProductionIterator struct {
	Symbol
	Production Production
	ln         int
	Idx        int
}

func (p *ProductionIterator) Next() bool {
	if p.ln == 0 {
		if p.Production == nil {
			return false
		}
		p.ln = p.Production.Symbols()
		if p.ln == 0 {
			return false
		}
	} else {
		p.Idx++
	}
	if p.Idx >= p.ln {
		return false
	}
	p.Symbol = p.Production.Symbol(p.Idx)
	return p.Symbol != nil
}

type ProductionsIterator struct {
	Production
	Productions Productions
	ln          int
	Idx         int
}

func (p *ProductionsIterator) Next() bool {
	if p.ln == 0 {
		if p.Productions == nil {
			return false
		}
		p.ln = p.Productions.Productions()
		if p.ln == 0 {
			return false
		}
	} else {
		p.Idx++
	}
	if p.Idx >= p.ln {
		return false
	}
	p.Production = p.Productions.Production(p.Idx)
	return p.Symbol != nil
}
