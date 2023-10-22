package markup

/*
#include <libxml/parser.h>
#include <libxml/xpath.h>
#include <libxml/xpathInternals.h>
#include <libxslt/extensions.h>
#include <libexslt/exslt.h>
#include "xslt_extensions.h"

*/
import "C"
import (
	"time"
	"unsafe"
)

//export FormatDateCallback
func FormatDateCallback(ctx C.xmlXPathParserContextPtr, nArgs C.int) {
	var (
		arg1 *C.xmlChar
		arg2 *C.xmlChar
		arg3 *C.xmlChar
	)
	if nArgs != 3 {
		filename := C.CString("missing __FILE__ macro")
		defer C.free(unsafe.Pointer(filename))
		C.xmlXPatherror(ctx, filename, 0, C.XPATH_INVALID_ARITY)
		if ctx != nil {
			ctx.error = C.XPATH_INVALID_ARITY
		}
		return
	}
	arg3 = C.xmlXPathPopString(ctx)
	arg2 = C.xmlXPathPopString(ctx)
	arg1 = C.xmlXPathPopString(ctx)
	date := C.GoString((*C.char)(unsafe.Pointer(arg1)))
	inLayout := C.GoString((*C.char)(unsafe.Pointer(arg2)))
	outLayout := C.GoString((*C.char)(unsafe.Pointer(arg3)))

	t, _ := time.Parse(inLayout, date)
	result := C.CString(t.Format(outLayout))

	C.valuePush(ctx, C.xmlXPathWrapString((*C.xmlChar)(unsafe.Pointer(result))))
}

func init() {
	C.exsltCommonRegister()
	C.exsltStrRegister()
}
