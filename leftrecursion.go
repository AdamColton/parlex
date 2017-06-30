package parlex

// IsLeftRecursive will check if a grammar is left recursive.
func IsLeftRecursive(grammar Grammar) bool {
	checked := make(map[string]bool)
	for _, nt := range grammar.NonTerminals() {
		isLR, _ := checkLeftRecursion(grammar, nt, checked, make(map[string]bool))
		if isLR {
			return true
		}
	}
	return false
}

func checkLeftRecursion(g Grammar, s Symbol, checked, stack map[string]bool) (bool, bool) {
	if checked[s.String()] {
		return false, false
	}

	prods := g.Productions(s)
	if prods == nil {
		// s is a terminal
		checked[s.String()] = true // add to checked to avoid the lookup in the future
		return false, false
	}

	if stack[s.String()] {
		return true, false
	}
	stack[s.String()] = true

	retCheckNext := false
	for i := prods.Iter(); i.Next(); {
		if i.Production.Symbols() == 0 {
			retCheckNext = true
			continue
		}
		for j, isLR, checkNext := i.Iter(), false, true; j.Next() && checkNext; {
			isLR, checkNext = checkLeftRecursion(g, j.Symbol, checked, stack)
			if isLR {
				return true, false
			}
		}
	}

	checked[s.String()] = true
	stack[s.String()] = false
	return false, retCheckNext
}
