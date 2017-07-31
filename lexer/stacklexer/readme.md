## Stack Lexer
The Stack Lexer can take some of the load off the parser and help simplify the
grammar. A stack lexer is comprised of several sub-lexers. Lexing rules can be
tagged to either push the current lexer on the stack and switch to different
lexer, or discard the current lexer and use one popped off the stack.

``` Sample Stack Lexer
== Main ==
  START innerLexer
  Shared
== InnerLexer ==
  STOP ^
  Shared
== Shared ==
  word  /\w+/
  space /\s+/ -
```

Each sub lexer is marked by two = followed by a word. The last two == are
optional and decorative.

### Lexing Rule Line Syntax
A lexer rule line has 4 parts.

  Name [/regex/ [(Groups)]] [^ | push-lexer] [-]

The last 4 parts are option and any combination can be included but the order
must be maintained.

The first part is the Symbol, which can be any non-whitespace characters except
"/".

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
a regex be defined. However, this isn't recommended anyway, the preferred method
is to always use names.

The second part also can optionally define which regex groups to use. This is
delimited by parenthesis. Any strings will be used literally and numbers will
identify the regex group.
```
== Main ==
  code /\s*(\d{3})[-:_](\d{3})/ (1|":"|2)
```

Will match "  123_456" but report the value as "123:456".

The third part is either the literal charcter "^" which indicates a pop
operation or the name of another sub-lexer which indicates a push operation.
More than one ^ can be used to indicate the number of layers of the stack to
pop.

The last segment is the literal character "-" which means that anything matching
should be discarded. Often useful to discard whitespace.

The order of rules indicates their priority.

### Inheiritance
When defining a line, if Name matches a sub-lexer, the rules of that sublexer
will be inheirited. This also means that a sub-lexer can be placed in a position
in the list to defin it's priority.

### Style Guide
It is recommended that subparsers be capitalized and symbols be lower-case
words. This is easier to read both in the lexer and when used later in a
grammar.