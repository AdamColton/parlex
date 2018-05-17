package tree

func Chain(r1, r2 Reduction) Reduction {
	if r1==nil{
		return r2
	}
	if r2==nil{
		return r1
	}
	return func(node *PN) {
		r1(node)
		r2(node)
	}
}

// Condition is a function that takes a node and return a bool. It can be used
// as the condition on an If in a Reduction Chain.
type Condition func(node *PN) bool

// If allows a Chain to perform conditional logic
func If(condition Condition, then, otherwise Reduction) Reduction {
	return func(node *PN) {
		if condition(node) {
			if then != nil {
				then(node)
			}
		} else {
			if otherwise != nil {
				otherwise(node)
			}
		}
	}
}

// If allows a Chain to perform conditional logic
func (r Reduction) If(condition Condition, then, otherwise Reduction) Reduction {
	return Chain(r, If(condition, then, otherwise))
}

// ChildIs returns true if the child at cIdx is of type kind
func ChildIs(cIdx int, kind string) Condition {
	return func(node *PN) bool {
		return node.ChildIs(cIdx, kind)
	}
}

// PromoteChild removes the node with the child at cIdx and replaces it's own
// lexeme with the value. The grandchildren are spliced into the replaced childs
// position. The cIdx value uses GetIdx.
func PromoteChild(cIdx int) Reduction {
	return func(node *PN) {
		node.PromoteChild(cIdx)
	}
}

// PromoteChildrenOf will remove the child at cIdx and splice in all it's
// children. If cIdx is negative, it will find the child relative to the end. If
// cIdx is out of bounds, no action will be taken.
func (r Reduction) PromoteChildrenOf(cIdx int) Reduction {
	return Chain(r, PromoteChildrenOf(cIdx))
}

// PromoteChildValue returns a Reduction that will replace the value of the node
// with the value from the child at cIdx and remove the child at cIdx. If cIdx
// is negative, it will find the child relative to the end. If cIdx is out of
// bounds, no action will be taken.
func (r Reduction) PromoteChildValue(cIdx int) Reduction {
	return Chain(r, PromoteChildValue(cIdx))
}

// PromoteGrandChildren will remove all the immediate children and replace them
// with the grand children.
func (r Reduction) PromoteGrandChildren() Reduction {
	return Chain(r, PromoteGrandChildren)
}

// RemoveAll removes all child that match any symbol in symbols
func (r Reduction) RemoveAll(symbols ...string) Reduction {
	return Chain(r, RemoveAll(symbols...))
}

// RemoveChild produces a Reduction for removing a child at cIdx.
func (r Reduction) RemoveChild(cIdx int) Reduction {
	return Chain(r, RemoveChild(cIdx))
}

// RemoveChildren produces a Reduction for removing a child at cIdx.
func (r Reduction) RemoveChildren(cIdxs ...int) Reduction {
	return Chain(r, RemoveChildren(cIdxs...))
}

// PromoteSingleChild fulfills Reduce. If the node has a single child, that
// child will be promoted to replace the node.
func (r Reduction) PromoteSingleChild() Reduction {
	return Chain(r, PromoteSingleChild)
}

// ReplaceWithChild replaces the node with the child at position cIdx, using
// GetIdx
func (r Reduction) ReplaceWithChild(cIdx int) Reduction {
	return Chain(r, ReplaceWithChild(cIdx))
}

// PromoteChild removes the node with the child at cIdx and replaces it's own
// lexeme with the value. The grandchildren are spliced into the replaced childs
// position. The cIdx value uses GetIdx.
func (r Reduction) PromoteChild(cIdx int) Reduction {
	return Chain(r, PromoteChild(cIdx))
}

// RemoveChild produces a Reduction for removing a child at cIdx.
func RemoveChild(cIdx int) Reduction {
	return func(node *PN) { node.RemoveChild(cIdx) }
}

// RemoveChildren produces a Reduction for removing a child at cIdx.
func RemoveChildren(cIdxs ...int) Reduction {
	return func(node *PN) { node.RemoveChildren(cIdxs...) }
}

// RemoveAll removes all child that match any symbol in symbols
func RemoveAll(symbols ...string) Reduction {
	return func(node *PN) {
		node.RemoveAll(symbols...)
	}
}

// PromoteChildValue returns a Reduction that will replace the value of the node
// with the value from the child at cIdx and remove the child at cIdx. If cIdx
// is negative, it will find the child relative to the end. If cIdx is out of
// bounds, no action will be taken.
func PromoteChildValue(cIdx int) Reduction {
	return func(node *PN) { node.PromoteChildValue(cIdx) }
}

// ReplaceWithChild replaces the node with the child at position cIdx, using
// GetIdx
func ReplaceWithChild(cIdx int) Reduction {
	return func(node *PN) { node.ReplaceWithChild(cIdx) }
}

// PromoteChildrenOf will remove the child at cIdx and splice in all it's
// children. If cIdx is negative, it will find the child relative to the end. If
// cIdx is out of bounds, no action will be taken.
func PromoteChildrenOf(cIdx int) Reduction {
	return func(node *PN) {
		node.PromoteChildrenOf(cIdx)
	}
}

// PromoteGrandChildren will remove all the immediate children and replace them
// with the grand children.
func PromoteGrandChildren(node *PN) {
	node.PromoteGrandChildren()
}

// PromoteSingleChild fulfills Reduce. If the node has a single child, that
// child will be promoted to replace the node.
func PromoteSingleChild(node *PN) {
	node.PromoteSingleChild()
}