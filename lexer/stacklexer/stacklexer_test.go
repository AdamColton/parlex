package stacklexer

import (
	"github.com/adamcolton/parlex/lexeme"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStacklexer(t *testing.T) {
	lxr, err := New(`
    == main ==
      START innerLexer
      outerword  /\w+/
      shared
    == innerLexer ==
      STOP ^
      innerword  /\w+/
      shared
    == shared ==
      space /\s+/ -
      nl /\n/ -
  `)
	assert.NoError(t, err)
	lxms := lxr.Lex("this \n START is \n a STOP test")
	if !assert.Len(t, lxms, 6) {
		return
	}

	excpected := []*lexeme.Lexeme{
		lexeme.String("outerword").Set("this").At(1, 1),
		lexeme.String("START").Set("START").At(2, 2),
		lexeme.String("innerword").Set("is").At(2, 8),
		lexeme.String("innerword").Set("a").At(3, 2),
		lexeme.String("STOP").Set("STOP").At(3, 4),
		lexeme.String("outerword").Set("test").At(3, 9),
	}

	for i, e := range excpected {
		lx := lxms[i]
		assert.Equal(t, e.K.String(), lx.Kind().String())
		assert.Equal(t, e.V, lx.Value())
		ex, ey := e.Pos()
		gx, gy := lx.Pos()
		if !assert.Equal(t, ex, gx, "Line") || !assert.Equal(t, ey, gy, "Col") {
			t.Error(i)
		}
	}
}

func TestStacklexerByPriority(t *testing.T) {
	lxr, err := New(`
    == main ==
      START innerLexer
      outerword  /\w+/
      shared
    == innerLexer ==
      STOP ^
      innerword  /\w+/
      shared
    == shared ==
      space /\s+/ -
      nl /\n/ -
  `)
	assert.NoError(t, err)
	lxr.ByPriority()
	lxms := lxr.Lex("this \n START is \n a STOP test")
	if !assert.Len(t, lxms, 6) {
		return
	}

	excpected := []*lexeme.Lexeme{
		lexeme.String("outerword").Set("this").At(1, 1),
		lexeme.String("START").Set("START").At(2, 2),
		lexeme.String("innerword").Set("is").At(2, 8),
		lexeme.String("innerword").Set("a").At(3, 2),
		lexeme.String("STOP").Set("STOP").At(3, 4),
		lexeme.String("outerword").Set("test").At(3, 9),
	}

	for i, e := range excpected {
		lx := lxms[i]
		assert.Equal(t, e.K.String(), lx.Kind().String())
		assert.Equal(t, e.V, lx.Value())
		ex, ey := e.Pos()
		gx, gy := lx.Pos()
		if !assert.Equal(t, ex, gx, "Line") || !assert.Equal(t, ey, gy, "Col") {
			t.Error(i)
		}
	}
}

func TestSubParserLine(t *testing.T) {
	tests := []struct {
		in       string
		expected []string
	}{
		{"test", []string{"test", "test", "", "", "", ""}},
		{"name /regex/", []string{"name /regex/", "name", "regex", "", "", ""}},
		{"name /regex/ ^", []string{"name /regex/ ^", "name", "regex", "", "^", ""}},
		{"name /regex/ sublexer", []string{"name /regex/ sublexer", "name", "regex", "", "sublexer", ""}},
		{"name /regex/ sublexer -", []string{"name /regex/ sublexer -", "name", "regex", "", "sublexer", "-"}},
		{"name /regex/ (matches) sublexer -", []string{"name /regex/ (matches) sublexer -", "name", "regex", "matches", "sublexer", "-"}},
	}

	for _, test := range tests {
		got := subParserLine.FindStringSubmatch(test.in)
		assert.Equal(t, test.expected, got)
	}
}

func TestStacklexerSubSection(t *testing.T) {
	lxr, err := New(`
    == main ==
      START innerLexer
      outerword  /\w+/
      shared
    == innerLexer ==
      STOP ^
      foo /foo\n(\w+)foo/ (1)
      innerword  /\w+/
      shared
    == shared ==
      space /\s+/ -
      nl /\n/ -
  `)
	assert.NoError(t, err)
	lxr.ByPriority()
	lxms := lxr.Lex("this \n START foo\nbarfoo is \n a STOP test")
	if !assert.Len(t, lxms, 7) {
		return
	}

	excpected := []*lexeme.Lexeme{
		lexeme.String("outerword").Set("this").At(1, 1),
		lexeme.String("START").Set("START").At(2, 2),
		lexeme.String("foo").Set("bar").At(2, 8),
		lexeme.String("innerword").Set("is").At(3, 8),
		lexeme.String("innerword").Set("a").At(4, 2),
		lexeme.String("STOP").Set("STOP").At(4, 4),
		lexeme.String("outerword").Set("test").At(4, 9),
	}

	for i, e := range excpected {
		lx := lxms[i]
		assert.Equal(t, e.K.String(), lx.Kind().String())
		assert.Equal(t, e.V, lx.Value())
		ex, ey := e.Pos()
		gx, gy := lx.Pos()
		if !assert.Equal(t, ex, gx, "Line: "+e.V) || !assert.Equal(t, ey, gy, "Col: "+e.V) {
			t.Error(i)
		}
	}
}
