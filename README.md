# gostatic
Static website generator.

[![Go Report Card](https://goreportcard.com/badge/github.com/esdinb/gostatic?style=flat-square)](https://goreportcard.com/report/github.com/golang-standards/project-layout)

# WIP

Based on [libxslt](https://gitlab.gnome.org/GNOME/libxslt) and [libxml2](https://gitlab.gnome.org/GNOME/libxml2). 

Uses [yuins goldmark](https://github.com/yuin/goldmark) Markdown parser.

Uses [evanws esbuild](https://github.com/evanw/esbuild) bundler.

CGO wrappers for libxml2 are based on [jbussdiekers golibxml](https://github.com/jbussdieker/golibxml) code.

## Requirements

Requires Go version 1.18 or later.

The libxml2 and libxslt libraries needs to be installed on the system.

The [lexbor](https://github.com/lexbor/lexbor) libraries needs to be installed on the system.

C library bindings uses CGO and require a C compiler to also be installed.

## Building

Clone the gostatic repository:

`git clone https://github.com/esdinb/gostatic.git`

In the gostatic directory run `go build` something like this:

`CGO_CFLAGS="-I/usr/local/include" CGO_LDFLAGS="-L/usr/local/lib -llexbor_static" go build -o gostatic cli/main.go`

Change `CGO_CFLAGS` and `CGO_LDFLAGS` to match the paths to the lexbor libraries.

