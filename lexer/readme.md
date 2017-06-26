## Lexer
[![GoDoc](https://godoc.org/github.com/AdamColton/parlex/lexer?status.svg)](https://godoc.org/github.com/AdamColton/parlex/lexer)

### Todo
* Rebuild old StackLexer with named lexers

StackLexer(map[string]string{
  "root":`
    START inner
    word /\w+/
    space /\s+/ -
  `,
  "inner":`
    STOP +
    innerWord /\w+/
    innerSpace /\s+/
  `
})

taking the form
name {state|+}? /regex/? -?