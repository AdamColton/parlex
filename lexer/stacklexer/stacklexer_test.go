package stacklexer

import (
	"github.com/adamcolton/parlex"
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
      foo /This will mask shared foo/
      shared
    == shared ==
      foo
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

func TestLexerError(t *testing.T) {
	lxr, err := New(`
    == main ==
      foo
      bar barLexer
      shared
    == barLexer ==
      glorp
      shared
    == shared ==
      space /\s+/ -
      nl /\n/ -
  `)
	assert.NoError(t, err)

	lxms := lxr.Lex("foo bar error1 foo error2")
	assert.Len(t, lxms, 5)

	errs := parlex.LexErrors(lxms)
	assert.Len(t, errs, 2)
}

func TestDoublePop(t *testing.T) {
	lxr, err := New(`
    == Main ==
      level1 Level1
      bar
      Shared
    == Level1 ==
      level2 Level2
      badBar /bar/
      Shared
    == Level2 ==
      exit ^^
      badBar /bar/
      Shared
    == Shared ==
      space /\s+/ -
      nl /\n/ -
  `)
	assert.NoError(t, err)

	lxms := lxr.Lex("bar level1 level2 exit bar")
	if assert.Len(t, lxms, 5) {
		assert.Equal(t, "level1", lxms[1].Kind().String())
		assert.Equal(t, "level2", lxms[2].Kind().String())
		assert.Equal(t, "exit", lxms[3].Kind().String())
		assert.Equal(t, "bar", lxms[4].Kind().String())
	}
}
