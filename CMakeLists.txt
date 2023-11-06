cmake_minimum_required(VERSION 3.15)

project(gostatic
    VERSION 0.0.1
    DESCRIPTION "gostatic")

add_custom_target(cli ALL
    COMMAND env CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS="${CGO_LDFLAGS}" go build -o gostatic -trimpath cli/main.go
)

