# A Simple Language Compiler

A compiler for a simple language, to help better understand the practical side of compilers.

```text
term = <input> | variable | literal
expression = E
relation = R
assign = variable = expression | variable = relation
instr = assign | <if> relation <then> instr | <if> relation <then> instr <else> instr | <goto> :label | <output> expression | <output> relation
label = :label

R  -> E>E | E<E | E==E | E!=E | E>=E | E<=E | !R
E  -> TE'
E' -> +TE' | -TE' | epsilon
T  -> FT'
T' -> *FT' | /FT' | %FT' | epsilon
F  -> term | (E) | epsilon
```

E.g. Program:
```
a = 10;
b = 20;

out = a + (b * 10)/input;
output out;
```
```
n = input;
di = input;

i = 0;
:l1;
if i >= n then goto :l2;
output i;
i = i + di;
goto :l1;

:l2 output i;
output n;
```
