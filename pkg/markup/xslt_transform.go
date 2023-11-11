package markup

/*
#include <libxml/tree.h>
#include <libxml/parser.h>
#include <libxslt/transform.h>
#include <libxslt/variables.h>
#include <libxslt/xsltutils.h>
#include <libxslt/extensions.h>
#include "xslt_transform.h"
#include "xslt_extensions.h"
#include "xml_error.h"

*/
import "C"
import (
	"log"
	"runtime/cgo"
	"unsafe"
)

const (
	XSLT_DEFAULT_VERSION string = C.XSLT_DEFAULT_VERSION
	XSLT_DEFAULT_URL            = C.XSLT_DEFAULT_URL
	XSLT_PARSE_OPTIONS          = C.XML_PARSE_NOENT | C.XML_PARSE_NOCDATA
)

type TransformContext struct {
	Ptr    C.xsltTransformContextPtr
	Logger *log.Logger
}

func NewTransformContext(style *Stylesheet, doc *Document, logger *log.Logger) *TransformContext {
	if ptr := C.xsltNewTransformContext(style.Ptr, doc.Ptr); ptr != nil {
		C.registerExtensionFunctions(ptr)
		C.xsltSetCtxtParseOptions(ptr, XSLT_PARSE_OPTIONS)
		handle := cgo.NewHandle(logger)
		C.set_xslt_transform_error_func(ptr, unsafe.Pointer(&handle))
		return &TransformContext{ptr, logger}
	}
	return nil
}

// https://mail.gnome.org/archives/xslt/2009-December/msg00002.html
func (t *TransformContext) ApplyStylesheet(style *Stylesheet, doc *Document, params []string, strparams []string) *Document {

	cparams := C.makeParamsArray(C.int(len(params) + 1))
	defer C.freeParamsArray(cparams, C.int(len(params)+1))
	for idx, param := range params {
		C.setParamsElement(cparams, C.CString(param), C.int(idx))
	}

	cstrparams := C.makeParamsArray(C.int(len(strparams) + 1))
	defer C.freeParamsArray(cstrparams, C.int(len(strparams)+1))
	for idx, strparam := range strparams {
		C.setParamsElement(cstrparams, C.CString(strparam), C.int(idx))
	}

	if C.xsltQuoteUserParams(t.Ptr, cstrparams) != -1 {
		if ptr := C.xsltApplyStylesheetUser(style.Ptr, doc.Ptr, cparams, nil, nil, t.Ptr); ptr != nil {
			return makeDoc(ptr)
		}
	}

	return nil
}

func (t *TransformContext) Free() {
	C.xsltFreeTransformContext(t.Ptr)
}

func ApplyStylesheet(style *Stylesheet, doc *Document) *Document {
	if ptr := C.xsltApplyStylesheet(style.Ptr, doc.Ptr, nil); ptr != nil {
		return makeDoc(ptr)
	}

	return nil
}
