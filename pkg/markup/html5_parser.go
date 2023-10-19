package markup

/*
#cgo pkg-config: libxml-2.0
#include <libxml/tree.h>
#include "html5_parser.h"

*/
import "C"
import (
	"os"

	"github.com/CannibalVox/cgoalloc"
)

type HTML5Parser struct {
	Context   *C.html5_parser_context_t
	Document  C.xmlDocPtr
	Allocator *cgoalloc.ArenaAllocator
}

func CreateHTML5Parser(doc *Document, node *Node) *HTML5Parser {
	var ptrn C.xmlNodePtr = nil
	if node != nil {
		ptrn = node.Ptr
	}
	if ptrc := C.html5_create_parser_context(doc.Ptr, ptrn); ptrc != nil {
		dfa := &cgoalloc.DefaultAllocator{}
		return &HTML5Parser{
			Context:   ptrc,
			Document:  doc.Ptr,
			Allocator: cgoalloc.CreateArenaAllocator(dfa),
		}
	}
	return nil
}

func (p *HTML5Parser) ParseChunk(chunk string) int {
	ptrc := (*C.char)(cgoalloc.CString(p.Allocator, chunk))
	if res := int(C.html5_parse_chunk(p.Context, ptrc, C.size_t(len(chunk)))); res < 0 {
		return -1
	}
	return len(chunk)
}

func (p *HTML5Parser) Terminate() int {
	if res := int(C.html5_parse_end_document(p.Context)); res < 0 {
		return -1
	}
	return 0
}

func (p *HTML5Parser) Free() {
	C.html5_destroy_parser_context(p.Context)
	p.Allocator.Destroy()
}

func (p *HTML5Parser) MyDoc() *Document {
	if ptrd := p.Document; ptrd != nil {
		return makeDoc(ptrd)
	}
	return nil
}

func ReadHTMLFile(srcPath string, options ParserOption) *Document {
	parser := CreateHTML5Parser(
		NewDoc("1.0"),
		nil,
	)
	if parser == nil {
		return nil
	}
	defer parser.Free()

	var chunk string
	if data, err := os.ReadFile(srcPath); err != nil {
		return nil
	} else {
		chunk = string(data)
	}
	if res := parser.ParseChunk(chunk); res < 0 {
		return nil
	}
	if res := parser.Terminate(); res < 0 {
		return nil
	}

	doc := parser.MyDoc()
	_ = doc.ProcessXInclude(options)

	return doc
}
