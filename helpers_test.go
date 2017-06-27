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
