package config

import "fmt"

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return fmt.Sprintf("context key %s\n", k.name)
}

var LoggerContextKey = contextKey{"logger"}
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
