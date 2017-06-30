package packrat

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
	"github.com/adamcolton/parlex/tree"
)

// Packrat is a Packrat parser
type Packrat struct {
	parlex.Grammar
}

type treeMarker struct {
	symbol stringsymbol.Symbol
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
func New(grmr parlex.Grammar) *Packrat {
	return &Packrat{
		Grammar: grmr,
	}
}

// Constructor fulfills parlex.ParserConstructor
func Constructor(grmr parlex.Grammar) (parlex.Parser, error) {
	return &Packrat{
		Grammar: grmr,
	}, nil
}

// Parse fulfills the parlex.Parser. The Packrat parser will try to parse the
// lexemes.
func (p *Packrat) Parse(lexemes []parlex.Lexeme) parlex.ParseNode {
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
		symbol: stringsymbol.Symbol(nts[0].String()),
	}
	op.addProds(start)

	var u *updater
	for op.stack != nil {
		u, op.stack = op.stack, op.stack.next
		u.update(op)
	}

	var accept treeKey
	accept.symbol = stringsymbol.Symbol(nts[0].String())
	accept.end = len(lexemes)
	accepted := op.memo[accept]
	return accepted.toPN(lexemes, op.memo)
}

func (op *prOp) addProds(root treeMarker) {
	if op.queued[root] {
		return
	}
	op.queued[root] = true
	prods := op.grmr.Productions(root.symbol)
	if prods == nil {
		return
	}
	ln := prods.Productions()
	for pri := 0; pri < ln; pri++ {
		prod := prods.Production(pri)
		if prod.Symbols() == 0 {
			var nilTreeDef treeDef
			nilTreeDef.treeMarker = root
			nilTreeDef.end = root.start
			nilTreeDef.priority = pri
			op.addToMemo(nilTreeDef)
			continue
		}

		prodStart := treeMarker{
			symbol: stringsymbol.Symbol(prod.Symbol(0).String()),
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

	if ln+1 == extended.prod.Symbols() {
		op.addToMemo(extended.treeDef)
		return
	}

	requires := treeMarker{
		symbol: stringsymbol.Symbol(extended.prod.Symbol(ln + 1).String()),
		start:  extended.end,
	}

	if extended.end < len(op.lxms) && requires.symbol.String() == op.lxms[extended.end].Kind().String() {
		var td treeDef
		td.treeMarker = requires
		td.end = requires.start + 1
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
	} else if td.comparePriority(&old, op) == 1 && !op.createsCircularDep(td, &td) {
		op.memo[td.treeKey] = td
	}
}

func (op *prOp) createsCircularDep(node treeDef, root *treeDef) bool {
	for _, ck := range node.children {
		if ck == root.treeKey || op.createsCircularDep(op.memo[ck], root) {
			return true
		}
	}
	return false
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

	for i, ck1 := range td.children {
		ck2 := td2.children[i]
		if ck1 == ck2 {
			continue
		}
		c1, c2 := op.memo[ck1], op.memo[ck2]
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
	} else if requires.start < len(op.lxms) && requires.symbol.String() == op.lxms[requires.start].Kind().String() {
		var td treeDef
		td.treeMarker = requires
		td.end = requires.start + 1
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
	if td.start < len(lxms) && lxms[td.start].Kind().String() == td.symbol.String() {
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
