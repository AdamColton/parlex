package parlex

// IsLeftRecursive will check if a grammar is left recursive.
func IsLeftRecursive(grammar Grammar) bool {
	checked := make(map[Symbol]bool)
	for _, nt := range grammar.NonTerminals() {
		isLR, _ := checkLeftRecursion(grammar, nt, checked, make(map[Symbol]bool))
		if isLR {
			return true
		}
	}
	return false
}

func checkLeftRecursion(g Grammar, s Symbol, checked, stack map[Symbol]bool) (bool, bool) {
	if checked[s] {
		return false, false
	}

	prods := g.Productions(s)
	if prods == nil {
		// s is a terminal
		checked[s] = true // add to checked to avoid the lookup in the future
		return false, false
	}

	if stack[s] {
		return true, false
	}
	stack[s] = true

	retCheckNext := false
	var prod Production
	for i := 0; i < prods.Productions(); i++ {
		prod = prods.Production(i)
		if prod.Symbols() == 0 {
			retCheckNext = true
			continue
		}
		for i, isLR, checkNext := 0, false, true; i < prod.Symbols() && checkNext; i++ {
			isLR, checkNext = checkLeftRecursion(g, prod.Symbol(i), checked, stack)
			if isLR {
				return true, false
			}
		}
	}

	checked[s] = true
	stack[s] = false
	return false, retCheckNext
}
