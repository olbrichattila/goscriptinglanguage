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

### Variable assingnmets
Example:

```
Aty programming language
-------------------------
-> let x = 10;
x = x + 1

```

### Definion of complex objects:

Exampe:
```
let foo = 15 - 3;
const obj = {
   x: 150,
   y: 130,
   foo,
   complex: {
    bar: true,
   },
};

print(1,5)
let f = obj.complex.bar;
foo = obj.foo() + 5

```

### internal Functions
```
print(1, 5)
time()
```

### User defined functions Functions, closures with variable scopes
Example:
```
let z = 35;
fn add (x, y) {
    let result = x + y;
    print(result)

    result
}

fn sub () {
    let x = 10;
    let y = 20;
    fn add (x, y) {
        let result = x + y;
        print(result)

        result
    }

    let foo = 45;
    add(x,foo)
}

const result = add(10, 4);

print(result)
print(result)
print(result)

print(add(5, 3))
print(add(2, 7))
print(add(5, 3))
print(add(1, z))


print(sub())
```

## What comes.
- loops
- conditions
- classes

and more..






