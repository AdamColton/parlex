package pike

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchWord(t *testing.T) {
	// ca*|ot
	b := newBuilder()

	b.wait()
	b.match('c')

	oBranch := b.defer_branch()
	b.set_rv(0, 1)
	loc := b.loc()
	b.wait()
	b.match('a')
	b.inc(1)
	b.branch(loc)
	jmp := b.defer_jump()

	oBranch()
	b.set_rv(0, 2)
	loc = b.loc()
	b.wait()
	b.match('o')
	b.inc(1)
	b.branch(loc)

	jmp()
	b.wait()
	b.match('t')
	b.accept()
	b.stop()

	p := b.close()

	op := p.run("cat")
	assert.Equal(t, 3, op.best)
	assert.Equal(t, uint32(1), op.bestState.workingState().readUint32(0))
	assert.Equal(t, uint32(1), op.bestState.workingState().readUint32(1))

	op = p.run("cot")
	assert.Equal(t, 3, op.best)
	assert.Equal(t, uint32(2), op.bestState.workingState().readUint32(0))
	assert.Equal(t, uint32(1), op.bestState.workingState().readUint32(1))

	op = p.run("caaat")
	assert.Equal(t, 5, op.best)
	assert.Equal(t, uint32(1), op.bestState.workingState().readUint32(0))
	assert.Equal(t, uint32(3), op.bestState.workingState().readUint32(1))
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

	op := p.run("cat")
	assert.Equal(t, 3, op.best)
	g := op.groups[1][0]
	assert.Equal(t, 1, g[0])
	assert.Equal(t, 2, g[1])

	op = p.run("caaat")
	assert.Equal(t, 5, op.best)
	g = op.groups[1][0]
	assert.Equal(t, 1, g[0])
	assert.Equal(t, 4, g[1])
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
	g := op.groups[1][0]
	assert.Equal(t, 0, g[0])
	assert.Equal(t, 3, g[1])

	op = p.run("123.4")
	assert.Equal(t, 5, op.best)
	g = op.groups[1][0]
	assert.Equal(t, 0, g[0])
	assert.Equal(t, 3, g[1])
	g = op.groups[2][0]
	assert.Equal(t, 4, g[0])
	assert.Equal(t, 5, g[1])
}
