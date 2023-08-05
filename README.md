# unexported
`unexported` provides a Go analyzer module.
The analyzer is designed to check that exported functions and types use only exported types in their signatures. 
This helps to enforce best practices in Go code and ensures that the public API doesn't rely on unexported types.

## Usage
```
go get github.com/la0rg/unexported

import "github.com/la0rg/unexported"
```
The `unexported.NewAnalyzer()` follows guidelines in the
[`golang.org/x/tools/go/analysis`][xanalysis] package. This should make it
easy to integrate `unexported` with your own analysis driver program.

## Flags
The "unexported" analyzer supports the following flags:
-skip-interfaces: When set, interfaces are excluded from the analysis.
-skip-types: When set, types are excluded from the analysis.
-skip-func-args: When set, function arguments are excluded from the analysis.
-skip-func-returns: When set, function return parameters are excluded from the analysis.

## Command-Line Tool
This repository also includes a command-line tool to run the "unexported" analyzer on your code. 
The tool can be found in the cmd/unexported directory. To use the tool, follow these steps:

1. Install:
```
go install github.com/la0rg/unexported/cmd/unexported@latest
```

2. Run the tool passing the path to your Go code:
```
unexported [flags] [packages]
```

## Example
Given the package:
```go
package storage

type person struct {
	Name string
}

func GetPerson() person {
	return person{Name: "John Doe"}
}

func SavePerson(p person) {
	// save to DB
}

type ExtendedPerson struct {
	Person person
}
```

running the `unexported ./...` command with the default flags will produce:
```
storage.go:7:18: unexported type person is used in the exported function GetPerson
storage.go:11:17: unexported type person is used in the exported function SavePerson
storage.go:15:6: unexported type person is used in the exported type declaration ExtendedPerson
```

## Contributing
Contributions to this project are welcome! If you find any bugs, have feature requests, or want to contribute 
improvements, feel free to create issues and pull requests on the GitHub repository: https://github.com/la0rg/unexported