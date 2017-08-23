package regexgram

import (
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/parser/packrat"
	"github.com/adamcolton/parlex/tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		lex   string
		kinds []string
		vals  []string
	}{
		{
			lex:   "rule -> prod1 prod2",
			kinds: []string{"nl", "symbol", "rarr", "symbol", "symbol"},
			vals:  []string{"\n", "rule", "->", "prod1", "prod2"},
		},
	}

	for _, tc := range testCases {
		for i, lx := range lxr.Lex(tc.lex) {
			assert.Equal(t, tc.kinds[i], lx.Kind().String())
			assert.Equal(t, tc.vals[i], lx.Value())
		}
	}
}

func TestReduce(t *testing.T) {
	lxms := lxr.Lex(`
    rule -> prod1? prod2* (A B)*
    B -> C D E
    -> X Y
    -> W|X*|Y|Z|(A B)
  `)
	pn := prsr.Parse(lxms)
	pn = rdcr.Reduce(pn)
	expected, err := tree.New(`
      Grammar {
      Production: "rule" {
        OptSymbol {
          symbol: "prod1"
        }
        RepSymbol {
          symbol: "prod2"
        }
        RepSymbol {
          Group {
            symbol: "A"
            symbol: "B"
          }
        }
      }
      Production: "B" {
        symbol: "C"
        symbol: "D"
        symbol: "E"
      }
      ContinueProd {
        symbol: "X"
        symbol: "Y"
      }
      ContinueProd {
        OrSymbol {
          symbol: "W"
          RepSymbol {
            symbol: "X"
          }
          symbol: "Y"
          symbol: "Z"
          Group {
            symbol: "A"
            symbol: "B"
          }
        }
      }
    }
  `)
	assert.NoError(t, err)
	if expected.String() != pn.(*tree.PN).String() {
		t.Error("\n" + pn.(*tree.PN).String() + "====\n" + expected.String())
	}
}

func TestEval(t *testing.T) {
	grammarString := `
    // This is a comment
    rule   -> prod1 prod2
    rule2  -> foo? bar?
    rule3  -> A|B C
    rule4  -> A|B C?
           -> X
    rule5  -> (J K L)? M
    MoreOr -> (Group|OptSymbol|RepSymbol|symbol) (or MoreOr)?
    rule6  -> (A B)* C
    rule7  -> value (comma value)*
  `
	grmr, rdcr, err := New(grammarString)
	assert.NotNil(t, rdcr)
	assert.NoError(t, err)
	expected, err := grammar.New(`
    rule           -> prod1 prod2
    rule2          -> foo bar
                   -> foo
                   -> bar
                   -> 
    rule3          -> A C
                   -> B C
    rule4          -> A C
                   -> A
                   -> B C
                   -> B
                   -> X
    rule5          -> J K L M
                   -> M
    MoreOr         -> Group or MoreOr
                   -> Group
                   -> OptSymbol or MoreOr
                   -> OptSymbol
                   -> RepSymbol or MoreOr
                   -> RepSymbol
                   -> symbol or MoreOr
                   -> symbol
    rule6          -> (A_B)* C
    rule7          -> value (comma_value)*
    (A_B)*         -> A B (A_B)*
                   -> 
    (comma_value)* -> comma value (comma_value)*
                   ->
  `)
	assert.NoError(t, err)
	if expected.String() != grmr.String() {
		t.Error("\n" + grmr.String() + "====\n" + expected.String())
	}
}

func TestMergeRules(t *testing.T) {
	a := rules{
		rule{"A", "B"},
		rule{"C"},
		rule{"D", "E"},
	}
	b := rules{
		rule{},
		rule{"X"},
		rule{"Y", "Z"},
	}

	expected := rules{
		rule{"A", "B"},
		rule{"A", "B", "X"},
		rule{"A", "B", "Y", "Z"},
		rule{"C"},
		rule{"C", "X"},
		rule{"C", "Y", "Z"},
		rule{"D", "E"},
		rule{"D", "E", "X"},
		rule{"D", "E", "Y", "Z"},
	}
	assert.Equal(t, expected, mergeRules(a, b))
}

func TestOuputReducer(t *testing.T) {
	lxr, err := simplelexer.New(`
    int   /\d+/
    word  /\w+/
    comma /,/
  `)
	assert.NoError(t, err)
	grmr, rdcr := Must(`
    List -> Value MoreValues*
    Value -> int|word
    MoreValues -> comma Value
  `)
	rdcr = tree.Merge(rdcr, tree.Reducer{
		"MoreValues": tree.ReplaceWithChild(1),
	})
	lxms := lxr.Lex("1,test,3,4")
	pn := packrat.New(grmr).Parse(lxms)
	pn = rdcr.Reduce(pn)

	expected, err := tree.New(`
    List {
      Value {
        int: "1"
      }
      Value {
        word: "test"
      }
      Value {
        int: "3"
      }
      Value {
        int: "4"
      }
    }
  `)

	assert.Equal(t, expected.String(), pn.(*tree.PN).String())
}

func TestDontRepeatRepeat(t *testing.T) {
	grammarString := `
    rule1 -> lexeme1* lexeme2
    rule2 -> lexeme1* lexeme3
  `
	grmr, rdcr, err := New(grammarString)
	assert.NotNil(t, rdcr)
	assert.NoError(t, err)

	expectedStr := `
    rule1    -> lexeme1* lexeme2
    rule2    -> lexeme1* lexeme3
    lexeme1* -> lexeme1 lexeme1*
             ->
  `
	expectGrmr, err := grammar.New(expectedStr)
	assert.NoError(t, err)
	assert.Equal(t, expectGrmr.String(), grmr.String())
}
