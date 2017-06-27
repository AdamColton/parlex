package parlex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type lx struct {
	k Symbol
	v string
}

func (l *lx) Kind() Symbol             { return l.k }
func (l *lx) Value() string            { return l.v }
func (l *lx) Pos() (line int, col int) { return 0, 0 }

func TestLexemeString(t *testing.T) {
	lxs := []Lexeme{
		&lx{k: "int", v: "1"},
		&lx{k: "op", v: "+"},
		&lx{k: "int", v: "2"},
		&lx{k: "op", v: "*"},
		&lx{k: "int", v: "3"},
	}

	expected := "[int: 1, op: +, int: 2, op: *, int: 3]"
	assert.Equal(t, expected, LexemeString(lxs...))
	assert.Equal(t, "int: 1", LexemeString(lxs[0]))
	assert.Equal(t, "", LexemeString())
}

type testGrammar struct {
	order       []Symbol
	productions map[Symbol]Productions
	cur         Symbol
}

func (tg *testGrammar) Productions(symbol Symbol) Productions {
	return tg.productions[symbol]
}

func (tg *testGrammar) NonTerminals() []Symbol {
	return tg.order
}

func (tg *testGrammar) new(symbol Symbol) {
	tg.cur = symbol
	tg.order = append(tg.order, symbol)
}

func (tg *testGrammar) add(symbols ...Symbol) {
	tg.productions[tg.cur] = append(tg.productions[tg.cur], Production(symbols))
}

func (tg *testGrammar) reset() {
	tg.productions = make(map[Symbol]Productions)
	tg.order = nil
}

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
