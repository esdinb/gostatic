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

*/
import "C"

const (
	XSLT_DEFAULT_VERSION string = C.XSLT_DEFAULT_VERSION
	XSLT_DEFAULT_URL            = C.XSLT_DEFAULT_URL
	XSLT_PARSE_OPTIONS          = C.XML_PARSE_NOENT | C.XML_PARSE_NOCDATA
)

func NewTransformContext(style *Stylesheet, doc *Document) *TransformContext {
	if ptr := C.xsltNewTransformContext(style.Ptr, doc.Ptr); ptr == nil {
		return nil
	} else {
		return &TransformContext{ptr}
	}
}

type TransformContext struct {
	Ptr C.xsltTransformContextPtr
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

func ApplyStylesheetUser(style *Stylesheet, doc *Document, params []string, strparams []string) *Document {

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

	// https://mail.gnome.org/archives/xslt/2009-December/msg00002.html
	if ctx := C.xsltNewTransformContext(style.Ptr, doc.Ptr); ctx != nil {
		defer C.xsltFreeTransformContext(ctx)
		C.registerExtensionFunctions(ctx)
		C.xsltSetCtxtParseOptions(ctx, XSLT_PARSE_OPTIONS)
		if C.xsltQuoteUserParams(ctx, cstrparams) != -1 {
			if ptr := C.xsltApplyStylesheetUser(style.Ptr, doc.Ptr, cparams, nil, nil, ctx); ptr != nil {
				return makeDoc(ptr)
			}
		}
	}

	return nil
}
