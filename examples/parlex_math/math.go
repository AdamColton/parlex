package main

import (
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/parser/packrat"
	"github.com/adamcolton/parlex/tree"
	"os"
	"strconv"
	"strings"
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
    -> op2 number
    -> P
  P -> ( E )
`

var lxr = parlex.MustLexer(simplelexer.New(lexerRules))
var grmr = parlex.MustGrammar(grammar.New(grammarRules))
var prsr = packrat.New(grmr)

var reducer = tree.Reducer{
	"E": func(node *tree.PN) {
		if !node.PromoteSingleChild() {
			node.PromoteChild(1)
		}
	},
	"P": tree.ReplaceWithChild(1),
}

var runner = parlex.New(lxr, prsr, reducer)

func main() {
	tr, err := runner.Run(strings.Join(os.Args[1:], " "))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Println(eval(tr))
}

func eval(node parlex.ParseNode) float64 {
	switch node.Kind().String() {
	case "number":
		i, _ := strconv.ParseFloat(node.Value(), 64)
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
