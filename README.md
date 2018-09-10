# optionGen
---

optionGen is a tool to generate go Struct option for test, mock or more flexible.

## Installation
```bash
go get -u github.com/xsam/optionGen/cmd/optionGen
```

optionGen require [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) to format code which is generated. So you may confirm that `goimports` has been installed

```bash
go get golang.org/x/tools/cmd/goimports
```

## Example
To generate struct option, you need write a function declaration to tell optionGen how to generate. function Name consists of "_" prefix, struct name and "OptionDeclaration" suffix. In this function, just return a variable which type is `map[string]interface{}`.

The key of the map means option name, and the value of the map should consist of two parts, one for option type(except func type), and the other option default value.

Looks like this:
```go
"sounds": string("Meow")
"food":   (*string)(nil)
// function type is special
"Walk":   func() {
	log.Println("Walking")
}
```

A `Cat` struct option declaration may like this:
```go
func _CatOptionDeclaration() interface{} {
	return map[string]interface{}{
		"sounds": string("Meow"),
	}
}
```

Once you finished you declaration. You can add this line in your code
```go
//go:generate optionGen
```
and use `go generate` command to generate option.

Here is the sample result generate by `optionGen`

```go
package main

import "log"

type CatOptions struct {
	sounds string
}

type CatOp func(option *CatOptions)

func CatOpWith_sounds(value string) CatOp { return func(option *CatOptions) { option.sounds = value } }

func _NewCatOptions() CatOptions {
	return CatOptions{
		sounds: "Meow",
	}
}
```

To use the generated code. you could add `Options` struct to your struct

```go
type Cat struct {
	options CatOptions
}
```

And write a new function

```go
func NewCat(option ... CatOp) *Cat {
	cat := Cat{
		options: _NewCatOptions(),
	}

	for _, op := range option {
		op(&cat.options)
	}
	return &cat
}
```

For more example. see the [example](https://github.com/XSAM/optionGen/blob/master/example/cat.go) folder

Enjoy!