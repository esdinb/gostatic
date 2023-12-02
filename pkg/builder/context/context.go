package builder_context

import (
	"context"
	"fmt"
	"log"
	"os"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return fmt.Sprintf("context key %s\n", k.name)
}

var LoggerContextKey = contextKey{"logger"}
var LoggerHandleContextKey = contextKey{"loggerhandle"}
var BuildPathContextKey = contextKey{"buildpath"}
var InPathContextKey = contextKey{"inpath"}
var OutPathContextKey = contextKey{"outpath"}
var RootPathContextKey = contextKey{"rootpath"}
var InFileContextKey = contextKey{"infilepath"}
var OutFileContextKey = contextKey{"outfilepath"}
var DocumentContextKey = contextKey{"documentpath"}
var ParamsContextKey = contextKey{"paramspath"}
var StringParamsContextKey = contextKey{"strparamspath"}
var FormatterContextKey = contextKey{"formatterpath"}

func NewBuildContext() context.Context {
	logger := log.New(os.Stderr, "🐙 ", 0)
	ctx := context.WithValue(context.Background(), LoggerContextKey, logger)
	return ctx
}
