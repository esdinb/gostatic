package transformer

import (
	"context"
)

type Status int

const (
	Stop Status = iota + 1
	Continue
)

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
