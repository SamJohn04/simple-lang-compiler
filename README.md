# A Simple Language Compiler

A compiler for a simple language, to help better understand the practical side of compilers.

```text
term = <input> | variable | literal
expression = E
relation = R
assign = variable = expression | variable = relation
instr = assign |
    <if> relation { instr } |
    <if> relation { instr } <else> { instr } |
    <if> relation { instr } <else> <if> relation { instr } ... |
    <while> relation { instr } |
    <output> expression |
    <output> relation

R  -> E>E | E<E | E==E | E!=E | E>=E | E<=E | !R | E
E  -> TE'
E' -> +TE' | -TE' | epsilon
T  -> FT'
T' -> *FT' | /FT' | %FT' | epsilon
F  -> term | (E) | epsilon
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
