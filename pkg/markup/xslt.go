package markup

/*
#include <stdlib.h>
#include <libxslt/xslt.h>
#include <libxslt/xsltInternals.h>
*/
import "C"
import "unsafe"

type Stylesheet struct {
	Ptr C.xsltStylesheetPtr
}

func (s *Stylesheet) Free() {
	C.xsltFreeStylesheet(s.Ptr)
}

func ParseStylesheetFile(filename string) *Stylesheet {
	ptrf := C.CString(filename)
	defer C.free(unsafe.Pointer(ptrf))
	ptrs := C.xsltParseStylesheetFile((*C.xmlChar)(unsafe.Pointer(ptrf)))
	return &Stylesheet{ptrs}
}

func ParseStylesheetDoc(doc *Document) *Stylesheet {
	if ptr := C.xsltParseStylesheetDoc(doc.Ptr); ptr == nil {
		return nil
	} else {
		return &Stylesheet{ptr}
	}
}

func LoadStylesheetPI(doc *Document) *Stylesheet {
	if ptr := C.xsltLoadStylesheetPI(doc.Ptr); ptr == nil {
		return nil
	} else {
		return &Stylesheet{ptr}
	}
}
