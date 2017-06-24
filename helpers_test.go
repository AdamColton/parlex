package parlex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexemeString(t *testing.T) {
	lxs := []Lexeme{
		&L{K: "int", V: "1"},
		&L{K: "op", V: "+"},
		&L{K: "int", V: "2"},
		&L{K: "op", V: "*"},
		&L{K: "int", V: "3"},
	}

	expected := "[int: 1, op: +, int: 2, op: *, int: 3]"
	assert.Equal(t, expected, LexemeString(lxs...))
	assert.Equal(t, "int: 1", LexemeString(lxs[0]))
	assert.Equal(t, "", LexemeString())
}
