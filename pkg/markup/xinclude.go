package markup

/*
#cgo pkg-config: libxml-2.0
#include <libxml/xinclude.h>
*/
import "C"

func ProcessXInclude(doc *Document, flags ParserOption) int {
	return int(C.xmlXIncludeProcessFlags(doc.Ptr, C.int(flags)))
}
