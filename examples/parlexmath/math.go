package parlexmath

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/parser/packrat"
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
    -> op2 number
    -> P
  P -> ( E )
`

var lxr = parlex.MustLexer(simplelexer.New(lexerRules))
var grmr = parlex.MustGrammar(grammar.New(grammarRules))
var prsr = packrat.New(grmr)

var reducer = tree.Reducer{
	"E": func(node *tree.PN) {
		switch node.Children() {
		case 1, 2:
			node.PromoteChild(0)
		case 3:
			node.PromoteChild(1)
		}
	},
	"P": tree.ReplaceWithChild(1),
}

var runner = parlex.New(lxr, prsr, reducer)

// Eval takes an expression and tries to evaluate it. If the evaluation is
// successful, the value is returned. If not, an error is returned.
func Eval(expr string) (float64, error) {
	tr, err := runner.Run(expr)
	if err != nil {
		return 0, err
	}
	return eval(tr), nil
}

func eval(node parlex.ParseNode) float64 {
	switch node.Kind().String() {
	case "number":
		i, _ := strconv.ParseFloat(node.Value(), 64)
		return i
	case "op2":
		if node.Children() == 1 {
			a := eval(node.Child(0))
			if node.Value() == "-" {
				a = -a
			}
			return a
		}
		fallthrough
	case "op1":
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
