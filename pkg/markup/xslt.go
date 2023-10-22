package markup

/*
#include <libxslt/xslt.h>
#include <libxslt/xsltInternals.h>

static inline xmlChar *to_xmlcharptr(const char *s) { return (xmlChar *)s; }

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
	ptrs := C.xsltParseStylesheetFile(C.to_xmlcharptr(ptrf))
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
