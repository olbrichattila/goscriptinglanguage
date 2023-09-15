# Golang scripting language

## This project is under development, how to create your own language in golang.
This is not intented to be a full featured language (or who knows), only an example how a language interpreter works.

Current status:

Lexer implemented to tokenize the code
interpreter can evaluate experssions:

Example:

`go run .` will display the prompt

```
Aty programming language
-------------------------
-> 

```

### Try: Evaluate experssion:
Example:
```
1 + 5 + 6 /2 * (15 - 2)   
&{1 45}
```

Supported arithmetic operations, -, +, *, /, %

### Variable declaration:

```
let foo = 1;
```

### constant declaration:

```
const x = 8;

```

### Evaluate expression works with variables.
Example:

```
Aty programming language
-------------------------
-> x + foo * 3
&{1 11}
```

## What comes.
- Variable assignments,
- loops
- conditions
- functions
- closures
- classes
- structures
and more..






