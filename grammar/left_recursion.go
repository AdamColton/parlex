package grammar

import (
	"github.com/adamcolton/parlex"
)

type lrOp struct {
	in        parlex.Grammar
	out       *Grammar
	done      map[parlex.Symbol]bool
	cur       parlex.Symbol
	hasDirect bool
}

// RemoveLeftRecursion will convert a grammar with left recursion into one
// without.
func RemoveLeftRecursion(grammar parlex.Grammar) parlex.Grammar {
	nts := grammar.NonTerminals()
	op := &lrOp{
		in:   grammar,
		out:  Empty(),
		done: make(map[parlex.Symbol]bool, len(nts)),
	}
	for _, op.cur = range nts {
		op.hasDirect = false
		for _, prod := range grammar.Productions(op.cur) {
			op.safeAdd(prod)
		}
		if op.hasDirect {
			op.removeDirectLeftRecursion()
		}
		op.done[op.cur] = true
	}
	return op.out
}

func (op *lrOp) safeAdd(prod parlex.Production) {
	var first parlex.Symbol
	if len(prod) > 0 {
		first = prod[0]
	}
	if !op.done[first] {
		op.hasDirect = op.hasDirect || first == op.cur
		op.out.Add(op.cur, prod)
		return
	}

	tail := prod[1:]
	for _, lead := range op.in.Productions(first) {
		newProd := make(parlex.Production, len(lead)+len(tail))
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
	op.out.Add(newSym, nil)
}
