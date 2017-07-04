package scalc

import (
	//"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/tree"
	//"github.com/adamcolton/parlex/symbol/stringsymbol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	lxm := lxr.Lex("1 2 +")
	if !assert.Len(t, lxm, 3) {
		return
	}

	assert.Equal(t, "int", lxm[0].Kind().String())
	assert.Equal(t, "int", lxm[1].Kind().String())
	assert.Equal(t, "bop", lxm[2].Kind().String())
}

type testCase struct {
	expr   string
	expect []string
}

func makeCase(expr string, expect ...string) testCase {
	if expect == nil {
		expect = make([]string, 0)
	}
	return testCase{
		expr:   expr,
		expect: expect,
	}
}

func TestCases(t *testing.T) {
	tests := []testCase{
		makeCase("2", "2"),
		makeCase("2.1", "2.1"),
		makeCase("-2", "-2"),
		makeCase("-2.1", "-2.1"),
		makeCase("1 2", "1", "2"),
		makeCase("1 2 3", "1", "2", "3"),
		makeCase("(2)", "2"),
		makeCase("1 (2)", "1", "2"),
		makeCase("(1) (2)", "1", "2"),
		makeCase("(1) 3 (2)", "1", "3", "2"),
		makeCase("2 --", "-2"),
		makeCase("2 3 +", "5"),
		makeCase("2 3 -", "-1"),
		makeCase("2 3 *", "6"),
		makeCase("2 3 /", "1"),
		makeCase("2 3 ^", "8"),
		makeCase("14 4 %", "2"),
		makeCase("2 2 3 sum", "7"),
		makeCase("2 2 3 sum *", "10"),
		makeCase("2 5 3 max", "5"),
		makeCase("3 (2 5 3 max) avg", "4"),
		makeCase("1 1 1 len", "3"),
		makeCase("1 2 3 min", "1"),
		makeCase("1 2 3 first", "3"),
		makeCase("1 2 3 last", "1"),
		makeCase("1 2 swap", "2", "1"),
		makeCase("1 2 drop", "1"),
		makeCase("1 drop"),
		makeCase("1 2 3 clear"),
		makeCase("1 2.000 *", "2.000"),
		makeCase("1 3.0 /", "0.3"),
		makeCase("1 2 / 1 2 / +", "1"),
	}

	for _, tt := range tests {
		pn := prsr.Parse(lxr.Lex(tt.expr))
		pn = rdcr.Reduce(pn)
		if !assert.Equal(t, tt.expect, eval(pn.(*tree.PN)), tt.expr) {
			t.Error(pn)
		}
	}
}

func TestPad(t *testing.T) {
	t.Skip()
	pn := prsr.Parse(lxr.Lex("1 drop"))
	t.Error(pn)
	pn = rdcr.Reduce(pn)
	t.Error(pn)
	t.Error(eval(pn.(*tree.PN)))
}
