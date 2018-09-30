# goportfinder

Trying out a portscanner in Golang.

## Usage

``` go

./goportfinder:
  -ip string
        insert ip (default "127.0.0.1")

```

## Installation

`go get github.com/x1mdev/goportfinder`

### Example

` ./goportfinder -ip=127.0.0.1 `

This will scan 3000 ports by default on the target provider with the `-ip` flag.

It's fast.