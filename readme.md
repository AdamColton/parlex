## Parlex
[![GoDoc](https://godoc.org/github.com/AdamColton/parlex?status.svg)](https://godoc.org/github.com/AdamColton/parlex)

For a quick demo check out [json.go](https://github.com/AdamColton/parlex/blob/master/examples/parlex_json)

The core parlex package defines a common language around parsing. It provides
interfaces for Lexer, Lexeme, Parser, Grammar, ParseNode and Reducer. It also
provides concrete types for Symbol, Production and Productions.

Sub-packages exist to fulfill all of the interfaces.

### Nil productions in Grammar
Some grammars make use of Nil or Epsilon values, that is a non-terminal that can
be skipped. The correct way to define this is:

Given:
```
NonTerminal -> NIL
NIL         ->
```

A NIL Production should be represented by a Production of length 0, not by nil.