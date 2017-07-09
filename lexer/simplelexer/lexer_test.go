package simplelexer

import (
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
	s := "this is \na test"
	lxr, err := New(`
    test
    word  /\w+/
    space /\s+/
  `)
	lxr.
		InsertStart("START_K", "START_V").
		InsertEnd("END_K", "END_V")
	assert.NoError(t, err)
	lxs := lxr.Lex(s)

	if assert.Equal(t, 9, len(lxs)) {
		assert.Equal(t, "START_K", lxs[0].Kind().String())
		assert.Equal(t, "START_V", lxs[0].Value())
		assert.Equal(t, "word", lxs[1].Kind().String())
		assert.Equal(t, "this", lxs[1].Value())
		assert.Equal(t, "space", lxs[2].Kind().String())
		assert.Equal(t, " ", lxs[2].Value())
		// confirm that test takes priority over word
		assert.Equal(t, "test", lxs[7].Kind().String())
		assert.Equal(t, "test", lxs[7].Value())
		assert.Equal(t, "END_K", lxs[8].Kind().String())
		assert.Equal(t, "END_V", lxs[8].Value())
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
		assert.Equal(t, "word", lxs[1].Kind().String())
		assert.Equal(t, "is", lxs[1].Value())
		// confirm that test takes priority over word
		assert.Equal(t, "test", lxs[3].Kind().String())
		assert.Equal(t, "test", lxs[3].Value())
	}
}
