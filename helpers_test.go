package parlex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type symbol string

func (s symbol) String() string { return string(s) }

type production []symbol

func (p production) Symbols() int {
	return len(p)
}
func (p production) Symbol(i int) Symbol {
	if i < len(p) {
		return p[i]
	}
	return nil
}

func (p production) Iter() *ProductionIterator {
	return &ProductionIterator{
		Production: p,
	}
}

type productions []production

func (p productions) Productions() int {
	return len(p)
}

func (p productions) Production(i int) Production {
	if i < len(p) {
		return p[i]
	}
	return nil
}

func (p productions) Iter() *ProductionsIterator {
	return &ProductionsIterator{
		Productions: p,
	}
}

type lx struct {
	k symbol
	v string
}

func (l *lx) Kind() Symbol             { return l.k }
func (l *lx) Value() string            { return l.v }
func (l *lx) Pos() (line int, col int) { return 0, 0 }

type testGrammar struct {
	order       []Symbol
	productions map[symbol]productions
	cur         symbol
}

func (tg *testGrammar) Productions(s Symbol) Productions {
	return tg.productions[s.(symbol)]
}

func (tg *testGrammar) NonTerminals() []Symbol {
	return tg.order
}

func (tg *testGrammar) new(s symbol) {
	tg.cur = s
	tg.order = append(tg.order, s)
}

func (tg *testGrammar) add(symbols ...symbol) {
	tg.productions[tg.cur] = append(tg.productions[tg.cur], production(symbols))
}

func (tg *testGrammar) reset() {
	tg.productions = make(map[symbol]productions)
	tg.order = nil
}

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
