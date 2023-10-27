package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			fmt.Println(setting.Key, setting.Value)
		}
	}
}
