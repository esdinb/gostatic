package transformer

import (
	"gostatic/pkg/markup"
)

type Status int

const (
	Stop Status = iota + 1
	Continue
)

type Context struct {
	InPath   string
	OutPath  string
	RootPath string
	Document *markup.Document
}

type TransformerFunc func(*Context, []string) (*Context, Status, error)

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
