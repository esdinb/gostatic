# gostatic
Static website generator.

# WIP

Based on [libxml2](https://gitlab.gnome.org/GNOME/libxml2) and [libxslt](https://gitlab.gnome.org/GNOME/libxslt).

Uses a fork of [jbussdiekers golibxml](https://github.com/jbussdieker/golibxml). Import is aliased to `../golibxml`.

Uses [yuins goldmark](https://github.com/yuin/goldmark) Markdown parser.

## Building

Clone the gostatic repository:

`git clone https://github.com/esdinb/gostatic.git`

Clone the golibxml fork:

`git clone https://github.com/esdinb/golibxml.git`

Change to gostatic directory and run `make`:

`cd gostatic && make`

