package parlex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type lxErr struct {
	*lx
}

func (l *lxErr) Error() string {
	return "ERROR"
}

func TestLexErrors(t *testing.T) {
	lxs := []Lexeme{
		&lx{k: "int", v: "1"},
		&lx{k: "op", v: "+"},
		&lxErr{&lx{k: "int", v: "2"}},
		&lx{k: "op", v: "*"},
		&lx{k: "int", v: "3"},
	}

	errs := LexErrors(lxs)
	if assert.Len(t, errs, 1) {
		assert.Equal(t, "2", errs[0].Value())
	}
}

func TestLexemeList(t *testing.T) {
	lxs := []Lexeme{
		&lx{k: "int", v: "1"},
		&lx{k: "op", v: "+"},
		&lxErr{&lx{k: "int", v: "2"}},
		&lx{k: "op", v: "*", l: 12, c: 4},
		&lx{k: "int", v: "3"},
	}

	str := LexemeList(lxs)

	assert.Equal(t, "int: \"1\"\nop: \"+\"\nint: \"2\"\nop: \"*\" (12, 4)\nint: \"3\"", str)

	assert.Equal(t, "", LexemeList(nil))
}

func TestLexemeString(t *testing.T) {
	tt := map[string]struct {
		lxs      []Lexeme
		expected string
	}{
		"empty": {},
		"one": {
			lxs: []Lexeme{
				&lx{k: "int", v: "1"},
			},
			expected: "int: 1",
		},
		"no-val": {
			lxs: []Lexeme{
				&lx{k: "int"},
			},
			expected: "int",
		},
		"many": {
			lxs: []Lexeme{
				&lx{k: "int", v: "1"},
				&lx{k: "op", v: "+"},
				&lx{k: "int", v: "2"},
				&lx{k: "op", v: "*"},
				&lx{k: "int", v: "3"},
			},
			expected: "[int: 1, op: +, int: 2, op: *, int: 3]",
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, LexemeString(tc.lxs...))
		})
	}
}

func TestMustLexer(t *testing.T) {
	l := &testLexer{}
	assert.Equal(t, l, MustLexer(l, nil))

	defer func() {
		assert.Equal(t, testErr, recover())
	}()
	MustLexer(l, testErr)
}
