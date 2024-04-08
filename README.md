# gostatic
Static website generator.

[![Go Report Card](https://goreportcard.com/badge/github.com/esdinb/gostatic?style=flat-square)](https://goreportcard.com/report/github.com/golang-standards/project-layout)

# WIP

## Requirements

Requires Go version 1.21 or later.

### Go libraries:

- [goldmark](https://github.com/yuin/goldmark)
- [esbuild](https://github.com/evanw/esbuild)

### C libraries:

- [libxml2](https://gitlab.gnome.org/GNOME/libxml2)
- [libxslt](https://gitlab.gnome.org/GNOME/libxslt)
- [lexbor](https://github.com/lexbor/lexbor)

C library bindings uses CGO and require a C compiler to also be installed.

## Building

Clone the gostatic repository:

`git clone https://github.com/esdinb/gostatic.git`

In the gostatic directory run `go build` something like this on macOS:

`CGO_CFLAGS="-I/usr/local/include" CGO_LDFLAGS="-lxml2 -lxslt -lexslt -L/usr/local/opt/lexbor/lib -llexbor_static" go build -o gostatic cli/main.go`

Or using CMake:

`cmake -DCGO_CFLAGS="-I/usr/local/opt/lexbor/include" -DCGO_LDFLAGS="-lxml2 -lxslt -lexslt -L/usr/local/opt/lexbor/lib -llexbor_static" . && make`

Change `CGO_CFLAGS` and `CGO_LDFLAGS` to match the paths to the libxml2, libxslt and lexbor libraries on your system.

