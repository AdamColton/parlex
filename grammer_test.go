package parlex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRecursive(t *testing.T) {
	g := &testGrammar{}
	g.reset()
	g.new("A")
	g.add("B", "C")
	g.add("x")
	g.new("B")
	g.add("y")
	g.add("w")
	g.new("C")
	g.add("z")
	g.add("A")
	assert.False(t, IsLeftRecursive(g))

	g.new("A")
	g.add("B", "C")
	g.add("x")
	g.new("B")
	g.add("Y")
	g.new("C")
	g.add("z")
	g.new("Y")
	g.add("A")
	assert.True(t, IsLeftRecursive(g))

	g.new("A")
	g.add("B", "C")
	g.add("x")
	g.new("B")
	g.add("w")
	g.add()
	g.new("C")
	g.add("A")
	g.add("a")
	assert.True(t, IsLeftRecursive(g))
}

func TestMustGrammar(t *testing.T) {
	g := &testGrammar{}
	assert.Equal(t, g, MustGrammar(g, nil))

	defer func() {
		assert.Equal(t, testErr, recover())
	}()
	MustGrammar(g, testErr)
}
