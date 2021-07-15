# ⛳ Golfcart

[![Go](https://github.com/healeycodes/golfcart/actions/workflows/go.yml/badge.svg)](https://github.com/healeycodes/golfcart/actions/workflows/go.yml)

 * [Getting Started](#getting-started)
 * [Scope Rules](#scope-rules)
 * [Usage](#usage)
 * [Building and tests](#building-and-tests)
 * [Contributions](#contributions)
 * [License](#license)

Golfcart is a minimal programming language inspired by Ink, JavaScript, and Python – implemented in Go. It's a toy programming language that I built to use for Advent of Code 2021. Another motivation was to learn how to write an interpreter from scratch.

```javascript
// Here's the classic interview question FizzBuzz
for i = 1; i < 101; i = i + 1 {
    log(if i % 3 == 0 and i % 5 == 0 {
       "FizzBuzz"
    } else if i % 3 == 0 {
       "Fizz"
    } else if i % 5 == 0 {
       "Buzz"
    } else {
       str(i)
    })
}
```

Golfcart is a dynamic strongly typed language with support for bools, strings, numbers (float64), lists, dicts, and nil (null). There is full support for closures and functions can alter any variable in a higher scope.

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

For Golfcart, I began with a desire to design a small programming language that didn't use semi-colons or automatic semicolon insertion. So, no statements, and everything should be an expression that evaluates to a value. For example:
- `if/else if/else` evaluates to the successful branch
- A variable declaration evaluates to the value
- Setting a dict value evaluates to the value
- A for loop evaluates to the number of times the condition expression succeeded

```javascript
assert(
    // This runs five times
    for i = 0; i < 5; i = i + 1 {}, 5
)
```

## Getting Started

A Golfcart program is a series of expressions. Line breaks are optional and there are no semi-colons. The final expression is sent to stdout.

```javascript
a = 1 b = 2 assert(a + b, 3) // A successful assert() evaluates to nil
```

There are seven types. A type-check can be performed with `type()`.

```javascript
// Bools
true or false
true and true

// Numbers
1
1.1 + 1.1 // 2.2

// Strings
"multi-line
string"
"1" + "2" // "12"

// Lists
[1, 2]
nums = [3, 4]
nums.append(5) // [3, 4, 5]
[0] + [1] // [0, 1]

// Dicts
{a: 1} // Accessed by `.a` or `["a"]` like JavaScript
       // Values can be any type
keys({a: 1}) // ["a"]

// Functions
_ => nil // All user-defined functions are anonymous, assignable by variable
n => n + 1
sum = (x, y) => x + y

// Nil
nil
nil == nil // true
```

The Fibonacci sequence.

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
- [Programs that purposefully throw errors](https://github.com/healeycodes/golfcart/tree/main/example%20programs/error%20programs)

(All the above are used as part of Golfcart's test suite)

## Scope Rules

Let's talk about stack frames in Golfcart. A stack frame is a map of variables in scope. It's a recursive structure, every stack frame has a parent apart from the global frame. All functions are anonymous and create closures. Any variable referenced in a higher scope can be altered. Examples explain this better than words.

```javascript
a = 1
a_function = () => a = 2 // Closure created
a_function() // When called, `a` is changed
a // 2

if true {
    // `b` is not defined in a higher scope
    // So, `b` is declared only within this scope
    b = 3
}
b // Error: cannot find value for key 'b'

c = nil
if true {
    // This assignment recursively looks in higher scopes for `c`
    // it's found and that value is altered
    c = 4
}
c // 4
```

## Usage

Pass a Golfcart program as the first command-line argument

```bash
$ ./golfcart-linux program.golf
```

Run the binary with no command-line arguments to open the REPL.

```bash
$ ./golfcart-linux 

      .-::":-.
    .'''..''..'.
   /..''..''..''\
  ;'..''..''..''.;
  ;'..''..''..'..;
   \..''..''..''/
    '.''..''...'
      '-..::-' Golfcart v0.1
λ 
```

Use `-ebnf` to print the Extended Backus–Naur form grammar to stdout and quit.

Use `-version` to print the version to stdout and quit.

## Building and tests

Create releases.

```bash
./build.sh
```

Run all tests (they also run via GitHub Action on commit).

```bash
go test ./...
```

## Contributions

More than welcome! Raise an issue with a bug report/feature proposal and let's chat.

## License

MIT.

