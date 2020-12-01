// Package packrat implements a packrat parser based on
// http://web.cs.ucla.edu/~todd/research/pepm08.pdf . It can handle left
// recursion.
package packrat

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/symbol/setsymbol"
	"github.com/adamcolton/parlex/tree"
)

// Packrat is a Packrat parser
type Packrat struct {
	parlex.Grammar
}

type treeMarker struct {
	idx   int
	start int
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
	lxms     []*lexeme.Lexeme
	memo     map[treeKey]treeDef
	markers  map[treeMarker][]treeDef
	partials map[treeMarker][]treePartial
	queued   map[treeMarker]bool
	nonterms []bool
	stack    *updater
	set      *setsymbol.Set
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
	set := setsymbol.New()
	set.LoadGrammar(p.Grammar)
	op := &prOp{
		grmr:     p.Grammar,
		lxms:     set.LoadLexemes(lexemes),
		memo:     make(map[treeKey]treeDef),
		markers:  make(map[treeMarker][]treeDef),     // maps marker to treeDef containing that marker
		partials: make(map[treeMarker][]treePartial), // maps a marker to a treePartial looking for that marker
		queued:   make(map[treeMarker]bool),
		set:      set,
		nonterms: make([]bool, set.Size()),
	}
	for _, nonterm := range p.Grammar.NonTerminals() {
		op.nonterms[op.set.Symbol(nonterm).Idx()] = true
	}

	start := treeMarker{
		idx: op.set.Symbol(nts[0]).Idx(),
	}
	op.addProds(start)

	var u *updater
	for op.stack != nil {
		u, op.stack = op.stack, op.stack.next
		u.update(op)
	}

	var accept treeKey
	accept.idx = start.idx
	accept.end = len(lexemes)
	accepted := op.memo[accept]
	if accepted.end != accept.end {
		return nil
	}
	return accepted.toPN(op.lxms, op.memo, op.set)
}

func (op *prOp) addProds(root treeMarker) {
	if op.queued[root] {
		return
	}
	op.queued[root] = true
	rootSymbol := op.set.ByIdx(root.idx)
	prods := op.grmr.Productions(rootSymbol)
	if prods == nil {
		return
	}
	for i := prods.Iter(); i.Next(); {
		if i.Symbols() == 0 {
			var nilTreeDef treeDef
			nilTreeDef.treeMarker = root
			nilTreeDef.end = root.start
			nilTreeDef.priority = i.Idx
			op.addToMemo(nilTreeDef)
			continue
		}

		prodStart := treeMarker{
			idx:   op.set.Symbol(i.Symbol(0)).Idx(),
			start: root.start,
		}
		var prodPartial treePartial
		prodPartial.treeMarker = root
		prodPartial.prod = i.Production
		prodPartial.priority = i.Idx

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
		idx:   op.set.Symbol(extended.prod.Symbol(ln + 1)).Idx(),
		start: extended.end,
	}

	if td := op.checkNonTerminal(extended.treeMarker); td != nil {
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

func (op *prOp) addPartial(tp treePartial, requires treeMarker) {
	if op.nonterms[requires.idx] {
		op.partials[requires] = append(op.partials[requires], tp)
	} else {
		op.checkNonTerminal(requires)
	}

	for _, td := range op.markers[requires] {
		op.push(tp, td.treeKey)
	}

	op.addProds(requires)
}

func (op *prOp) checkNonTerminal(at treeMarker) *treeDef {
	matchesNonterminal := at.start < len(op.lxms) && at.idx == op.lxms[at.start].K.(*setsymbol.Symbol).Idx()
	if !matchesNonterminal {
		return nil
	}
	var td treeDef
	td.treeMarker = at
	td.end = at.start + 1
	op.addToMemo(td)
	return &td
}

func (op *prOp) push(tp treePartial, tk treeKey) {
	op.stack = &updater{
		next:      op.stack,
		base:      tp,
		extension: tk,
	}
}

func (td *treeDef) toPN(lxms []*lexeme.Lexeme, memo map[treeKey]treeDef, set *setsymbol.Set) *tree.PN {
	var lx *lexeme.Lexeme
	var setPos bool
	if td.start < len(lxms) && lxms[td.start].K.(*setsymbol.Symbol).Idx() == td.idx {
		lx = lxms[td.start]
	} else {
		lx = lexeme.New(set.ByIdx(td.idx))
		setPos = true
	}
	pn := &tree.PN{
		Lexeme: lx,
		C:      make([]*tree.PN, len(td.children)),
	}
	for i, c := range td.children {
		ct := memo[c]
		cpn := ct.toPN(lxms, memo, set)
		cpn.P = pn
		pn.C[i] = cpn
	}
	if setPos && len(pn.C) > 0 {
		lx.L, lx.C = pn.C[0].Pos()
	}
	return pn
}
