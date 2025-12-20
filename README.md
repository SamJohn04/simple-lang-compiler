# A Simple Language Compiler

A compiler for a simple language, to help better understand the practical side of compilers.

## Parser

For this compiler, an LL\(1\) parser has been chosen. The parser uses the following grammar.

```text
I -> I1;I | ε

I1 -> variable=E |
    let I6 |
    if R { I } I4 |
    while R { I } |
    output E

I4 -> else I7 | ε
I7 -> if R { I } I4 | { I }

I6 -> mut variable I8 | variable=E
I8 -> =E | ε

R -> ER1E
R1 -> > | < | == | != | >= | <=

E -> TE1
E1 -> +TE1 | -TE1 | ε
T -> FT1
T1 -> *FT1 | /FT1 | %FT1 | ε
F -> input | variable | literal | (E) | -F
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
