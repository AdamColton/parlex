package grammar

import (
	"errors"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
	"strings"
)

// ErrBadGrammar is thrown when a Grammar Definition cannot be parsed to a
// Grammar
var ErrBadGrammar = errors.New("Bad Grammar")

// trimAndSplit takes a string of the form "A"
func trimAndSplit(str string) (prod stringsymbol.Production) {
	//TODO: add some tests for this
	strs := strings.Split(strings.TrimSpace(str), " ")
	if l := len(strs); l > 0 {
		symbols := make([]stringsymbol.Symbol, 0, l)
		for _, str := range strs {
			str = strings.TrimSpace(str) // trim tabs
			if str != "" {               // skip empty caused by two spaces in a row
				symbols = append(symbols, stringsymbol.Symbol(str))
			}
		}
		prod = stringsymbol.Production(symbols)
	}
	return
}

// productionFromLine takes a line and returns the non-terminal, and production,
// if the line is malformed, it will return an error. A blank line will not
// return an error but will return nil for the production. Either side may be
// blank
//   E -> E op E
//     -> int
//     ->
//   P ->
func productionFromLine(line string) (nt stringsymbol.Symbol, prod stringsymbol.Production, err error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return
	}
	p := strings.Split(line, "->")
	if l := len(p); l == 1 {
		prod = trimAndSplit(p[0])
	} else if l == 2 {
		nt = stringsymbol.Symbol(strings.TrimSpace(p[0]))
		prod = trimAndSplit(p[1])
	} else if l > 2 {
		err = ErrBadGrammar
	}
	return
}

// Grammar implements parlex.Grammar. It represents a set of production rules
// for a context free grammar.
type Grammar struct {
	order       []stringsymbol.Symbol
	productions map[stringsymbol.Symbol]stringsymbol.Productions
	longest     int
	totalCount  int
}

// New Grammar. The productions string should have one rule per line. A rule
// has the form "NonTerminal -> A B C" where A,B and C are symbols for either
// terminals or non-terminals. If there are multiple productions for a non-
// terminal, each row after the first can omit the non-terminal, as in "-> D E".
func New(productions string) (*Grammar, error) {
	g := &Grammar{
		productions: make(map[stringsymbol.Symbol]stringsymbol.Productions),
		longest:     -1,
	}
	cur := stringsymbol.Symbol("START") // current non terminal
	for _, line := range strings.Split(productions, "\n") {
		nt, prod, err := productionFromLine(line)
		if err != nil {
			return nil, err
		}
		if prod == nil {
			continue
		}
		if nt != "" {
			cur = nt
		}
		prods, defined := g.productions[cur]
		if !defined {
			g.order = append(g.order, cur)
		}
		g.productions[cur] = append(prods, prod)
	}
	return g, nil
}

// Empty Grammar. The productions string should have one rule per line. A rule
// has the form "NonTerminal -> A B C" where A,B and C are symbols for either
// terminals or non-terminals. If there are multiple productions for a non-
// terminal, each row after the first can omit the non-terminal, as in "-> D E".
func Empty() *Grammar {
	return &Grammar{
		productions: make(map[stringsymbol.Symbol]stringsymbol.Productions),
		longest:     -1,
	}
}

// Productions returns the productions for the given symbol. If the symbol is
// not a non-terminal in the Grammar, nil is returned. It is part of the
// parlex.Grammer interface.
func (g *Grammar) Productions(symbol parlex.Symbol) parlex.Productions {
	var p parlex.Productions
	s, ok := symbol.(stringsymbol.Symbol)
	if ok {
		p, ok = g.productions[s]
	} else {
		p, ok = g.productions[stringsymbol.Symbol(symbol.String())]
	}
	if ok {
		return p
	}
	return nil
}

// NonTerminals returns the non-terminals for the grammar. The first symbol in
// the list is the start symbol. The returned slice is a copy of the underlying
// structure and can be safely modified. It is part of the parlex.Grammer
// interface.
func (g *Grammar) NonTerminals() []parlex.Symbol {
	r := make([]parlex.Symbol, len(g.order))
	for i := range r {
		r[i] = g.order[i]
	}
	return r
}

// Add a production to the grammar.
func (g *Grammar) Add(from parlex.Symbol, to parlex.Production) {
	f, t := stringsymbol.CastSymbol(from), stringsymbol.CastProduction(to)
	if to == nil {
		to = make(stringsymbol.Production, 0)
	}
	prod, defined := g.productions[f]
	if !defined {
		g.order = append(g.order, f)
	}
	g.productions[f] = append(prod, t)
}

// String converts the grammar to a string. It aligns all the ->'s. The output
// of Grammar.String() can be used to define a copy of the grammar.
func (g *Grammar) String() string {
	return parlex.FormatGrammar(g)
}
