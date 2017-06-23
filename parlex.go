package parlex

// Run performs the lexing, parsing and reducing for an input
func Run(input string, lexer Lexer, parser Parser, reducer Reducer) ParseNode {
	lexemes := lexer.Lex(input)
	tree := parser.Parse(lexemes)
	if reducer != nil {
		tree = reducer.Reduce(tree)
	}
	return tree
}
