package main

import (
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer"
	"github.com/adamcolton/parlex/parser"
	"github.com/adamcolton/parlex/tree"
	"strconv"
)

const lexerRules = `
  space  /\s+/ -
  number /\d*\.?\d+/
  op1 /[\*\/]/
  op2 /[\+\-]/
  ( /\(/
  ) /\)/
`
const grammarRules = `
  E -> E op2 E
    -> E op1 E
    -> number
    -> P
  P -> ( E )
`

var reducer = tree.Reducer{
	"E": func(node *tree.PN) {
		if ln := len(node.C); ln == 1 {
			node.PromoteSingleChild()
		} else if ln == 3 {
			if k := node.C[1].Kind(); k == "op1" || k == "op2" {
				node.PromoteChild(1)
			}
		}
	},
	"P": func(node *tree.PN) {
		*node = *(node.C[1])
	},
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	lxr, err := lexer.New(lexerRules)
	check(err)

	grmr, err := grammar.New(grammarRules)
	check(err)

	p := parser.Packrat(grmr)

	ok := true
	expected := map[string]int{
		"1+2+3":         6,
		"2*2+3":         7,
		"1+2*3":         7,
		"1+2-3":         0,
		"1*(2+3)":       5,
		"2*(2+3*2)-2*3": 10,
	}
	for str, i := range expected {
		tr := parlex.Run(str, lxr, p, reducer)
		ei := eval(tr)
		fmt.Println(str, "\n", tr, ei, "=", i)
		ok = ok && ei == i
	}
	fmt.Println(ok)
}

func eval(node parlex.ParseNode) int {
	switch node.Kind() {
	case "number":
		i, err := strconv.Atoi(node.Value())
		check(err)
		return i
	case "op1", "op2":
		a := eval(node.Child(0))
		b := eval(node.Child(1))
		switch node.Value() {
		case "+":
			return a + b
		case "*":
			return a * b
		case "/":
			return a / b
		case "-":
			return a - b
		}
	}
	return 0
}
