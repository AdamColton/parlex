package parlex

// IsLeftRecursive will check if a grammar is left recursive.
func IsLeftRecursive(grammar Grammar) bool {
	op := &lrOp{
		Grammar: grammar,
		checked: make(map[string]bool),
	}
	for _, nt := range grammar.NonTerminals() {
		op.stack = make(map[string]bool)
		isLR, _ := op.check(nt)
		if isLR {
			return true
		}
	}
	return false
}

type lrOp struct {
	Grammar
	checked, stack map[string]bool
}

func (op *lrOp) check(s Symbol) (isLR, checkNext bool) {
	if op.checked[s.String()] {
		return false, false
	}

	prods := op.Productions(s)
	if prods == nil {
		// s is a terminal
		op.checked[s.String()] = true // add to checked to avoid the lookup in the future
		return false, false
	}

	if op.stack[s.String()] {
		return true, false
	}
	op.stack[s.String()] = true

	retCheckNext := false
	for i := prods.Iter(); i.Next(); {
		if i.Production.Symbols() == 0 {
			retCheckNext = true
			continue
		}
		for j, isLR, checkNext := i.Iter(), false, true; j.Next() && checkNext; {
			isLR, checkNext = op.check(j.Symbol)
			if isLR {
				return true, false
			}
		}
	}

	op.checked[s.String()] = true
	op.stack[s.String()] = false
	return false, retCheckNext
}
