package grammar

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
)

type lrOp struct {
	in        parlex.Grammar
	out       *Grammar
	done      map[stringsymbol.Symbol]bool
	cur       stringsymbol.Symbol
	hasDirect bool
}

// RemoveLeftRecursion will convert a grammar with left recursion into one
// without.
func RemoveLeftRecursion(grammar parlex.Grammar) parlex.Grammar {
	nts := grammar.NonTerminals()
	op := &lrOp{
		in:   grammar,
		out:  Empty(),
		done: make(map[stringsymbol.Symbol]bool, len(nts)),
	}
	for _, s := range nts {
		op.cur = stringsymbol.CastSymbol(s)
		op.hasDirect = false
		prods := grammar.Productions(op.cur)
		ln := prods.Productions()
		for i := 0; i < ln; i++ {
			op.safeAdd(stringsymbol.CastProduction(prods.Production(i)))
		}
		if op.hasDirect {
			op.removeDirectLeftRecursion()
		}
		op.done[op.cur] = true
	}
	return op.out
}

func (op *lrOp) safeAdd(prod stringsymbol.Production) {
	var first stringsymbol.Symbol
	if len(prod) > 0 {
		first = prod[0]
	}
	if !op.done[first] {
		op.hasDirect = op.hasDirect || first == op.cur
		op.out.Add(op.cur, prod)
		return
	}

	tail := prod[1:]
	prods := op.in.Productions(first)
	ln := prods.Productions()
	var lead stringsymbol.Production
	for i := 0; i < ln; i++ {
		lead = stringsymbol.CastProduction(prods.Production(i))
		newProd := make(stringsymbol.Production, len(lead)+len(tail))
		copy(newProd, lead)
		copy(newProd[len(lead):], tail)
		op.safeAdd(newProd)
	}
}

func (op *lrOp) removeDirectLeftRecursion() {
	newSym := op.cur + "'"

	prods := op.out.productions[op.cur]
	op.out.productions[op.cur] = nil
	for _, prod := range prods {
		if len(prod) == 0 || prod[0] != op.cur {
			op.out.productions[op.cur] = append(op.out.productions[op.cur], append(prod, newSym))
		} else {
			for len(prod) > 0 && prod[0] == op.cur {
				prod = prod[1:]
			}
			if len(prod) > 0 {
				prod = append(prod, newSym)
				op.out.Add(newSym, prod)
			}
		}
	}
	op.out.Add(newSym, make(stringsymbol.Production, 0, 0))
}
