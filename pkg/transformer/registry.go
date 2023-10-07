package transformer

import (
	"context"
)

type Status int

const (
	Stop Status = iota + 1
	Continue
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "gostatic transform context key " + k.name }

var InPathContextKey = contextKey{"inpath"}
var OutPathContextKey = contextKey{"outpath"}
var RootPathContextKey = contextKey{"rootpath"}
var DocumentContextKey = contextKey{"documentpath"}
var ParamsContextKey = contextKey{"paramspath"}
var StringParamsContextKey = contextKey{"strparamspath"}

type TransformerFunc func(context.Context, []string) (context.Context, Status, error)

type registry map[string]TransformerFunc

func (r *registry) Register(name string, transformer TransformerFunc) {
	(*r)[name] = transformer
}

func (r *registry) Lookup(name string) TransformerFunc {
	if transformer, ok := (*r)[name]; ok {
		return transformer
	} else {
		return nil
	}
}

var Registry = registry{}
