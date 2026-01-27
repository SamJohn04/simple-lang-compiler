# A Simple Language Compiler

A compiler for a simple language, to help better understand the practical side of compilers.

## Datatypes

As of right now, the compiler accepts:

- integers
- floating point numbers
- characters
- booleans

Strings are accepted only as the first parameter of an `output` call.

## Comments

Comments start with `//`, and are ignored by the lexer.

## Parser

For this compiler, an LL\(1\) parser has been chosen.

E.g. Program:
```
output "Enter a number: ";
let n = input;

let mut i = 2;
let mut fib1 = 1;
let mut fib2 = 1;

output "%d: %d\n", 0, fib1;
output "%d: %d\n", 1, fib2;

let mut temp;
while i < n {
    temp = fib1 + fib2;
    fib1 = fib2;
    fib2 = temp;
    output "%d: %d\n", i, temp;

    i = i + 1;
};
```
```
let n = input;
let di = input;

let mut i = 0;
while i < n {
    output "%d", i;
    i = i + di;
};

output "%d", i;
output "%d", n;
```

## Output

The compiler converts the given code to C-language. As such, do not use any keywords in C as identifiers.
