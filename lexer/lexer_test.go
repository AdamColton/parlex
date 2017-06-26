package lexer

import (
	"github.com/adamcolton/parlex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexStr(t *testing.T) {
	tests := []struct {
		str, kind, regex string
		hasMinus         bool
	}{
		{"test", "test", "", false},
		{"foo /bar/ -", "foo", "bar", true},
		{"foo/bar/ -", "foo", "bar", true},             // space not actually required
		{"   foo  \t  /bar/ -   ", "foo", "bar", true}, // leading and trailing spaces are ignored
	}

	for _, test := range tests {
		m := lexStr.FindStringSubmatch(test.str)
		assert.Equal(t, 4, len(m))
		assert.Equal(t, test.kind, m[1])
		assert.Equal(t, test.regex, m[2])
		assert.True(t, test.hasMinus == (m[3] == "-"))
	}

	m := lexStr.FindStringSubmatch("  /  ")
	assert.Equal(t, 0, len(m))
}

func TestLexerString(t *testing.T) {
	str := `
    test
    word  /\w+/
    space /\s+/ -
  `
	lxr1, err := New(str)
	assert.NoError(t, err)
	assert.NotNil(t, lxr1)
	lxr2, err := New(lxr1.String())
	assert.NoError(t, err)
	assert.Equal(t, lxr1.String(), lxr2.String())
}

func TestLex(t *testing.T) {
	s := "this is a test"
	lxr, err := New(`
    test
    word  /\w+/
    space /\s+/
  `)
	assert.NoError(t, err)
	lxs := lxr.Lex(s)

	if assert.Equal(t, 7, len(lxs)) {
		assert.Equal(t, parlex.Symbol("word"), lxs[0].Kind())
		assert.Equal(t, "this", lxs[0].Value())
		assert.Equal(t, parlex.Symbol("space"), lxs[1].Kind())
		assert.Equal(t, " ", lxs[1].Value())
		// confirm that test takes priority over word
		assert.Equal(t, parlex.Symbol("test"), lxs[6].Kind())
		assert.Equal(t, "test", lxs[6].Value())
	}
}

func TestLexDiscard(t *testing.T) {
	s := "this is a test"
	lxr, err := New(`
    test
    word /\w+/
    space /\s+/ -
  `)
	assert.NoError(t, err)
	lxs := lxr.Lex(s)

	if assert.Equal(t, 4, len(lxs)) {
		assert.Equal(t, parlex.Symbol("word"), lxs[1].Kind())
		assert.Equal(t, "is", lxs[1].Value())
		// confirm that test takes priority over word
		assert.Equal(t, parlex.Symbol("test"), lxs[3].Kind())
		assert.Equal(t, "test", lxs[3].Value())
	}
}
