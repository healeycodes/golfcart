# ⛳ golfcart

[![Go](https://github.com/healeycodes/golfcart/actions/workflows/go.yml/badge.svg)](https://github.com/healeycodes/golfcart/actions/workflows/go.yml)

golfcart is a minimal programming language inspired by Ink, JavaScript, and Python. It's implemented in Go.

```javascript
for i = 1; i < 101; i = i + 1 {
    if i % 3 == 0 and i % 5 == 0 {
        log("FizzBuzz")
    } else if i % 3 == 0 {
        log("Fizz")
    } else if i % 5 == 0 {
        log("Buzz")
    } else {
        log(str(i))
    }
}
```

golfcart is a dynamic typed language with support for bools, strings, numbers (float64), lists, dicts, and nil (null). There is full support for closures and functions can alter any variable in a higher scope.

```javascript
counter = () => {
    n = 0
    () => {
        n = n + 1
        n
    }
}

my_counter = counter()
my_counter() // 1

assert(my_counter(), 2)
```

## Getting started

A golfcart program is a series of expressions. Linebreaks are optional and there are no semi-colons.

```javascript
a = 1 b = 2 assert(a + b, 3)
```

There are seven types. A type-check can be performed with `type()`.

```javascript
// bools
true or false
true and true

// numbers
1
1.1

// strings
"multi-line
string"

// lists
[1, 2]

// dicts
{a: 1} // accessed by `.a` or `["a"]` like JavaScript

// functions
() => nil
n => n + 1
sum = (x, y) => x + y

// nil
nil
```

The fibonacci sequence.

```javascript
// Naive
t = time()
fib = n => if n == 0 {
    0
} else if n == 1 {
    1
} else {
    fib(n - 1) + fib(n - 2)
}
fib(20)
log("fib: " + str(time() - t))

// With memoization 
t = time()
cache = {"0": 0, "1": 1}
fib_memo = n => if cache[n] != nil {
    cache[n]
} else {
    cache[n] = fib_memo(n - 1) + fib_memo(n - 2)
}
fib_memo(20)
log("fib_memo: " + str(time() - t))
```

For more detailed examples, see:
- [Example programs](https://github.com/healeycodes/golfcart/tree/main/example%20programs)
- [Specification programs](https://github.com/healeycodes/golfcart/tree/main/example%20programs/spec%20programs)
- [Programs that purposefully throw errors](https://github.com/healeycodes/golfcart/tree/main/example%20programs/error%20programs).

## Motivations

This is a toy programming language that I built to take part in Advent of Code 2021 using my own programming language. And to learn how to write an interpreter from scratch. Previously, I read Crafting Interpreters and implemented the Lox programming language using Python, and partially ported Ink using Rust.

I wanted to design a programming language that didn't use semi-colons or automatic semicolon insertion. So, no statements and everything should be an expression that evaluates to a value.
- `if` evaluates to the successful branch
- a variable declaration evalutes to the value
- setting a map value evalutes to the value
- a for loop evalutes to the number of times the condition expression succeeded

```javascript
assert(
    for i = 0; i < 5; i = i + 1 {}, 5
)
```

I didn't realise how restrctive my design goals were. A problem I ran into early was accessing an item from a literal.

```javascript
[0] // A list with the number zero
[0] // This is a parsing error!
    // Why? Well, are we evaluating another list
    // or trying to access the zeroth element of the literal on line 1

// Instead, you do:
a = [0]
a[0] // Evaluates to `0`

// or
([0])[0] // Evaluates to `0`
```
