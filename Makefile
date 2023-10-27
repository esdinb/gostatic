CGO_CFLAGS ?= -I/usr/local/opt/libxml2/include -I/usr/local/opt/libxslt/include -I/usr/local/include
CGO_LDFLAGS ?= -L/usr/local/opt/libxml2/lib -lxml2 -L/usr/local/opt/libxslt/lib -lxslt -lexslt -L/usr/local/lib -llexbor_static
BUILD ?= $(git describe)

cli:
	git describe > cli/cmd/version.txt
	CGO_ENABLED=1 CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" go build -o gostatic -trimpath cli/main.go
.PHONY: cli

static:
	git describe > cli/cmd/version.txt
	CGO_ENABLED=1 CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" go build -o gostatic -trimpath -ldflags="-linkmode=external -extldflags=-static -w -s" cli/main.go
.PHONY: static

install: cli
	cp ./gostatic ~/go/bin
.PHONY: install

compressed: static
	upx --brute gostatic
.PHONY: compressed
