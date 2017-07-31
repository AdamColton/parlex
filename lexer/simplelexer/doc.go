// Package simplelexer provides a very basic lexer that fulfills parlex.Lexer.
// The easiest way to define an lexer is using a string. Each line in the string
// will define one rule. A line may start with white space, then it must have a
// word, which will be the symbol. If the symbol is the same as the value, that
// is all that's required. If the symbol is to match a regular expression, the
// regexp should be delimited by /'s. A rule can optionally end with "-" to
// indicate that the value should be dropped, which is often helpful to
// eliminate whitespace.
//
// An example of the simple lexer can be seen in
// parlex/examples/parlexmath
package simplelexer
