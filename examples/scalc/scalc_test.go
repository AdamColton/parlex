package scalc

import (
	"github.com/adamcolton/parlex/tree"
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

func eval(node *tree.PN) []string {
	stack := evalStack(node)
	out := make([]string, len(stack))
	for i, s := range stack {
		out[i] = s.String()
	}
	return out
}

func TestCases(t *testing.T) {
	//t.Skip()
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
		makeCase("-2 abs", "2"),
		makeCase("2 -- abs", "2"),
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
		makeCase("1 2 + swap", "3"), // ??? is this correct
		makeCase("1 2 drop", "1"),
		makeCase("1 drop"),
		makeCase("1 2 3 clear"),
		makeCase("1 2.000 *", "2.000"),
		makeCase("1 3.0 /", "0.3"),
		makeCase("1 2 / 1 2 / +", "1"),
		makeCase("1 2 >", "0"),
		makeCase("2 1 >", "1"),
		makeCase("1 2 <", "1"),
		makeCase("2 1 <", "0"),
		makeCase("2 1 =", "0"),
		makeCase("2 2 =", "1"),
		makeCase("2 2 cmpr", "0"),
		makeCase("1 2 cmpr", "-1"),
		makeCase("2 1 cmpr", "1"),
		makeCase("2 1 0 ?", "2"),
		makeCase("2 1 3 ?", "1"),
		makeCase("2 3 + * 3 ?", "6"),
		makeCase("2 3 + * 0 ?", "5"),
		makeCase("3 abs -- 1 ?", "-3"),
		makeCase("2 3 - + 1 ? * 0 ?", "5"),
		makeCase("2 3 - + 0 ? * 0 ?", "-1"),
		makeCase("2 3 - + 1 ? * 1 ?", "6"),
		makeCase("2 3 - + 0 ? * 1 ?", "6"),
		makeCase("2 3 swap drop 1 ?", "2"),
		makeCase("2 3 swap drop 0 ?", "3", "2"),
		makeCase("-5--+6+*1+2++3=?", "30"),
	}

	var tt testCase
	for _, tt = range tests {
		pn := Parse(tt.expr)
		if pn == nil {
			t.Error("Could not parse", tt.expr)
			continue
		}
		tpn := pn.(*tree.PN)
		str := tpn.String()
		if !assert.Equal(t, tt.expect, eval(tpn), tt.expr) {
			t.Error(str)
			t.Error(pn)
		}
	}
}

func TestParseFailsAsNil(t *testing.T) {
	pn := Parse("you can't parse me!")
	assert.True(t, pn == nil)
}

func TestPad(t *testing.T) {
	t.Skip()
	pn := prsr.Parse(lxr.Lex("2 3 swap drop 1 ?"))
	//t.Error(pn)
	pn = rdcr.Reduce(pn)
	t.Error(pn)
	t.Error(eval(pn.(*tree.PN)))
}
