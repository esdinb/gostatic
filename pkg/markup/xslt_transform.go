package markup

/*
#cgo pkg-config: libxslt
#include <libxml/tree.h>
#include <libxslt/transform.h>

*/
import "C"
import "unsafe"

const (
	XSLT_DEFAULT_VERSION string = C.XSLT_DEFAULT_VERSION
	XSLT_DEFAULT_URL            = C.XSLT_DEFAULT_URL
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

func ApplyStylesheet(style *Stylesheet, doc *Document, params []string) *Document {
	cparams := C.malloc(C.size_t(len(params)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	a := (*[1<<30 - 1]*C.char)(cparams)
	for idx, param := range params {
		a[idx] = C.CString(param)
	}
	if ptr := C.xsltApplyStylesheet(style.Ptr, doc.Ptr, (**C.char)(cparams)); ptr == nil {
		return nil
	} else {
		return makeDoc(ptr)
	}
}
