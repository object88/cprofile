# cprofile

Tool to read a process profile and produce useful debugging information

## Problem

Would like to be able to easily observe performance of Golang programs, sorting and filtering to discover hotspots, with a cross-platform interactive GUI.

Ideally, this application would read in a static profile file, and also be able to live monitor a process that's been tooled, either from its start, or for an indeterminate period.

## Notes

For setting up process monitoring over http, set up client process using [standard go directions](https://golang.org/pkg/net/http/pprof/), and make request to `http://localhost:8081/debug/pprof/goroutine`.  This will download a `.gz` file which can be fed into `go tool pprof`.

Need to determine whether repeated requests to `goroutine` are cumulative, etc.

## Possible tactic

Write add-in for profiling, not using `net/http/pprod`.  Instead, application must accept a command line flag to inform it of an endpoint, which it will use to open a web socket, and push profile information to.

Maybe provide option to indicate whether to start immediately, or wait to start.

## Building

`go build -o bin/cprofile main/main.go`

## Vendoring

Use the `dep` tool for vendoring.

`dep ensure`
