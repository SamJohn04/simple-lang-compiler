# A Simple Language Compiler

A compiler for a simple language, to help better understand the practical side of compilers.

## Parser

```text
INS -> I;INS | epsilon

I -> variable = E |
    if R { INS } J |
    while R { INS } |
    output E

J -> else M |
    epsilon

M -> if R { INS } J |
    { INS }

R -> ER'
R' -> >E | <E | ==E | !=E | >=E | <=E

E -> TE'
E' -> +TE' | -TE' | epsilon
T -> FT'
T' -> *FT' | /FT' | %FT' | epsilon
F -> input | variable | literal | (E)
```

E.g. Program:
```
let a = 10;
let b = 20;

let out = a + (b * 10)/input;
output out;
```
```
let n = input;
let di = input;

let mut i = 0;
while i < n {
    output i;
    i = i + di;
};

output i;
output n;
```
