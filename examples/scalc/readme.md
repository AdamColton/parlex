## Stack Calculator

### Numbers

Numbers can begin with + or -. They must have an integer part and may have a
decimal part. The out put will be limited to the longest decimal part.

### Uniary Operators
-- returns the negative value

drop drops the top value on the stack.

@ copies the next symbol on the stack

! if the next value is <=0 it is replaced with 1, otherwise 0.

### Binary operators
+ Addition

- Subtraction

* Multiplaction

/ Division

^ Exponent

% takes the modulus

= Puts a 1 on the stack if the next two values are equal, and a zero if they
are not

comp compares the next two values and puts a 1 if the first is greater, a 0 if
they are equal and -1 if the second is greater.

swap takes the next two symbols on the stack and swaps them.

### Stack Operators
Stack operators consume an entire stack. A sub stack can be defined by
surrounding it in parenthesis.

sum will sum the entire stack

len will return the length of the stack.

avg will average the values in the stack.

min will return the minimum value in the stack

max will return the maximum value in the stack

### Stack Manipulation operators
These may both consume some of the stack and add values to the stack.

if will take the next value and if it is a 0