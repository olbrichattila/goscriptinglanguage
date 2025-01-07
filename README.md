# Golang Scripting Language

## This project is under development – How to create your own language in Go

This is not intended to be a fully featured language (though who knows what the future holds), but rather an example of how a language interpreter works.

### Current Status:
- Lexer implemented to tokenize the code
- Interpreter can evaluate expressions


Example:

`go run .` will display the prompt

```
Aty programming language
-------------------------
-> 

```

### Try: Evaluate expression:
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

### Variable assignments
Example:

```
Aty programming language
-------------------------
-> let x = 10;
x = x + 1

```

### Definition of complex objects:

Example:
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
numToStr(5)
strToNum("53")
input()
round(num, decimals)
rand(10) // parameter is optional, 0 to range
fileWrite(fileName, content)
fileRead(filename)

```
### Example of num to str
```
let mod = 0;
for (let i = 0; i < 100; i = i + 1) {
    mod = i % 3
    if (mod == 0) {
        print(numToStr(i) + " can be divided by 3")
    }

    if (mod != 0) {
        print(i)
    }
}
```

## Example input (reads from console)
```
print("What is your name? ")
let name = input();
print("Hello " + name)
```

## Example of file read, write
```
if (fileWrite("test.txt", "This is the test file content")) {
    println("File written succesfully")
}

let content = fileRead("test.txt");

println("The content is", content)
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

if (5 != 4) {
    print(5)
}

if (5 != 5) {
    print(5)
}
```

### Else
Examples:
```
if (1 == 2) {
    println("if")
} elseif (2 == 2) {
    println("elseif")
} elseif (2 == 4) {
    println("else2if")
} else {
    println("else3")
}


if (1 == 2) {
    println("if")
} elseif (2 == 3) {
    println("elseif")
} elseif (2 == 4) {
    println("else2if")
} else {
    println("else3")
}

if (1 == 1) {
    println("if1")
} else {
    println("if2")
}


if (1 == 2) {
    println("if3")
} else {
    println("if4")
}

```

### For loop
There are three variations of the for loop. One behaves like an incremental for loop, while the other two act as do-while and while loops. Therefore, there is no need to have separate while and do loops.

Example:
### Incremental for
```
for (let i = 0; i < 10 ; i = i + 1) {
    print(i)
}
```
### do for
```
let i = 0;
for (i < 10) {
    print(i)
    i = i + 1
}

```
### while for
```
let i = 0;
for {
    print(i)
    i = i + 1
}(i < 10)
```

## Break and continue
```
let i = 0;
for {
    i = i + 1
    if (i <= 10) {
        continue
    }

    if (i > 20) {
        break
    }
    
    println(i)
}
```

## Switch case:
### Numbers
```
for (let t = 1; t < 30; t = t + 1) {
    switch (t) {
        case 1:
            println("case 1")
            break
        case 3:
        case 4:
            println("case 3, 4")
            break
        case 5:
            println("case 5")
            break
        case 15:
            println("case 15")
            break
        default:
            println("  default")
    }
}
```
### Strings
```
let s = "test";
switch (s) {
    case "":
        println("empty string")
        break
    case "test":
        println("test string")
        break
    case "nontest":
        println("case 15")
        break
    default:
        println("  string default")
}

```

## String assignment and comparision
```
const a = "Arnold";
const b = "Bruno";

print(a + " " + b)
print ((a < b))
print ((a > b))
print ((a == b))
print ((a != b))
print ((a <= b))
print ((a >= b))
```

## What comes.

- classes

and more..

## About me:
Learn more about me on my personal website. https://attilaolbrich.co.uk/menu/my-story
Check out my latest blog blog at my personal page. https://attilaolbrich.co.uk/blog/1/single

