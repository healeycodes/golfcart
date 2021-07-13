# ⛳ golfcart

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

It's a dynamic typed language with support for bools, strings, numbers (float64), lists, dicts, and nil (null). There is full support for closures and functions can alter any variable in a higher scope.

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

Next, the fibonacci sequence.

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
