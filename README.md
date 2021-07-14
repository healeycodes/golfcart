# ⛳ Golfcart

[![Go](https://github.com/healeycodes/golfcart/actions/workflows/go.yml/badge.svg)](https://github.com/healeycodes/golfcart/actions/workflows/go.yml)

Golfcart is a minimal programming language inspired by Ink, JavaScript, and Python – implemented in Go.

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

It's a dynamic strongly typed language with support for bools, strings, numbers (float64), lists, dicts, and nil (null). There is full support for closures and functions can alter any variable in a higher scope.

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

A Golfcart program is a series of expressions. Linebreaks are optional and there are no semi-colons. The final expression is sent to stdout.

```javascript
a = 1 b = 2 assert(a + b, 3)
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

// Dicts
{a: 1} // Accessed by `.a` or `["a"]` like JavaScript
keys({a: 1}) // ["a"]

// Functions
_ => nil // All user-defined functions are anonymous, assignable by variable
n => n + 1
sum = (x, y) => x + y

// Nil
nil
nil == nil // true
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
- [Programs that purposefully throw errors](https://github.com/healeycodes/golfcart/tree/main/example%20programs/error%20programs)

(All the above are used as part of Golfcart's test suite)

## Motivations

This is a toy programming language that I built to use for Advent of Code 2021. Another motivation was to learn how to write an interpreter from scratch. Previously, I read Crafting Interpreters and [implemented the Lox programming language](https://github.com/healeycodes/hoot-language) using Python, and [partially ported Ink](https://github.com/healeycodes/quill) using Rust. Another introduction to interpreters I enjoyed was [A Frontend Programmer's Guide to Languages](https://thatjdanisso.cool/programming-languages). The [Ink blog](https://dotink.co/posts/) is also great.

I wanted to design a programming language that didn't use semi-colons or automatic semicolon insertion. So, no statements and everything should be an expression that evaluates to a value. For example:
- `if/else if/else` evaluates to the successful branch
- a variable declaration evalutes to the value
- setting a dict value evalutes to the value
- a for loop evalutes to the number of times the condition expression succeeded

```javascript
assert(
    for i = 0; i < 5; i = i + 1 {}, 5
)
```

However, I didn't realise how restrictive my design goals were. A problem I ran into early was accessing an item from a literal.

```javascript
[1][0] // This evaluates to [0] because Golfcart thinks it's two lists

// Instead, you do:
a = [1]
a[0]
```

While it's too late to add semi-colon separated statements to Golfcart, I have a new found appreciation for `;`.

The main problem with Golfcart is that there are differences between how Golfcart programs run in my head vs. in the interpreter. This is because I jumped to implementing the language and didn't spend enough time designing. Linus Lee has some interesting notes on designing small interpreters in [Build your own programming language](https://thesephist.com/posts/pl/#impl).

> In this phase, I usually keep a text document where I experiment with syntax by writing small programs in the new as-yet-nonexistent language. I write notes and questions for myself in comments in the code, and try to implement small algorithms and data structures like recursive mathematical functions, sorting algorithms, and small web servers.

If I had more predefinied programs to start with (to run as tests), I would have noticed the divergence of how programs are actually evaluated early enough to re-think the design. This project's example programs were written after the fact within the confines of the language's limitations.

Ultimately, I've learned a lot and this won't be my last language!

## Implementation

Golfcart is a tree-walk interpreter. Its one dependancy is the [Participle](https://github.com/alecthomas/participle) parsing library, which consumes a parser grammer written using Go structs and a RegEx-like syntax to create a syntax tree (see [parser.go](https://github.com/healeycodes/golfcart/blob/main/pkg/golfcart/parse.go)). This library let me move fast and refactor parsing bugs without headaches.

A piece of source code is turned into tokens by Participle's lexer. The lexer uses token definitions. For example, Golfcart's identifier RegEx is: ```{"Ident", `[\w]+`, nil}```). These tokens are parsed into a syntax tree using struct definitions.

Here's a list literal:

```go
type ListLiteral struct {
	Pos lexer.Position

	Expressions *[]Expression `"[" ( @@ ( "," @@ )* )? "]"`
}
```

Once the source code has been built into the syntax tree, each node of this tree is walked — as in _tree-walk_ (see [eval.go](https://github.com/healeycodes/golfcart/blob/main/pkg/golfcart/eval.go)). The code archiecture is similar to [Ink](https://github.com/thesephist/ink)'s — the way stack frames work is similar, and I used a near-identical interface for values. 

```go
type Value interface {
	String() string
	Equals(Value) bool, error
}
```

<br>

Whenever I implemented a new type of value or syntax, I added a specification program with an assertion. When I came across a bug, I sometimes wrote an error program to purposefully throw an error. This project's tests `go test ./...` ensure that the specification programs and example programs run without any errors (an `assert()` call throws an error and quits) and that the error programs all error out.

The Participle library provides line-numbers for each lexer token. These are added to the value structs during evaluation and so some Golfcart errors have line-numbers (all have error text, and hopefully enough information to find and fix). Golfcart lacks a mature stack trace.

There are runtime functions (e.g. input/output, type assertions and casts, keys/values, etc.) in [runtime.go](https://github.com/healeycodes/golfcart/blob/main/pkg/golfcart/runtime.go) and the REPL can be found in [run.go](https://github.com/healeycodes/golfcart/blob/main/pkg/golfcart/run.go).


## Contributions/issues

More than welcome! Raise an issue with a bug report/feature proposal and let's chat.

## License

MIT.

