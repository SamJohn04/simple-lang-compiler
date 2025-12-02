# A Simple Language Compiler

A compiler for a simple language, to help better understand the practical side of compilers.

```text
term = <input> | variable | literal
expression = E
relation = R
assign = variable = expression | variable = relation
instr = assign |
    <if> relation { INS } |
    <if> relation { INS } <else> { INS } |
    <if> relation { INS } <else> <if> relation { INS } ... |
    <while> relation { INS } |
    <output> expression |
    <output> relation

INS -> instr;INS | epsilon

R -> !R | E>E | E<E | E==E | E!=E | E>=E | E<=E
E -> E+T | E-T | T
T -> T*F | T/F | T%F | F
F -> term | (E)

PROGRAM -> INS
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
}

output i;
output n;
```
