# gostatic
Static website generator.

# WIP

Based on [libxml2](https://gitlab.gnome.org/GNOME/libxml2) and [libxslt](https://gitlab.gnome.org/GNOME/libxslt).

CGO wrappers for libxml2 are based on [jbussdiekers golibxml](https://github.com/jbussdieker/golibxml) code.

Uses [yuins goldmark](https://github.com/yuin/goldmark) Markdown parser.

## Requirements

The libxml2 and libxslt libraries needs to be installed on the system.

The [lexbor](https://github.com/lexbor/lexbor) library needs to be installed on the system.

C library bindings uses CGO and require a C compiler to also be installed.

## Building

Clone the gostatic repository:

`git clone https://github.com/esdinb/gostatic.git`

In the gostatic directory run `go build` something like this:

`CGO_CFLAGS="-I/usr/local/include" CGO_LDFLAGS="-L/usr/local/lib -llexbor_static" go build -o gostatic cli/main.go`

Change `CGO_CFLAGS` and `CGO_LDFLAGS` to match the paths to the lexbor libraries.


