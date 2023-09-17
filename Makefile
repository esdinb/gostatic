cli:
	go build -o gostatic -ldflags="-s -w" cli/main.go
.PHONY: cli

stripped: cli
	upx --brute gostatic
.PHONY: stripped
