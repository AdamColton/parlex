## Stack Calculator

### Numbers

Numbers can begin with + or -. They must have an integer part and may have a
decimal part. The out put will be limited to the longest decimal part.

### Unary Operators
The Unary operators are -- and abs. -- performs negation and abs takes the
absolute value.

### Binary operators
As expected +, -, *, / as well as ^ for exponent, and % for modulus. There are
also a set of binary operations for comparison >,< and = will return 1 for true
and 0 for false. The cmpr operator will return 1 if a>b, 0 if a==b and -1 if
a&lt;b.

### Ternary Operator
? acts as the ternary operator. It is always of the form
  T T E ?
The first two can be anything as long as the types match; unary operators,
binary operators, stack operators, stack manipuation operators or expressions.
E must be an expression. If E evalutes to >0, the second T (top of the stack)
will be returned, other wise the first T will be returned.

### Stack Operators
Stack operators consume an entire stack and return a value. A sub stack can be
defined by surrounding it in parenthesis.

  1 (2 3 sum) -> 1 5

The stack operators are len, sum, avg, min, max, first, last. It should be
obvious what they do.

### Stack Manipulation operators
These change the structure of the stack. They are swap, drop and clear. Swap
will swap the top two values. Drop will drop the first value. Clear will clear
the whole stack, which is useful in interactive mode.

### Precision
The degree of precision is tracked and used in the return values.
  1 + 2.0 -> 3.0
  1.0 / 3 -> 0.3

One gotcha to be aware of is that this uses the Go rounding built into fmt. It
has one piece of unexpeced behavoir; when rounding a trailing 5, if the
preceding digit is odd it rounds down and if it is even, it rounds up.
  (25 10 /) (15 10 /) -> 2 2

### Command line
The command line tool is "scalc". Running it with no input will enter
interactive mode. Type "exit" to exit. Running scalc with input will evaluate
the input. Running "scalc parse [expression]" will show the parse tree for the
expression.