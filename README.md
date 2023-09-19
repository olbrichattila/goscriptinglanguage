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

## Conditional expressions
```
let x;
x = (1 == 1)
print(x)

x = (1 > 1)
print(x)

x = (1 < 1)
print(x)

x = (5 + 5 == 10)
print(x)

const foo = 51;
x = (foo == 51)
print(x)

x = (foo -1 == 50)
print(x)

x = (foo > 0)
print(x)
```

## If statement 
```
if (5 == 5) {
    print(true)
}

if (5 + 5 > 7) {
    print(true)
}

if (5 + 5 < 7) {
    print(false)
}

const foo = 50;

if (foo == 50) {
    print(true, 50)
}

if (foo > 10) {
    print(true, 10)
}

if (foo >= 10) {
    print(true, 10)
}

if (foo < 10) {
    print(false, 10)
}

if (foo <= 10) {
    print(false, 10)
}

let z = 20;

if (z == 20) {
    print(20)
}

z = 25
if (z == 25) {
    print(20)
    print(25)
}
```

## What comes.
- else, if else
- loops
- classes

and more..

