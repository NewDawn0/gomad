# gomad - A Go Monad Library for Typed Error Handling

## Overview

`gomad` is a Go library that provides a `TypedErrMonad`, a monad that stores either a value of type `T` or an error. The library allows you to bind the monad to functions that can return either just a value of type `T`, just an error, or both. Additionally, it provides the `ValueOr` function to retrieve the monad's value or a default value when an error occurs.

## Features

- **Typed Error Handling**: The monad stores a value of type `T` or an error.
- **`Bind` Method**: The `Bind` method allows you to chain operations, applying functions that return either a value of type `T`, an error, or both.
- **`ValueOr` Function**: The `ValueOr` function returns the stored value of the monad or a default value if an error occurred.

## Installation

To use the `gomad` library, simply copy the code into your Go project, or use `go get` to fetch the package.

```bash
go get github.com/NewDawn0/gomad
```

## Example Usage

```go
package main

import (
	"fmt"
	"github.com/NewDawn0/gomad"
)

// Function that returns a value and error
func a(x int) (int error) {
    if x == 42 {
        return 100, nil
    }
    return 0, fmt.Errorf("Number not 42")
}
// Errorous function
// Functions can accept arguments of anytype
fnError := func(x string) (int, error) {
    return 0, fmt.Errorf("something went wrong")
}

func main() {
	// Create a new monad with a value
	monad := gomad.NewTypedMonad(42)

	// Bind the function to the monad
	monad = monad.Bind(a, monad.Val)

	// Retrieve the result or a default value if error occurred
	result := monad.ValueOr(0)
	fmt.Println(result) // Output: 100

	// Bind the error function to the monad
	monad = monad.Bind(fnError, "test")

	// Retrieve the result or a default value if error occurred
	result = monad.ValueOr(0)
	fmt.Println(result) // Output: 0 (default value)
}
```

## Key Concepts

### `TypedErrMonad`

The core of the library is the `TypedErrMonad`, which has two fields:

- `Val`: Stores the value of type `T`.
- `Err`: Stores an error if one occurs.

The monad is used to represent a value that may or may not be valid, allowing you to chain operations that may fail without having to manually check errors at each step.

### Methods

#### `NewTypedMonad[T any](val T) TypedErrMonad[T]`

Creates a new `TypedErrMonad` containing a value of type `T`.

#### `Bind(f interface{}, args ...interface{}) TypedErrMonad[T]`

The `Bind` method allows you to bind a function to the monad. The function `f` must return either:

- A value of type `T`, or
- A value of type `T` and an error, or
- Just an error.

It handles the propagation of errors and returns a new monad based on the result of the function.

#### `ValueOr(val T) T`

The `ValueOr` method allows you to retrieve the value stored in the monad, or return a default value `val` if an error occurred.

## Error Handling

If an error occurs during a bind operation, the monad will propagate the error and prevent further operations. You can use `ValueOr` to provide a fallback value when an error is encountered.

### Example of Error Propagation

```go
fn := func(x int) (int, error) {
	return 0, fmt.Errorf("something went wrong")
}
fn2 := func(x int) (int, error) {
	return 0, fmt.Errorf("something went wrong here aswell")
}

monad := gomad.NewTypedMonad(42)
monad = monad.Bind(fn, monad.Val)
monad = monad.Bind(fn2, monad.Val)

// Will return the default value 0
result := monad.ValueOr(0)
fmt.Println(result) // Output: 0, Err: "Something went wrong"
```

## Function Signature

The function passed to `Bind` must adhere to one of the following signatures:

- `T` (returns only a value of type `T`)
- `(T, error)` (returns a value of type `T` and an error)
- `error` (returns only an error)

## Handling Invalid Functions

If the function does not return a valid signature (not `T`, `T` and `error`, or just `error`), the `Bind` method will return a monad containing an error.

```go
// Invalid function signature, will return an error
invalidFn := func(x int) string {
	return "Invalid"
}

monad = monad.Bind(invalidFn, monad.Val)
// This will return an error in the monad
```

## License

`gomad` is released under the [MIT License](https://opensource.org/licenses/MIT). You are free to use, modify, and distribute the code.
