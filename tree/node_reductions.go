package tree

import (
	"github.com/adamcolton/parlex/lexeme"
)

// ReplaceWithChild replaces the node with the child at cIdx.
func (p *PN) ReplaceWithChild(cIdx int) bool {
	cIdx, _, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}
	*p = *(p.C[cIdx])
	return true
}

// RemoveAll removes all child that match any symbol in symbols
func (p *PN) RemoveAll(symbols ...string) {
	for i := 0; i < len(p.C); i++ {
		if p.ChildAt(i, symbols...) {
			p.RemoveChild(i)
			i--
		}
	}
}

// GetIdx takes a child index and performs a conversion and check on it. If the
// index is negative, it will convert it to a positive position relative to the
// end (so -1 is the last child). The second int returned is the number of
// children. And the last value is bool indicating if cIdx is between 0 and len.
func (p *PN) GetIdx(cIdx int) (int, int, bool) {
	l := len(p.C)
	if cIdx < 0 {
		cIdx = l + cIdx
	}
	return cIdx, l, cIdx >= 0 && cIdx < l
}

// RemoveChildren calls RemoveChild repeatedly for each index given. Note that
// the relative positions may change so if you wanted to remove what were
// initially indexes 1, 3 and 5, you would need to either call
// RemoveChildren(5,3,1) or RemoveChildren(1,2,3)
func (p *PN) RemoveChildren(cIdxs ...int) {
	for _, cIdx := range cIdxs {
		p.RemoveChild(cIdx)
	}
}

// RemoveChild remove the child at cIdx. If cIdx is negative, it will find the
// child relative to the end. If cIdx is out of bounds, no action will be taken.
func (p *PN) RemoveChild(cIdx int) bool {
	cIdx, l, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}
	if cIdx == l-1 {
		p.C = p.C[:cIdx]
	} else {
		p.C = append(p.C[0:cIdx], p.C[cIdx+1:]...)
	}
	return true
}

// PromoteSingleChild ; if the node has a single child, that child will be
// promoted to replace the node.
func (p *PN) PromoteSingleChild() bool {
	if len(p.C) == 1 {
		p.PromoteChild(0)
		return true
	}
	return false
}

// PromoteChildValue returns a Reduction that will replace the value of the node
// with the value from the child at cIdx and remove the child at cIdx. If cIdx
// is negative, it will find the child relative to the end. If cIdx is out of
// bounds, no action will be taken.
func (p *PN) PromoteChildValue(cIdx int) {
	l := len(p.C)
	if cIdx < 0 {
		cIdx = l + cIdx
	}
	if cIdx >= 0 && l > cIdx {
		ch := p.C[cIdx]
		p.Lexeme = lexeme.New(p.Kind()).Set(ch.Value()).At(ch.Pos())
	}
	p.RemoveChild(cIdx)
}

// ChildIs returns true if the child at cIdx is of type kind
func (p *PN) ChildIs(cIdx int, kind string) bool {
	cIdx, _, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}
	return p.Child(cIdx).Kind().String() == kind
}

// PromoteChildrenOf will remove the child at cIdx and splice in all it's
// children. If cIdx is negative, it will find the child relative to the end. If
// cIdx is out of bounds, no action will be taken.
func (p *PN) PromoteChildrenOf(cIdx int) bool {
	cIdx, l, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}
	if cIdx == l-1 {
		p.C = append(p.C[:cIdx], p.C[cIdx].C...)
	} else {
		p.C = append(p.C[:cIdx], append(p.C[cIdx].C, p.C[cIdx+1:]...)...)
	}
	return true
}

// PromoteGrandChildren will remove all the immediate children and replace them
// with the grand children.
func (p *PN) PromoteGrandChildren() {
	ct := 0
	for _, child := range p.C {
		ct += len(child.C)
	}
	newChildren := make([]*PN, 0, ct)
	for _, child := range p.C {
		for _, grandChild := range child.C {
			grandChild.P = p
		}
		newChildren = append(newChildren, child.C...)
	}
	p.C = newChildren
}

// PromoteChild removes the node with the child at cIdx and replaces it's own
// lexeme with the value. The grandchildren are spliced into the replaced childs
// position. The cIdx value uses GetIdx.
func (p *PN) PromoteChild(cIdx int) bool {
	cIdx, l, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}

	p.Lexeme = p.C[cIdx].Lexeme
	tail := p.C[cIdx].C
	if cIdx+1 < l {
		tail = append(tail, p.C[cIdx+1:]...)
	}
	p.C = append(p.C[0:cIdx], tail...)
	for _, c := range p.C {
		c.P = p
	}
	return true
}
