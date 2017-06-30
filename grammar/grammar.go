package grammar

import (
	"errors"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/symbol/setsymbol"
	"strings"
)

// ErrBadGrammar is thrown when a Grammar Definition cannot be parsed to a
// Grammar
var ErrBadGrammar = errors.New("Bad Grammar")

// trimAndSplit takes a string of the form " A  B  C " and converts it to a
// production
func (g *Grammar) trimAndSplit(str string) *setsymbol.Production {
	//TODO: add some tests for this
	strs := strings.Split(strings.TrimSpace(str), " ")
	symbols := make([]parlex.Symbol, 0, len(strs))
	for _, str := range strs {
		str = strings.TrimSpace(str) // trim tabs
		if str != "" {               // skip empty caused by two spaces in a row
			symbols = append(symbols, g.set.Str(str))
		}
	}

	return g.set.Production(symbols...)
}

// productionFromLine takes a line and returns the non-terminal, and production,
// if the line is malformed, it will return an error. A blank line will not
// return an error but will return nil for the production. Either side may be
// blank
//   E -> E op E
//     -> int
//     ->
//   P ->
func (g *Grammar) productionFromLine(line string) (nt *setsymbol.Symbol, prod *setsymbol.Production, err error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return
	}
	p := strings.Split(line, "->")
	if l := len(p); l == 1 {
		prod = g.trimAndSplit(p[0])
	} else if l == 2 {
		ntstr := strings.TrimSpace(p[0])
		if ntstr != "" {
			nt = g.set.Str(ntstr)
		}
		prod = g.trimAndSplit(p[1])
	} else if l > 2 {
		err = ErrBadGrammar
	}
	return
}

// Grammar implements parlex.Grammar. It represents a set of production rules
// for a context free grammar.
type Grammar struct {
	order       []int
	productions []*setsymbol.Productions
	longest     int
	totalCount  int
	set         *setsymbol.Set
}

// New Grammar. The productions string should have one rule per line. A rule
// has the form "NonTerminal -> A B C" where A,B and C are symbols for either
// terminals or non-terminals. If there are multiple productions for a non-
// terminal, each row after the first can omit the non-terminal, as in "-> D E".
func New(productions string) (*Grammar, error) {
	g := &Grammar{
		longest: -1,
		set:     setsymbol.New(),
	}
	cur := -1
	for _, line := range strings.Split(productions, "\n") {
		nt, prod, err := g.productionFromLine(line)
		if err != nil {
			return nil, err
		}
		if prod == nil {
			continue
		}
		if nt != nil {
			cur = nt.Idx()
		} else if cur == -1 {
			cur = g.set.Str("START").Idx()
		}
		var prods *setsymbol.Productions
		if cur < len(g.productions) {
			prods = g.productions[cur]
		} else {
			g.productions = append(g.productions, make([]*setsymbol.Productions, 1+cur-len(g.productions))...)
		}
		if prods == nil {
			g.productions[cur] = g.set.Productions(prod)
			g.order = append(g.order, cur)
		} else {
			prods.AddProductions(prod)
		}
	}
	return g, nil
}

// Empty Grammar. The productions string should have one rule per line. A rule
// has the form "NonTerminal -> A B C" where A,B and C are symbols for either
// terminals or non-terminals. If there are multiple productions for a non-
// terminal, each row after the first can omit the non-terminal, as in "-> D E".
func Empty() *Grammar {
	return &Grammar{
		longest: -1,
		set:     setsymbol.New(),
	}
}

// Productions returns the productions for the given symbol. If the symbol is
// not a non-terminal in the Grammar, nil is returned. It is part of the
// parlex.Grammer interface.
func (g *Grammar) Productions(symbol parlex.Symbol) parlex.Productions {
	idx := g.set.Symbol(symbol).Idx()
	if idx >= len(g.productions) {
		return nil
	}
	p := g.productions[idx]
	if p == nil {
		return nil
	}
	return p
}

// NonTerminals returns the non-terminals for the grammar. The first symbol in
// the list is the start symbol. The returned slice is a copy of the underlying
// structure and can be safely modified. It is part of the parlex.Grammer
// interface.
func (g *Grammar) NonTerminals() []parlex.Symbol {
	r := make([]parlex.Symbol, len(g.order))
	for i := range r {
		r[i] = g.set.ByIdx(g.order[i])
	}
	return r
}

// Add a production to the grammar.
func (g *Grammar) Add(from parlex.Symbol, to parlex.Production) {
	f := g.set.Symbol(from).Idx()
	if to == nil {
		to = g.set.Production()
	}

	var prods *setsymbol.Productions
	if f < len(g.productions) {
		prods = g.productions[f]
	} else {
		g.productions = append(g.productions, make([]*setsymbol.Productions, 1+f-len(g.productions))...)
	}

	if prods == nil {
		g.productions[f] = g.set.Productions(to)
		g.order = append(g.order, f)
	} else {
		prods.AddProductions(to)
	}
}

// String converts the grammar to a string. It aligns all the ->'s. The output
// of Grammar.String() can be used to define a copy of the grammar.
func (g *Grammar) String() string {
	return parlex.FormatGrammar(g)
}
