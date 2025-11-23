# wbmath

wbmath is a Go utility library that provides a collection of general-purpose 
math helpers and a rational-number implementation (`Fraction`) for precise 
arithmetic with fractions. More features may be added over time.

Module path: `github.com/bogersw/wbmath`

## Overview

The library contains:

- General math helpers in the package `wbmath` (examples: `Gcd`, `PowInt`, `Round`, `IsInteger`).
- A `fraction` subpackage that implements a `Fraction` type and utilities for creating 
and manipulating rational numbers (constructors, arithmetic operations, simplification, 
string formatting, evaluation to float, etc.).

The API is intentionally small and idiomatic Go. The `fraction` package stores numerators and denominators as non-negative integers and tracks sign separately.

## Install

Use `go get` (or add the module to your `go.mod`):

```bash
go get github.com/bogersw/wbmath
```

Then import where you need it:

```go
import (
    "fmt"
    "github.com/bogersw/wbmath"
    "github.com/bogersw/wbmath/fraction"
)
```

## Examples

Below are several short examples that demonstrate common use cases.

### Gcd (greatest common divisor)

The `Gcd` function computes the greatest common divisor using the Euclidean algorithm.

```go
package main

import (
    "fmt"
    "github.com/bogersw/wbmath"
)

func main() {
    fmt.Println(wbmath.Gcd(48, 18)) // Output: 6
    fmt.Println(wbmath.Gcd(0, 5))   // Output: 5
}
```

### Fractions (basic usage)

The `fraction` package exposes constructors and methods to work with rational numbers.

```go
package main

import (
    "fmt"
    "github.com/bogersw/wbmath/fraction"
)

func main() {
    // Create a fraction from integers
    f, err := fraction.New(2, 3)
    if err != nil {
        panic(err)
    }
    fmt.Println(f.String()) // "2/3"

    // Create from a string
    f2, err := fraction.NewFromString("4 / 6")
    if err != nil {
        panic(err)
    }
    // NewFromString parses and reduces where possible
    fmt.Println(f2.String()) // "2/3" (after simplification)

    // Arithmetic: add
    f.Add(f2).Simplify()
    fmt.Println(f.String()) // "1"

    // Multiply
    f3 := fraction.MustNew(3, 4)
    f3.Multiply(fraction.MustNew(2, 5)).Simplify()
    fmt.Println(f3.String()) // "3/10"

    // Evaluate to float
    fmt.Println(f3.Evaluate()) // 0.3

    // Inspect numerator/denominator
    if num, ok := f3.Numerator(); ok {
        fmt.Println("numerator:", num)
    }
}
```

Notes:
- Use `MustNew*` constructors when you want a panic on invalid input; 
otherwise use the error-returning constructors and handle errors.
- Many `Fraction` methods modify the receiver in-place and return the 
receiver to allow method chaining.

## Testing

Run the package tests with the standard Go tooling:

```bash
# run all tests in the module
go test ./...

# run only the wbmath package tests (verbose)
cd wbmath && go test -v
```

To see detailed output for the tests use `-v` as shown above.

## Contributing

Contributions are welcome. Please open issues or pull requests on the repository. Keep code idiomatic and include tests for new functionality.

## License

This project is licensed under the MIT License. See the `LICENSE` file for the 
full license text or visit <https://opensource.org/licenses/MIT>.
