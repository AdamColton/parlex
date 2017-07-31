package parlexjson

import (
	"bytes"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar/regexgram"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/parser/topdown"
	"github.com/adamcolton/parlex/tree"
)

const lexerRules = `
  space  /\s+/ -
  number /\d*\.?\d+/
  string /(\"([^\"\\]|(\\.))*\")/
  bool   /(true)|(false)/
  null
  comma  /,/
  lcb    /\{/
  rcb    /\}/
  lb     /\[/
  rb     /\]/
  colon  /:/
`

var lxr = parlex.MustLexer(simplelexer.New(lexerRules))

const grammarRules = `
  Value       -> string | number | bool | null | Array | Object
  Array       -> lb ( Value MoreVals* )? rb
  MoreVals    -> comma Value
  Object      -> lcb ( KeyVal MoreKeyVals* )? rcb
  MoreKeyVals -> comma KeyVal
  KeyVal      -> string colon Value
`

var grmr, grmrRdcr = regexgram.Must(grammarRules)
var prsr = parlex.MustParser(topdown.New(grmr))

var rdcr = tree.Merge(grmrRdcr, tree.Reducer{
	"Value":       tree.PromoteSingleChild,
	"Object":      tree.RemoveChildren(0, -1),               // remove { }
	"Array":       tree.RemoveChildren(0, -1),               // remove [ ]
	"KeyVal":      tree.PromoteChildValue(0).RemoveChild(0), // Promote key, remove :
	"MoreVals":    tree.ReplaceWithChild(1),
	"MoreKeyVals": tree.ReplaceWithChild(1),
})

var runner = parlex.New(lxr, prsr, rdcr)

// Format takes a json string and formats it, returning it as a string. If there
// is an error, that will be returned with an empty string.
func Format(s string) (string, error) {
	var buf bytes.Buffer
	tr, err := runner.Run(s)
	if err != nil {
		return "", err
	}
	prettyPrint(tr, &buf, "")
	return buf.String(), nil
}

func prettyPrint(node parlex.ParseNode, buf *bytes.Buffer, pad string) {
	switch node.Kind().String() {
	case "Array":
		buf.WriteString("[")
		cpad := pad + "  "
		prepend := false
		for i := 0; i < node.Children(); i++ {
			if i != 0 {
				buf.WriteString(",")
			}
			child := node.Child(i)
			if cs := child.Kind().String(); cs == "Array" || cs == "Object" {
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
		buf.WriteString(node.Value())
		buf.WriteString(": ")
		prettyPrint(node.Child(0), buf, pad)
	default:
		buf.WriteString(node.Value())
	}
}
