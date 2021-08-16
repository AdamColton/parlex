package pike

import "fmt"

type node interface {
	build(*builder)
	String() string
	Tree(string) string
}

func buildTree(n node) *prog {
	b := newBuilder()
	n.build(b)
	return b.close()
}

type rootNode struct {
	child node
}

func (rn rootNode) build(b *builder) {
	rn.child.build(b)
	b.accept()
	b.stop()
}

func (rn rootNode) String() string {
	return rn.child.String()
}

func (rn rootNode) Tree(ind string) string {
	return rn.child.Tree(ind)
}

type manyNodes []node

func (mn manyNodes) build(b *builder) {
	for _, n := range mn {
		n.build(b)
	}
}

func (mn manyNodes) String() string {
	out := ""
	for _, n := range mn {
		out += n.String()
	}
	return out
}

func (mn manyNodes) Tree(ind string) string {
	out := fmt.Sprintf("%s{\n", ind)
	sub := ind + "\t"
	for _, n := range mn {
		out += n.Tree(sub)
	}
	out += fmt.Sprintf("%s}\n", ind)
	return out
}

type matchNode [][]rune

func (mn matchNode) match(r rune) matchNode {
	return append(mn, []rune{r})
}

func (mn matchNode) matchRange(start, end rune) matchNode {
	return append(mn, []rune{start, end})
}

func (mn matchNode) build(b *builder) {
	for _, rs := range mn {
		if len(rs) == 1 {
			b.match(rs[0])
		} else {
			b.match_range(rs[0], rs[1])
		}
	}
}

func (mn matchNode) String() string {
	out := ""
	for _, rs := range mn {
		if len(rs) == 1 {
			out += string(rs[0])
		} else {
			out += fmt.Sprintf("[%v-%v]", rs[0], rs[1])
		}
	}
	if len(mn) > 1 {
		out = "(:" + out + ")"
	}
	return out
}

func (mn matchNode) Tree(ind string) string {
	out := ""
	sub := ind
	if len(mn) > 1 {
		out = fmt.Sprintf("%s{\n", ind)
		sub += "\t"
	}
	for _, rs := range mn {
		if len(rs) == 1 {
			out += fmt.Sprintf("%sMatch(%s)\n", sub, string(rs[0]))
		} else {
			out += fmt.Sprintf("%sMatch(%s-%s)\n", sub, string(rs[0]), string(rs[1]))
		}
	}
	if len(mn) > 1 {
		out += fmt.Sprintf("%s}\n", ind)
	}
	return out
}

type orNode [2]node

func (on orNode) build(b *builder) {
	db := b.defer_branch()
	on[0].build(b)
	dj := b.defer_jump()
	db()
	on[1].build(b)
	dj()
}

func (on orNode) String() string {
	return on[0].String() + "|" + on[1].String()
}

func (on orNode) Tree(ind string) string {
	sub := ind + "\t"
	return fmt.Sprintf("%sOr {\n%s%s%s}\n", ind, on[0].Tree(sub), on[1].Tree(sub), ind)
}

type kleeneStarNode struct {
	child node
}

func (ksn kleeneStarNode) build(b *builder) {
	loc := b.loc()
	db := b.defer_branch()
	ksn.child.build(b)
	b.jump(loc)
	db()
}

func (ksn kleeneStarNode) String() string {
	return ksn.child.String() + "*"
}

func (ksn kleeneStarNode) Tree(ind string) string {
	return fmt.Sprintf("%sMany {\n%s%s}\n", ind, ksn.child.Tree(ind+"\t"), ind)
}

type groupNode struct {
	idx   uint32
	child node
}

func (gn groupNode) build(b *builder) {
	b.startGroup(gn.idx)
	gn.child.build(b)
	b.closeGroup()
}

func (gn groupNode) String() string {
	return "(" + gn.child.String() + ")"
}

func (gn groupNode) Tree(ind string) string {
	return fmt.Sprintf("%sGroup {\n%s%s}\n", ind, gn.child.Tree(ind+"\t"), ind)
}

type minNode struct {
	reg, val uint32
	child    node
}

func (mn minNode) build(b *builder) {
	b.set_rv(mn.reg, 0)
	loc := b.loc()
	db := b.defer_branch()
	mn.child.build(b)
	b.inc(mn.reg)
	b.jump(loc)
	db()
	b.ck_gte_rv(mn.reg, mn.val)
}

func (mn minNode) String() string {
	return mn.child.String() + fmt.Sprintf("{%d,}", mn.val)
}

func (mn minNode) Tree(ind string) string {
	return fmt.Sprintf("%sMin(%d) {\n%s%s}\n", ind, mn.val, mn.child.Tree(ind+"\t"), ind)
}

type maxNode struct {
	reg, val uint32
	child    node
}

func (mn maxNode) build(b *builder) {
	b.set_rv(mn.reg, 0)
	loc := b.loc()
	db := b.defer_branch()
	b.ck_lt_rv(mn.reg, mn.val)
	mn.child.build(b)
	b.inc(mn.reg)
	b.jump(loc)
	db()
}

func (mn maxNode) String() string {
	return mn.child.String() + fmt.Sprintf("{,%d}", mn.val)
}

func (mn maxNode) Tree(ind string) string {
	return fmt.Sprintf("%sMin(%d) {\n%s%s}\n", ind, mn.val, mn.child.Tree(ind+"\t"), ind)
}

type minmaxNode struct {
	reg, min, max uint32
	child         node
}

func (mmn minmaxNode) build(b *builder) {
	b.set_rv(mmn.reg, 0)
	loc := b.loc()
	db := b.defer_branch()
	b.ck_lt_rv(mmn.reg, mmn.max)
	mmn.child.build(b)
	b.inc(mmn.reg)
	b.jump(loc)
	db()
	b.ck_gte_rv(mmn.reg, mmn.min)
}

func (mmn minmaxNode) String() string {
	return mmn.child.String() + fmt.Sprintf("{%d,%d}", mmn.min, mmn.max)
}

func (mmn minmaxNode) Tree(ind string) string {
	return fmt.Sprintf("%sMin(%d, %d) {\n%s%s}\n", ind, mmn.min, mmn.max, mmn.child.Tree(ind+"\t"), ind)
}
