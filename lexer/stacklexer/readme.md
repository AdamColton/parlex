## Stack Lexer
[![GoDoc](https://godoc.org/github.com/AdamColton/parlex/lexer/stacklexer?status.svg)](https://godoc.org/github.com/AdamColton/parlex/lexer/stacklexer)

The Stack Lexer can take some of the load off the parser and help simplify the
grammar. A stack lexer is comprised of several sub-lexers. Lexing rules can be
tagged to either push the current lexer on the stack and switch to different
lexer, or discard the current lexer and use one popped off the stack.

``` Sample Stack Lexer
== Main ==
  start innerLexer
  outerFoo /foo/
  Shared
== InnerLexer ==
  stop ^
  innerFoo /foo/
  Shared
== Shared ==
  word  /\w+/
  space /\s+/ -
```

Will lex
```
foo bar start foo bar baz stop foo
```
as
```
outerFoo:"foo" word:"bar" start:"start" innerFoo:"foo" word:"bar" word:"baz" stop:"stop"  outerFoo:"foo"
```

Each sub lexer is marked by == followed by a word. The trailing == is
optional and decorative.

### Lexing Rule Line Syntax
A lexer rule line has 4 parts.

  Name [/regex/ [(Groups)]] [^ | push-lexer] [-]

The last 4 parts are optional and any combination can be included but the order
must be maintained.

The first part is the Symbol, which can be any non-whitespace characters except
"/". However, the best practice is to use camel-case words beginning with lower
case for the symbols and camel-case words beginning with upper case for the
lexer names.

The second part is the regular expression which starts and ends with "/". If no
regular expression is given, the Symbol is used as the regular expression. This
means that a Symbol that is a regular expression operator will cause an error
and needs a regular expression manually defined. For example

```
== Main ==
  x
  + /\+/
```
The first will match "x" and the second will match "+", but the second requires
a regex be defined. However, this isn't recommended anyway, because it
violates the naming standard.

The second part also can optionally define which regex groups to use. This is
delimited by parenthesis. Any strings will be used literally and numbers will
identify the regex group.
```
== Main ==
  code /\s*(\d{3})[-:_](\d{3})/ (1|":"|2)
```

Will match "  123_456" but report the value as "123:456".

The third part is either the literal character "^" which indicates a pop
operation or the name of another sub-lexer which indicates a push operation.
More than one ^ can be used to indicate the number of layers of the stack to
pop.

The last segment is the literal character "-" which means that anything matching
should be discarded. Often useful to discard whitespace.

The order of rules indicates their priority.

### Inheritance
When defining a line, if Name matches a sub-lexer, the rules of that sublexer
will be inherited. This also means that a sub-lexer can be placed in a position
in the list to define it's priority.
