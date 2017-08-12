## Parlex
[![GoDoc](https://godoc.org/github.com/AdamColton/parlex?status.svg)](https://godoc.org/github.com/AdamColton/parlex)

For a quick demo check out [json.go](https://github.com/AdamColton/parlex/blob/master/examples/parlexjson/json.go)

The core parlex package defines common interfaces around parsing. Sub-packages
exist to fulfill all of the interfaces.

The core package also implements some helpers than use only the interface
specification.

A few highlights; The
[stacklexer](https://github.com/AdamColton/parlex/tree/master/lexer/stacklexer)
provides a powerful, easy to use lexer that can reduce the complexity of a
grammar. The
[regexgram](https://github.com/AdamColton/parlex/tree/master/grammar/regexgram)
package supports some regex operators when defining a grammar. The
[packrat](https://github.com/AdamColton/parlex/tree/master/parser/packrat)
parser is a fairly efficient parser that can handle left recursion.