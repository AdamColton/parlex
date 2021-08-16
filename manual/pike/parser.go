package pike

import (
	"strconv"
)

type parseOp struct {
	reg, group uint32
	cur        node
	stack      []node
	*Reader
}

func parse(str string) node {
	op := &parseOp{
		cur:    rootNode{},
		Reader: NewStringReader(str),
	}
	for op.Inc(); !op.Done(); op.Inc() {
		switch op.R {
		case '*':
			op.modLast(func(n node) node {
				return kleeneStarNode{n}
			})
		case '+':
			op.modLast(func(n node) node {
				return oneOrMoreNode{n}
			})
		case '{':
			op.minmax()
		case '|':
			op.modLast(func(n node) node {
				return orNode{n, nil}
			})
		case '(':
			op.startGroup()
		case ')':
			for {
				_, isGroup := op.cur.(groupNode)
				if isGroup {
					break
				}
				op.pop()
			}
			op.pop()
		case '.':
			op.append(anyNode{})
		default:
			op.append(matchNode{}.match(op.R))
		}
	}
	for len(op.stack) > 0 {
		op.pop()
	}
	return op.cur
}

func (op *parseOp) startGroup() {
	g := groupNode{
		idx: op.group,
	}
	op.group++
	switch c := op.cur.(type) {
	case rootNode:
		break
	case manyNodes:
		op.cur = append(c, g)
	default:
		op.cur = manyNodes{op.cur, g}
	}
	op.push(g)
	g.child = manyNodes{}
	op.push(manyNodes{})
}

func (op *parseOp) minmax() {
	var rng [2][2]int
	op.Inc() // consume {
	op.Consume(' ')
	rng[0][0] = op.Idx
	op.ConsumeRange('0', '9')
	rng[0][1] = op.Idx
	if op.R != ',' {
		panic("bad range match (no comma)")
	}
	op.Inc() // consume ,
	op.Consume(' ')
	rng[1][0] = op.Idx
	op.ConsumeRange('0', '9')
	rng[1][1] = op.Idx
	if op.R != '}' {
		panic("bad range match (no lcb)")
	}

	if rng[0][0] == rng[0][1] {
		s, e := rng[1][0], rng[1][1]
		str := string(op.input[s:e])
		max, _ := strconv.Atoi(str)
		op.modLast(func(n node) node {
			r := op.reg
			op.reg++
			return maxNode{
				reg:   r,
				val:   uint32(max),
				child: n,
			}

		})
	} else if rng[1][0] == rng[1][1] {
		s, e := rng[0][0], rng[0][1]
		str := string(op.input[s:e])
		min, _ := strconv.Atoi(str)
		op.modLast(func(n node) node {
			r := op.reg
			op.reg++
			return minNode{
				reg:   r,
				val:   uint32(min),
				child: n,
			}

		})
	} else {
		s, e := rng[0][0], rng[0][1]
		str := string(op.input[s:e])
		min, _ := strconv.Atoi(str)
		s, e = rng[1][0], rng[1][1]
		str = string(op.input[s:e])
		max, _ := strconv.Atoi(str)
		op.modLast(func(n node) node {
			r := op.reg
			op.reg++
			return minmaxNode{
				reg:   r,
				min:   uint32(min),
				max:   uint32(max),
				child: n,
			}

		})
	}
}

func (op *parseOp) append(n node) {
	switch c := op.cur.(type) {
	case rootNode, groupNode:
		op.push(n)
	case manyNodes:
		op.cur = append(c, n)
	case orNode:
		if c[1] == nil {
			c[1] = n
			op.cur = c
		} else {
			op.cur = manyNodes{c, n}
		}
	default:
		op.cur = manyNodes{c, n}
	}
}

func (op *parseOp) push(n node) {
	op.stack = append(op.stack, op.cur)
	op.cur = n
}

// pop moves up to the next node that can accept a child
func (op *parseOp) pop() {
	n := op.cur
	ln := len(op.stack) - 1
	op.cur = op.stack[ln]
	op.stack = op.stack[:ln]

	switch c := op.cur.(type) {
	case rootNode:
		c.child = n
		op.cur = c
	case manyNodes:
		if m, isMany := n.(manyNodes); isMany {
			c = append(c[:len(c)-1], m...)
		} else {
			c[len(c)-1] = n
		}
		op.cur = c
	case groupNode:
		if m, isMany := n.(manyNodes); isMany && len(m) == 1 {
			c.child = m[0]
		} else {
			c.child = n
		}
		op.cur = c
	case matchNode:
		op.cur = manyNodes{op.cur, n}
	default:
		panic(c)
	}
}

func (op *parseOp) modLast(fn func(node) node) {
	switch c := op.cur.(type) {
	case rootNode:
		panic("cannot mod root")
	case manyNodes:
		ln := len(c) - 1
		c[ln] = fn(c[ln])
		op.cur = c
		op.push(c[ln])
	case orNode:
		if c[1] == nil {
			panic("bad or mod")
		}
		c[1] = fn(c[1])
		op.cur = c
	default:
		op.cur = fn(op.cur)
	}
}
