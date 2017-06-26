package main

import (
	"bytes"
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer"
	"github.com/adamcolton/parlex/parser"
	"github.com/adamcolton/parlex/tree"
	"os"
	"strings"
)

const lexerRules = `
  space  /\s+/ -
  number /\d*\.?\d+/
  string /(\"([^\"\\]|(\\.))*\")/
  bool   /(true)|(false)/
  null
  ,
  {
  }
  [      /\[/
  ]
  :
`

var lxr = parlex.MustLexer(lexer.New(lexerRules))

const grammarRules = `
  Value         -> string
                -> number
                -> bool
                -> null
                -> Array
                -> Object
  Object        -> { ObjectDef }
  ObjectDef     -> KeyVal , ObjectDef
                -> KeyVal
  KeyVal        -> Key : Value
  Key           -> string
  Array         -> [ ArrayContents ]
  ArrayContents -> Value , ArrayContents
                -> Value
`

var grmr = parlex.MustGrammar(grammar.New(grammarRules))
var prsr = parlex.MustParser(parser.TopDown(grmr))

func reduceList(node *tree.PN) {
	node.RemoveChild(1)       // remove comma
	node.PromoteChildrenOf(1) // promote recursive portion of list
}

func reduceListWrapper(node *tree.PN) {
	node.RemoveChildren(0, -1)  // remove wrapper
	node.PromoteGrandChildren() // promote rows in ObjectDef
}

var reducer = tree.Reducer{
	"Value":         tree.PromoteSingleChild,
	"ObjectDef":     reduceList,
	"Object":        reduceListWrapper,
	"KeyVal":        tree.RemoveChild(1), // remove :
	"Key":           tree.PromoteChildValue(0),
	"ArrayContents": reduceList,
	"Array":         reduceListWrapper,
}

var runner = parlex.New(lxr, prsr, reducer)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	s := strings.Join(os.Args[1:], " ")
	var buf bytes.Buffer
	tr := runner.Run(s)
	prettyPrint(tr, &buf, "")
	fmt.Println(buf.String())
}

func prettyPrint(node parlex.ParseNode, buf *bytes.Buffer, pad string) {
	switch node.Kind() {
	case "string", "number", "bool", "null", "Key":
		buf.WriteString(node.Value())
	case "Array":
		buf.WriteString("[")
		cpad := pad + "  "
		prepend := false
		for i := 0; i < node.Children(); i++ {
			if i != 0 {
				buf.WriteString(",")
			}
			child := node.Child(i)
			if child.Kind() == "Array" || child.Kind() == "Object" {
				buf.WriteString("\n")
				buf.WriteString(cpad)
				prepend = true
			}
			prettyPrint(child, buf, cpad)
		}
		if prepend {
			buf.WriteString("\n")
			buf.WriteString(pad)
		}
		buf.WriteString("]")
	case "Object":
		if node.Children() == 0 {
			buf.WriteString("{}")
			return
		}
		if node.Children() == 1 {
			buf.WriteString("{")
			prettyPrint(node.Child(0), buf, pad)
			buf.WriteString("}")
			return
		}
		buf.WriteString("{\n")
		cpad := pad + "  "
		for i := 0; i < node.Children(); i++ {
			if i != 0 {
				buf.WriteString(",\n")
			}
			buf.WriteString(cpad)
			prettyPrint(node.Child(i), buf, cpad)
		}
		buf.WriteString("\n")
		buf.WriteString(pad)
		buf.WriteString("}")
	case "KeyVal":
		prettyPrint(node.Child(0), buf, pad)
		buf.WriteString(": ")
		prettyPrint(node.Child(1), buf, pad)
	}

}
