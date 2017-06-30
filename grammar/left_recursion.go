package grammar

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/symbol/setsymbol"
)

type lrOp struct {
	in        parlex.Grammar
	out       *Grammar
	done      []bool
	cur       *setsymbol.Symbol
	hasDirect bool
	set       *setsymbol.Set
}

// RemoveLeftRecursion will convert a grammar with left recursion into one
// without.
func RemoveLeftRecursion(grammar parlex.Grammar) *Grammar {
	nts := grammar.NonTerminals()
	op := &lrOp{
		in:   grammar,
		out:  Empty(),
		done: make([]bool, len(nts)),
		set:  setsymbol.New(),
	}
	for _, s := range nts {
		op.cur = op.set.Symbol(s)
		op.hasDirect = false
		prods := grammar.Productions(op.cur)
		ln := prods.Productions()
		for i := 0; i < ln; i++ {
			op.safeAdd(op.set.CastProduction(prods.Production(i)))
		}
		if op.hasDirect {
			op.removeDirectLeftRecursion()
		}
		idx := op.cur.Idx()
		if idx >= len(op.done) {
			op.done = append(op.done, make([]bool, 1+idx+len(op.done))...)
		}
		op.done[op.cur.Idx()] = true
	}
	return op.out
}

func (op *lrOp) safeAdd(prod *setsymbol.Production) {
	var first *setsymbol.Symbol
	if prod.Symbols() > 0 {
		first = prod.Symbol(0).(*setsymbol.Symbol)
	}

	// We can directly add a production if it begins with a symobl we haven't
	// processed yet
	if first == nil || first.Idx() >= len(op.done) || !op.done[first.Idx()] {
		op.hasDirect = op.hasDirect || (first != nil && first.Idx() == op.cur.Idx())
		op.out.Add(op.cur, prod)
		return
	}

	tail := op.getTail(prod)
	prods := op.in.Productions(first)
	ln := prods.Productions()
	var lead *setsymbol.Production
	for i := 0; i < ln; i++ {
		lead = op.set.CastProduction(prods.Production(i))
		newProd := op.set.Production()
		ln := lead.Symbols()
		for j := 0; j < ln; j++ {
			newProd.AddSymbols(lead.Symbol(j))
		}
		ln = tail.Symbols()
		for j := 0; j < ln; j++ {
			newProd.AddSymbols(tail.Symbol(j))
		}
		op.safeAdd(newProd)
	}
}

func (l *lrOp) getTail(prod *setsymbol.Production) *setsymbol.Production {
	tail := l.set.Production()
	ln := prod.Symbols()
	for i := 1; i < ln; i++ {
		tail.AddSymbols(prod.Symbol(i))
	}
	return tail
}

func (op *lrOp) removeDirectLeftRecursion() {
	newSym := op.set.Str(op.cur.String() + "'")
	nsIdx := newSym.Idx()
	if nsIdx >= len(op.out.productions) || op.out.productions[nsIdx] == nil {
		op.out.order = append(op.out.order, nsIdx)
	}

	prods := op.out.productions[op.cur.Idx()]
	op.out.productions[op.cur.Idx()] = nil
	ln := prods.Productions()
	for i := 0; i < ln; i++ {
		prod := prods.Production(i).(*setsymbol.Production)
		if prod.Symbols() == 0 || prod.Symbol(0).(*setsymbol.Symbol).Idx() != op.cur.Idx() {
			prod.AddSymbols(newSym)
			op.directAdd(op.cur, prod)
		} else {
			for prod.Symbols() > 0 && prod.Symbol(0).(*setsymbol.Symbol).Idx() == op.cur.Idx() {
				prod = op.getTail(prod)
			}
			if prod.Symbols() > 0 {
				prod.AddSymbols(newSym)
				op.directAdd(newSym, prod)
			}
		}
	}
	op.directAdd(newSym, op.set.Production())
}

func (op *lrOp) directAdd(from *setsymbol.Symbol, to *setsymbol.Production) {
	f := from.Idx()
	var prods *setsymbol.Productions
	if f < len(op.out.productions) {
		prods = op.out.productions[f]
	} else {
		op.out.productions = append(op.out.productions, make([]*setsymbol.Productions, 1+f-len(op.out.productions))...)
	}

	if prods == nil {
		op.out.productions[f] = op.out.set.Productions(to)
	} else {
		prods.AddProductions(to)
	}
}
