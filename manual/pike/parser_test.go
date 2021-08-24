package pike

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserTable(t *testing.T) {
	type input struct {
		val    string
		ln     int
		groups map[uint32][][2]int
	}
	tt := map[string]struct {
		expected string
		inputs   []input
	}{
		"c(a*|o)t": {
			inputs: []input{
				{
					val:    "cat",
					groups: map[uint32][][2]int{0: {{1, 2}}},
				},
				{
					val:    "cot",
					groups: map[uint32][][2]int{0: {{1, 2}}},
				},
				{
					val:    "caat",
					groups: map[uint32][][2]int{0: {{1, 3}}},
				},
				{
					val:    "caaat",
					groups: map[uint32][][2]int{0: {{1, 4}}},
				},
			},
		},
		"a|b{2,2}": { // the {2,2} should be associated with "b" not "a|b"
			inputs: []input{
				{val: "a"},
				{val: "bb"},
				{val: "ba", ln: -1},
				{val: "b", ln: -1},
			},
		},
		"ca{2,3}t": {
			inputs: []input{
				{val: "cat", ln: -1},
				{val: "caat"},
				{val: "caaat"},
				{val: "caaaat", ln: -1},
			},
		},
		"ca{,3}t": {
			inputs: []input{
				{val: "cat"},
				{val: "caat"},
				{val: "caaat"},
				{val: "caaaat", ln: -1},
			},
		},
		"ca{3,}t": {
			inputs: []input{
				{val: "cat", ln: -1},
				{val: "caat", ln: -1},
				{val: "caaat"},
				{val: "caaaat"},
			},
		},
		"ca*t": {
			inputs: []input{
				{val: "ct"},
				{val: "cat"},
				{val: "caaat"},
				{val: "caaaat"},
				{val: "aaat", ln: -1},
			},
		},
		"ca|ot": {
			inputs: []input{
				{val: "cat"},
				{val: "cot"},
				{val: "caot", ln: -1},
			},
		},
		"cat": {
			inputs: []input{
				{val: "cat"},
				{val: "ct", ln: -1},
				{val: "cot", ln: -1},
				{val: "caot", ln: -1},
			},
		},
		"c.t": {
			inputs: []input{
				{val: "cat"},
				{val: "ct", ln: -1},
				{val: "cot"},
				{val: "caot", ln: -1},
			},
		},
		"ca+t": {
			inputs: []input{
				{val: "cat"},
				{val: "ct", ln: -1},
				{val: "caat"},
				{val: "caot", ln: -1},
			},
		},
		"ca?t": {
			inputs: []input{
				{val: "cat"},
				{val: "ct"},
				{val: "caat", ln: -1},
				{val: "caot", ln: -1},
			},
		},
		"c(a*|o{2,3}){1,2}t": {
			inputs: []input{
				{val: "ct"},
				{val: "caat"},
				{val: "coot"},
				{val: "cot", ln: -1},
			},
		},
		"a?a?a{2,2}": {
			inputs: []input{
				{val: "aa"},
			},
		},
	}

	out, _ := os.Create("out.txt")

	for re, tc := range tt {
		t.Run(re, func(t *testing.T) {
			fmt.Fprintln(out, re)
			n := parse(re)
			fmt.Fprintln(out, n.Tree(""))
			exp := tc.expected
			if exp == "" {
				exp = re
			}
			assert.Equal(t, exp, n.String())
			p := buildTree(n)
			p.optimize()
			fmt.Fprintln(out, printCode(p.code))

			for _, i := range tc.inputs {
				t.Run(i.val, func(t *testing.T) {
					op := p.run(i.val)
					ln := i.ln
					if ln == 0 {
						ln = len(i.val)
					}
					assert.Equal(t, ln, op.best)
					if ln > -1 {
						g := i.groups
						if g != nil {
							assert.Equal(t, g, op.groupMap.toMap(op.bestGroups))
						}
					} else {
						assert.Equal(t, groupID(0), op.bestGroups)
					}
				})
			}
		})
	}
}

func BenchmarkPathological(b *testing.B) {
	gen := func(i int) (string, string) {
		re := strings.Repeat("a?", i) + fmt.Sprintf("a{%d,%d}", i, i)
		str := strings.Repeat("a", i)
		return re, str
	}

	i := 700
	re, str := gen(i)
	p := compile(re)
	b.ResetTimer()
	op := p.run(str)
	assert.Equal(b, i, op.best)
}
