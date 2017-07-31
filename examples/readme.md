## Parlex Examples

These examples should help show how parlex can be used.

The parlexmath example uses the simplelexer with the packrat parser. It can take
a simple math expression and evaluate it, taking into account order of
operations.

The parlexjson example uses the stacklexer with the topdown parser. It can take
a json string and format it or return an error indicating that there is a lexing
or parsing error.

The scalc example is meant to show a bit more involved example. It can take a
stack expression and evaluate it, returning a stack of precision floats.