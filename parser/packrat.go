package parser

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/tree"
)

// PR is a Packrat parser
type PR struct {
	parlex.Grammar
}

// Packrat returns a Packrat parser
func Packrat(grmr parlex.Grammar) *PR {
	return &PR{
		Grammar: grmr,
	}
}

// Parse fulfills the parlex.Parser. The Packrat parser will try to parse the
// lexemes.
func (p *PR) Parse(lexemes []parlex.Lexeme) parlex.ParseNode {
	nts := p.Grammar.NonTerminals()
	if len(nts) == 0 {
		return nil
	}
	nonterm := make(map[parlex.Symbol]bool, len(nts))
	for _, nt := range nts {
		nonterm[nt] = true
	}
	op := &prOp{
		grmr:     p.Grammar,
		lxms:     lexemes,
		memo:     make(map[treeKey]treeDef),
		markers:  make(map[treeMarker][]treeDef),     // maps marker to treeDef containing that marker
		partials: make(map[treeMarker][]treePartial), // maps a marker to a treePartial looking for that marker
		queued:   make(map[treeMarker]bool),
		nonterm:  nonterm,
	}

	start := treeMarker{
		symbol: nts[0],
	}
	op.addProds(start)

	for op.stack != nil {
		u := op.stack
		op.stack = u.next
		u.update(op)
	}

	var accept treeKey
	accept.symbol = nts[0]
	accept.end = len(lexemes)

	return op.memo[accept].toPN(lexemes, op.memo)
}

func (op *prOp) addProds(marker treeMarker) {
	if op.queued[marker] {
		return
	}
	op.queued[marker] = true
	for pri, prod := range op.grmr.Productions(marker.symbol) {
		tm := treeMarker{
			symbol: prod[0],
			start:  marker.start,
		}
		var tp treePartial
		tp.treeMarker = marker
		tp.prod = prod
		tp.priority = pri

		op.addPartial(tm, tp)
	}
}

type prOp struct {
	grmr     parlex.Grammar
	lxms     []parlex.Lexeme
	memo     map[treeKey]treeDef
	markers  map[treeMarker][]treeDef
	partials map[treeMarker][]treePartial
	queued   map[treeMarker]bool
	nonterm  map[parlex.Symbol]bool
	stack    *updater
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
	next    *updater
	partial treePartial
	key     treeKey
}

func (u *updater) update(op *prOp) {
	var tp treePartial
	tp.treeKey = u.partial.treeKey
	tp.prod = u.partial.prod
	tp.priority = u.partial.priority
	ln := len(u.partial.children)
	tp.children = make([]treeKey, ln+1)
	copy(tp.children, u.partial.children)
	tp.children[ln] = u.key
	tp.end = u.key.end

	if ln+1 == len(tp.prod) {
		op.addToMemo(tp.treeDef)
		return
	}

	tm := treeMarker{
		symbol: tp.prod[ln+1],
		start:  tp.end,
	}
	if tp.end < len(op.lxms) && tm.symbol == op.lxms[tp.end].Kind() {
		var td treeDef
		td.treeMarker = tm
		td.end = tm.start + 1
		op.addToMemo(td)
		(&updater{
			partial: tp,
			key:     td.treeKey,
		}).update(op)
		return
	}

	op.addPartial(tm, tp)

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
		c1 := op.memo[ck]
		c2 := op.memo[td2.children[i]]
		p := c1.comparePriority(&c2, op)
		if p != 0 {
			return p
		}
	}
	return 0
}

func (op *prOp) addPartial(tm treeMarker, tp treePartial) {
	if tm.start >= len(op.lxms) {
		return
	}
	if op.nonterm[tm.symbol] {
		op.partials[tm] = append(op.partials[tm], tp)
	} else if tm.symbol == op.lxms[tm.start].Kind() {
		var td treeDef
		td.treeMarker = tm
		td.end = tm.start + 1
		op.addToMemo(td)
	}

	for _, td := range op.markers[tm] {
		op.push(tp, td.treeKey)
	}

	op.addProds(tm)
}

func (op *prOp) push(tp treePartial, tk treeKey) {
	op.stack = &updater{
		next:    op.stack,
		partial: tp,
		key:     tk,
	}
}

func (td treeDef) toPN(lxms []parlex.Lexeme, memo map[treeKey]treeDef) *tree.PN {
	if td.start == td.end {
		return nil
	}
	var lx parlex.Lexeme
	if lxms[td.start].Kind() == td.symbol {
		lx = lxms[td.start]
	} else {
		lx = &parlex.L{K: td.symbol}
	}
	pn := &tree.PN{
		Lexeme: lx,
		C:      make([]*tree.PN, len(td.children)),
	}
	for i, c := range td.children {
		cpn := memo[c].toPN(lxms, memo)
		cpn.P = pn
		pn.C[i] = cpn
	}
	return pn
}
