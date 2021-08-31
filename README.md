# dive

[![pkg.go.dev][gopkg-badge]][gopkg]

`dive` finds low readability if-blocks such as below.

* Too long if blocks
* Too many returns in a block
* A loop is in if-block 
* Deeply nest 

## Install

You can get `dive` by `go install` command (Go 1.16 and higher).

```bash
$ go install github.com/gostaticanalysis/dive/cmd/dive@latest
```

## How to use

`dive` run with `go vet` as below when Go is 1.12 and higher.

```bash
$ go vet -vettool=$(which dive) ./...
```

## Analyze with golang.org/x/tools/go/analysis

You can use [dive.Analyzer](https://pkg.go.dev/github.com/gostaticanalysis/dive/#Analyzer) with [unitchecker](https://golang.org/x/tools/go/analysis/unitchecker).

<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/gostaticanalysis/dive
[gopkg-badge]: https://pkg.go.dev/badge/github.com/gostaticanalysis/dive?status.svg
