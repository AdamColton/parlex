package packrat

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/tree"
)

// PR is a Packrat parser
type PR struct {
	parlex.Grammar
}

type treeMarker struct {
	symbol parlex.Symbol
	start  int
}

type treeKey struct {
	treeMarker
	end int
}

type treeDef struct {
	treeKey
	children []treeKey
	priority int
}

type treePartial struct {
	treeDef
	prod parlex.Production
}

// updaters form a linked-list of things to process. An updater takes the
// treePartial and creates a new one extended with the tree key
type updater struct {
	next      *updater
	base      treePartial
	extension treeKey
}

// pack-rat parse operation
type prOp struct {
	grmr     parlex.Grammar
	lxms     []parlex.Lexeme
	memo     map[treeKey]treeDef
	markers  map[treeMarker][]treeDef
	partials map[treeMarker][]treePartial
	queued   map[treeMarker]bool
	stack    *updater
}

// New returns a Packrat parser
func New(grmr parlex.Grammar) *PR {
	return &PR{
		Grammar: grmr,
	}
}

// Constructor fulfills parlex.ParserConstructor
func Constructor(grmr parlex.Grammar) (parlex.Parser, error) {
	return &PR{
		Grammar: grmr,
	}, nil
}

// Parse fulfills the parlex.Parser. The Packrat parser will try to parse the
// lexemes.
func (p *PR) Parse(lexemes []parlex.Lexeme) parlex.ParseNode {
	nts := p.Grammar.NonTerminals()
	if len(nts) == 0 {
		return nil
	}
	op := &prOp{
		grmr:     p.Grammar,
		lxms:     lexemes,
		memo:     make(map[treeKey]treeDef),
		markers:  make(map[treeMarker][]treeDef),     // maps marker to treeDef containing that marker
		partials: make(map[treeMarker][]treePartial), // maps a marker to a treePartial looking for that marker
		queued:   make(map[treeMarker]bool),
	}

	start := treeMarker{
		symbol: nts[0],
	}
	op.addProds(start)

	var u *updater
	for op.stack != nil {
		u, op.stack = op.stack, op.stack.next
		u.update(op)
	}

	var accept treeKey
	accept.symbol = nts[0]
	accept.end = len(lexemes)
	accepted := op.memo[accept]
	return accepted.toPN(lexemes, op.memo)
}

func (op *prOp) addProds(root treeMarker) {
	if op.queued[root] {
		return
	}
	op.queued[root] = true
	for pri, prod := range op.grmr.Productions(root.symbol) {
		prodStart := treeMarker{
			symbol: prod[0],
			start:  root.start,
		}
		var prodPartial treePartial
		prodPartial.treeMarker = root
		prodPartial.prod = prod
		prodPartial.priority = pri

		op.addPartial(prodPartial, prodStart)
	}
}

func (u *updater) update(op *prOp) {
	extended := u.base
	ln := len(u.base.children)
	extended.children = make([]treeKey, ln+1)
	copy(extended.children, u.base.children)
	extended.children[ln] = u.extension
	extended.end = u.extension.end

	if ln+1 == len(extended.prod) {
		op.addToMemo(extended.treeDef)
		return
	}

	requires := treeMarker{
		symbol: extended.prod[ln+1],
		start:  extended.end,
	}

	if (extended.end < len(op.lxms) && requires.symbol == op.lxms[extended.end].Kind()) || requires.symbol == "NIL" {
		var td treeDef
		td.treeMarker = requires
		if requires.symbol == "NIL" {
			td.end = requires.start
		} else {
			td.end = requires.start + 1
		}
		op.addToMemo(td)
		(&updater{
			base:      extended,
			extension: td.treeKey,
		}).update(op)
		return
	}

	op.addPartial(extended, requires)
}

func (op *prOp) addToMemo(td treeDef) {
	old, ok := op.memo[td.treeKey]
	if !ok {
		op.memo[td.treeKey] = td
		op.markers[td.treeMarker] = append(op.markers[td.treeMarker], td)
		for _, tp := range op.partials[td.treeMarker] {
			op.push(tp, td.treeKey)
		}
	} else if td.comparePriority(&old, op) == 1 {
		op.memo[td.treeKey] = td
	}
}

//  1: td > td2    which actually means td.priority < td2.priority
//  0: td == td2   because 0 is the highest priority
// -1: td < td2
func (td *treeDef) comparePriority(td2 *treeDef, op *prOp) int8 {
	if td.priority != td2.priority {
		if td.priority < td2.priority {
			return 1
		}
		return -1
	}
	// they should both have the same number of children - equal priority means
	// equal production
	for i, ck := range td.children {
		c1, c2 := op.memo[ck], op.memo[td2.children[i]]
		p := c1.comparePriority(&c2, op)
		if p != 0 {
			return p
		}
	}
	return 0
}

func (op *prOp) nonterm(symbol parlex.Symbol) bool {
	return op.grmr.Productions(symbol) != nil
}

func (op *prOp) addPartial(tp treePartial, requires treeMarker) {
	if op.nonterm(requires.symbol) {
		op.partials[requires] = append(op.partials[requires], tp)
	} else if (requires.start < len(op.lxms) && requires.symbol == op.lxms[requires.start].Kind()) || requires.symbol == "NIL" {
		var td treeDef
		td.treeMarker = requires
		if requires.symbol == "NIL" {
			td.end = requires.start
		} else {
			td.end = requires.start + 1
		}
		op.addToMemo(td)
	}

	for _, td := range op.markers[requires] {
		op.push(tp, td.treeKey)
	}

	op.addProds(requires)
}

func (op *prOp) push(tp treePartial, tk treeKey) {
	op.stack = &updater{
		next:      op.stack,
		base:      tp,
		extension: tk,
	}
}

func (td *treeDef) toPN(lxms []parlex.Lexeme, memo map[treeKey]treeDef) *tree.PN {
	var lx parlex.Lexeme
	if lxms[td.start].Kind() == td.symbol {
		lx = lxms[td.start]
	} else {
		lx = lexeme.New(td.symbol)
	}
	pn := &tree.PN{
		Lexeme: lx,
		C:      make([]*tree.PN, len(td.children)),
	}
	for i, c := range td.children {
		ct := memo[c]
		cpn := ct.toPN(lxms, memo)
		cpn.P = pn
		pn.C[i] = cpn
	}
	return pn
}
