# A Simple Language Compiler

A compiler for a simple language, to help better understand the practical side of compilers.

## I/O

Access to the C `printf` and `scanf("%lld")` are provided using the `output` and `input` keywords respectively.

## Datatypes

As of right now, the compiler accepts:

- integers
- floating point numbers
- characters
- booleans
- arrays

Strings are accepted only as the first parameter of an `output` call.

## Integers

Integers are represented by the type `long long int` in C.

## Floating Point Numbers

Floating point numbers are represented by the type `double` in C.

## Characters

Characters are represented by the type `char` in C.

## Booleans

Booleans are to import the stdbool C-library, and use the type `bool`.

## Arrays

Arrays are to be used like so:
```
let arr = [1, 2, 3];
```

The array index begins at 0. All elements of the array must be of the same type.

Multidimensional arrays are represented as an array of arrays. All the inner arrays must have an equal number of elements of the same type.

## Comments

Comments start with `//`, and are ignored by the lexer.

Inline comments are possible. E.g.:
```
output "Hello World!\n"; // prints hello world
```

## Parser

For this compiler, an LL\(1\) parser has been chosen.

E.g. Program:
```
output "Enter a number: ";
let n = input;

let mut i = 2;
let mut fib1 = 1;
let mut fib2 = 1;

output "%lld: %lld\n", 0, fib1;
output "%lld: %lld\n", 1, fib2;

let mut temp;
while i < n {
    temp = fib1 + fib2;
    fib1 = fib2;
    fib2 = temp;
    output "%lld: %lld\n", i, temp;

    i = i + 1;
};
```
```
let n = input;
let di = input;

let mut i = 0;
while i < n {
    output "%lld", i;
    i = i + di;
};

output "%lld", i;
output "%lld", n;
```

## Output

The compiler converts the given code to C-language.
