package pike

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchWord(t *testing.T) {
	// cat
	b := newBuilder()

	b.wait()
	b.match('c')
	b.wait()
	b.match('a')
	b.wait()
	b.match('t')

	b.accept()
	b.stop()

	p := b.close()

	op := p.run("cat")
	assert.Equal(t, 3, op.best)

	op = p.run("cot")
	assert.Equal(t, -1, op.best)
}

func TestGroup(t *testing.T) {
	b := newBuilder()
	b.wait()
	b.match('c')
	b.startGroup(1)
	loc := b.loc()
	b.wait()
	b.match('a')
	b.branch(loc)
	b.closeGroup()
	b.wait()
	b.match('t')
	b.accept()
	b.stop()

	p := b.close()

	// op := p.run("cat")
	// assert.Equal(t, 3, op.best)
	// gs := op.groupMap.toMap(op.bestGroups)
	// g := gs[1][0]
	// assert.Equal(t, 1, g[0])
	// assert.Equal(t, 2, g[1])

	op := p.run("caaat")
	assert.Equal(t, 5, op.best)
	// gs = op.groupMap.toMap(op.bestGroups)
	// g = gs[1][0]
	// assert.Equal(t, 1, g[0])
	// assert.Equal(t, 4, g[1])
}

func TestMatchRange(t *testing.T) {
	b := newBuilder()
	b.startGroup(1)
	loc := b.loc()
	b.wait()
	b.match_range('0', '9')
	b.branch(loc)
	b.closeGroup()
	b.accept()
	b.wait()
	b.match('.')
	b.startGroup(2)
	loc = b.loc()
	b.wait()
	b.match_range('0', '9')
	b.branch(loc)
	b.closeGroup()
	b.accept()
	b.stop()

	p := b.close()

	op := p.run("123")
	assert.Equal(t, 3, op.best)
	gs := op.groupMap.toMap(op.bestGroups)
	g := gs[1][0]
	assert.Equal(t, 0, g[0])
	assert.Equal(t, 3, g[1])

	op = p.run("123.4")
	assert.Equal(t, 5, op.best)
	gs = op.groupMap.toMap(op.bestGroups)
	g = gs[1][0]
	assert.Equal(t, 0, g[0])
	assert.Equal(t, 3, g[1])
	g = gs[2][0]
	assert.Equal(t, 4, g[0])
	assert.Equal(t, 5, g[1])
}

func TestCounter(t *testing.T) {
	b := newBuilder()
	b.wait()
	b.match('c')
	b.startCounter()
	loc := b.loc()
	db := b.defer_branch()
	b.ck_lt_c(3)
	b.wait()
	b.match('a')
	b.incCounter()
	b.jump(loc)
	db()
	b.ck_gte_c(2)
	b.wait()
	b.match('t')
	b.accept()
	b.stop()

	p := b.close()

	op := p.run("cat")
	assert.Equal(t, -1, op.best)

	op = p.run("caat")
	assert.Equal(t, 4, op.best)

	op = p.run("caaat")
	assert.Equal(t, 5, op.best)

	op = p.run("caaaat")
	assert.Equal(t, -1, op.best)
}
