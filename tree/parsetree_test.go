package tree

import (
	"github.com/adamcolton/parlex"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseNode(t *testing.T) {
	treeNode := &PN{
		Lexeme: &parlex.L{
			K: "foo",
			V: "bar",
		},
	}

	parlexNode := parlex.ParseNode(treeNode)
	assert.NotNil(t, parlexNode)
}

func TestParseNodeString(t *testing.T) {
	pn1, err := New(`
		E {
			E {
				int: '1'
			}
			op: '+'
			E {
				int: '2'
			}
	  }
  `)
	assert.NoError(t, err)
	assert.Equal(t, pn1.sliceReq(), len(pn1.string("", nil)))

	pn2, err := New(pn1.String())
	assert.Equal(t, pn1.String(), pn2.String())
}

func TestReTreeLine(t *testing.T) {
	tests := []struct {
		str                         string
		valid                       bool
		kind, value                 string
		hasleftBrace, hasRightBrace bool
	}{
		{"} op : '+' {", true, "op", "+", true, true},
		{" } ", true, "", "", true, false},
		{"{", false, "", "", false, false},
	}

	for _, test := range tests {
		m := reTreeLine.FindStringSubmatch(test.str)
		t.Log("> ", strings.Join(m, "|"), "\n", test.str)
		if test.valid {
			if assert.Equal(t, len(test.str), len(m[0])) {
				assert.Equal(t, test.kind, m[2], test.str)
				assert.Equal(t, test.value, m[3], test.str)
				assert.True(t, test.hasleftBrace == (m[1] == "}"), test.str)
				assert.True(t, test.hasRightBrace == (m[4] == "{"), test.str)
			}
		} else {
			assert.Equal(t, 0, len(m[0]))
		}
	}

}
