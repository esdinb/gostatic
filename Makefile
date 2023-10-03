CGO_CFLAGS ?= -I/usr/local/include
CGO_LDFLAGS ?= -L/usr/local/lib -llexbor_static

cli:
	CGO_ENABLED=1 CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" go build -o gostatic -trimpath -ldflags="-w -s" cli/main.go
.PHONY: cli

static:
	CGO_ENABLED=1 CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" go build -o gostatic -trimpath -ldflags="-linkmode=external -extldflags=-static -w -s" cli/main.go
.PHONY: static

install: cli
	cp ./gostatic ~/go/bin
.PHONY: install

compressed: static
	upx --brute gostatic
.PHONY: compressed
