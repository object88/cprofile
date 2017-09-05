# cprofile

Tool to gather useful information about a Go program

## Problem

Would like to understand the composition of a partial or complete Go program (library or executable) via AST exporation.

Would like to be able to easily observe performance of Golang programs, sorting and filtering to discover hotspots, with a cross-platform interactive GUI.

## Tools

This program exposes several tools

### Globals

Using the `globals` command, one can discover whether there are any instances of global variables in the code.

``` shell
cprofile globals main
```

### Imports

The `imports` command will list the imported packages.  See the [complete documentation](./cmd/imports.md).

## Building

`go build -o bin/cprofile main/main.go`

## Vendoring

Use the `dep` tool for vendoring.

`dep ensure`
