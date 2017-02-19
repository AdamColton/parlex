## Parlex
[![GoDoc](https://godoc.org/github.com/AdamColton/parlex?status.svg)](https://godoc.org/github.com/AdamColton/parlex)

For a quick demo check out [json.go](https://github.com/AdamColton/parlex/blob/master/examples/json.go)

### Todo
Double check that parent is being set for *PN

#### Lexeme Positions
I have Lexemes stubbed out for Line and Column, but they're not being used

#### Lexing and Parsing Errors
Have not touched this yet. When we find an error at either the Lexing or Parsing
stage, we should give as much data as possible. It should also be possible to
define the Lexer, Grammar and Parser to be more helpful when defining an error.

#### Parlex Struct
Right now, I just have parlex.Run - that should really be a struct with parser,
lexer and reduce.

This would also let me check that the symbols agree. Every symbol in the parser
grammar should either exist as a non-terminal in the parser grammar or a
terminal in the lexer. Every reduce rule should map to either a terminal or non-
terminal. There should also be no overlap between terminals and non-terminals.