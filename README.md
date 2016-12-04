go-fish
======

A [><>](http://esolangs.org/wiki/Fish) interpreter written in Go.

Installation
---------------

Install [golang](http://golang.org/doc/install). To install or update go-fish on your system, run:

```
go install github.com/redstarcoder/go-fish
```

Usage
---------------

```
$ go-fish -h
Usage: go-fish [args] <file>
  -c	output the codebox each tick
  -h	display this help message
  -i value
    	set the initial stack (ex: '"Example" 10 "stack"')
  -m	run like the fishlanguage.com interpreter
  -s	output the stack each tick
  -t duration
    	time to sleep between ticks (ex: 100ms)
```

Acknowledgments
---------------

* [redstarcoder](https://github.com/redstarcoder) wrote this library.
