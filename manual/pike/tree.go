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
			b.wait()
			b.match(rs[0])
		} else {
			b.wait()
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
	b.startCounter()
	loc := b.loc()
	db := b.defer_branch()
	mn.child.build(b)
	b.incCounter()
	b.jump(loc)
	db()
	b.ck_gte_c(mn.val)
	b.closeCounter()
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
	b.startCounter()
	loc := b.loc()
	db := b.defer_branch()
	b.ck_lt_c(mn.val)
	mn.child.build(b)
	b.incCounter()
	b.jump(loc)
	db()
	b.closeCounter()
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
	b.startCounter()
	loc := b.loc()
	db := b.defer_branch()
	b.ck_lt_c(mmn.max)
	mmn.child.build(b)
	b.incCounter()
	b.jump(loc)
	db()
	b.ck_gte_c(mmn.min)
	b.closeCounter()
}

func (mmn minmaxNode) String() string {
	return mmn.child.String() + fmt.Sprintf("{%d,%d}", mmn.min, mmn.max)
}

func (mmn minmaxNode) Tree(ind string) string {
	return fmt.Sprintf("%sMinMax(%d, %d) {\n%s%s}\n", ind, mmn.min, mmn.max, mmn.child.Tree(ind+"\t"), ind)
}

type anyNode struct{}

func (an anyNode) build(b *builder) {
	b.wait()
}

func (an anyNode) String() string {
	return "."
}

func (an anyNode) Tree(ind string) string {
	return fmt.Sprintf("%sAny\n", ind)
}

type oneOrMoreNode struct {
	child node
}

func (omn oneOrMoreNode) build(b *builder) {
	loc := b.loc()
	omn.child.build(b)
	b.branch(loc)
}

func (omn oneOrMoreNode) String() string {
	return omn.child.String() + "+"
}

func (omn oneOrMoreNode) Tree(ind string) string {
	return fmt.Sprintf("%sOneOrMore{\n%s%s}\n", ind, omn.child.Tree(ind+"\t"), ind)
}

type oneOrZeroNode struct {
	child node
}

func (ozn oneOrZeroNode) build(b *builder) {
	db := b.defer_branch()
	ozn.child.build(b)
	db()
}

func (ozn oneOrZeroNode) String() string {
	return ozn.child.String() + "?"
}

func (ozn oneOrZeroNode) Tree(ind string) string {
	return fmt.Sprintf("%sOneOrZero{\n%s%s}\n", ind, ozn.child.Tree(ind+"\t"), ind)
}
