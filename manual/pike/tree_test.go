package pike

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeMatch(t *testing.T) {
	n := rootNode{
		matchNode{}.match('c').
			match('a').
			match('t'),
	}
	assert.Equal(t, "(:cat)", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, 3, op.best)
}

func TestTreeOr(t *testing.T) {
	n := rootNode{
		manyNodes{
			matchNode{}.match('c'),
			orNode{
				matchNode{}.match('a'),
				matchNode{}.match('o'),
			},
			matchNode{}.match('t'),
		},
	}
	assert.Equal(t, "ca|ot", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, 3, op.best)

	op = p.run("cot")
	assert.Equal(t, 3, op.best)

	op = p.run("cut")
	assert.Equal(t, -1, op.best)
}

func TestTreeKleeneStar(t *testing.T) {
	n := rootNode{
		manyNodes{
			matchNode{}.match('c'),
			orNode{
				kleeneStarNode{matchNode{}.match('a')},
				matchNode{}.match('o'),
			},
			matchNode{}.match('t'),
		},
	}
	assert.Equal(t, "ca*|ot", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, 3, op.best)

	op = p.run("cot")
	assert.Equal(t, 3, op.best)

	op = p.run("cut")
	assert.Equal(t, -1, op.best)

	op = p.run("caat")
	assert.Equal(t, 4, op.best)
	op = p.run("caaat")
	assert.Equal(t, 5, op.best)

	op = p.run("coot")
	assert.Equal(t, -1, op.best)
}

func TestTreeGroup(t *testing.T) {
	n := rootNode{
		manyNodes{
			matchNode{}.match('c'),
			groupNode{
				idx: 1,
				child: orNode{
					kleeneStarNode{matchNode{}.match('a')},
					matchNode{}.match('o'),
				},
			},
			matchNode{}.match('t'),
		},
	}
	assert.Equal(t, "c(a*|o)t", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, 3, op.best)
	g := op.groups[1][0]
	assert.Equal(t, 1, g[0])
	assert.Equal(t, 2, g[1])

	op = p.run("cot")
	assert.Equal(t, 3, op.best)
	g = op.groups[1][0]
	assert.Equal(t, 1, g[0])
	assert.Equal(t, 2, g[1])

	op = p.run("caat")
	assert.Equal(t, 4, op.best)
	g = op.groups[1][0]
	assert.Equal(t, 1, g[0])
	assert.Equal(t, 3, g[1])

	op = p.run("caaat")
	assert.Equal(t, 5, op.best)
	g = op.groups[1][0]
	assert.Equal(t, 1, g[0])
	assert.Equal(t, 4, g[1])
}

func TestTreeMin(t *testing.T) {
	n := rootNode{
		manyNodes{
			matchNode{}.match('c'),
			minNode{
				reg:   0,
				val:   3,
				child: matchNode{}.match('a'),
			},
			matchNode{}.match('t'),
		},
	}
	assert.Equal(t, "ca{3,}t", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, -1, op.best)

	op = p.run("caat")
	assert.Equal(t, -1, op.best)

	op = p.run("caaat")
	assert.Equal(t, 5, op.best)

	op = p.run("caaaat")
	assert.Equal(t, 6, op.best)
}

func TestTreeMax(t *testing.T) {
	n := rootNode{
		manyNodes{
			matchNode{}.match('c'),
			maxNode{
				reg:   0,
				val:   3,
				child: matchNode{}.match('a'),
			},
			matchNode{}.match('t'),
		},
	}
	assert.Equal(t, "ca{,3}t", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, 3, op.best)

	op = p.run("caat")
	assert.Equal(t, 4, op.best)

	op = p.run("caaat")
	assert.Equal(t, 5, op.best)

	op = p.run("caaaat")
	assert.Equal(t, -1, op.best)
}

func TestTreeMinMax(t *testing.T) {
	n := rootNode{
		manyNodes{
			matchNode{}.match('c'),
			minmaxNode{
				reg:   0,
				min:   2,
				max:   3,
				child: matchNode{}.match('a'),
			},
			matchNode{}.match('t'),
		},
	}
	assert.Equal(t, "ca{2,3}t", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, -1, op.best)

	op = p.run("caat")
	assert.Equal(t, 4, op.best)

	op = p.run("caaat")
	assert.Equal(t, 5, op.best)

	op = p.run("caaaat")
	assert.Equal(t, -1, op.best)
}

func TestTreeAny(t *testing.T) {
	n := rootNode{
		manyNodes{
			matchNode{}.match('c'),
			anyNode{},
			matchNode{}.match('t'),
		},
	}
	assert.Equal(t, "c.t", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, 3, op.best)

	op = p.run("cot")
	assert.Equal(t, 3, op.best)

	op = p.run("caat")
	assert.Equal(t, -1, op.best)
}

func TestTreeOneOrMore(t *testing.T) {
	n := rootNode{
		manyNodes{
			matchNode{}.match('c'),
			oneOrMoreNode{
				matchNode{}.match('a'),
			},
			matchNode{}.match('t'),
		},
	}
	assert.Equal(t, "ca+t", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, 3, op.best)

	op = p.run("ct")
	assert.Equal(t, -1, op.best)

	op = p.run("caat")
	assert.Equal(t, 4, op.best)

	op = p.run("caaat")
	assert.Equal(t, 5, op.best)
}

func TestTreeOneOrZero(t *testing.T) {
	n := rootNode{
		manyNodes{
			matchNode{}.match('c'),
			oneOrZeroNode{
				matchNode{}.match('a'),
			},
			matchNode{}.match('t'),
		},
	}
	assert.Equal(t, "ca?t", n.String())

	p := buildTree(n)

	op := p.run("cat")
	assert.Equal(t, 3, op.best)

	op = p.run("ct")
	assert.Equal(t, 2, op.best)

	op = p.run("caat")
	assert.Equal(t, -1, op.best)

	op = p.run("caaat")
	assert.Equal(t, -1, op.best)
}
