package markup

/*
#include <libxml/parser.h>
#include <libxslt/documents.h>
#include "xslt_loader.h"

*/
import "C"
import "unsafe"

type DocLoaderContext struct {
	Ptr unsafe.Pointer // the context, either a stylesheet or a transformation context
}

type DocLoaderFunc func(string, *Dict, ParserOption, *DocLoaderContext, LoadType) *Document

type LoadType int

const (
	LoadStart      LoadType = C.XSLT_LOAD_START
	LoadStylesheet          = C.XSLT_LOAD_STYLESHEET
	LoadDocument            = C.XSLT_LOAD_DOCUMENT
)

func makeDocLoaderContext(ptr unsafe.Pointer) *DocLoaderContext {
	return &DocLoaderContext{ptr}
}

//export go_loader_callback
func go_loader_callback(uri *C.xmlChar, dict C.xmlDictPtr, options C.int, ctxt unsafe.Pointer, loadType C.xsltLoadType) C.xmlDocPtr {
	doc := loader(C.GoString((*C.char)(unsafe.Pointer(uri))), makeDict(dict), ParserOption(options), makeDocLoaderContext(ctxt), LoadType(loadType))
	if doc == nil {
		return nil
	}
	return doc.Ptr
}

func DefaultLoader(uri string, dict *Dict, options ParserOption, ctxt *DocLoaderContext, loadType LoadType) *Document {
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))
	var cdict C.xmlDictPtr
	if dict == nil {
		cdict = nil
	} else {
		cdict = dict.Ptr
	}
	var cctxt unsafe.Pointer
	if ctxt == nil {
		cctxt = nil
	} else {
		cctxt = ctxt.Ptr
	}
	ptr := C.default_loader((*C.xmlChar)(unsafe.Pointer(curi)), cdict, C.int(options), cctxt, (C.xsltLoadType)(loadType))
	return makeDoc(ptr)
}

var loader DocLoaderFunc

func SetLoaderFunc(f DocLoaderFunc) {
	if f == nil {
		C.xsltSetLoaderFunc(nil)
	} else {
		loader = f
		C.xsltSetLoaderFunc((C.xsltDocLoaderFunc)(C.custom_loader))
	}
}

func init() {
	C.save_default_loader()
}
