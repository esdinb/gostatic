package markup

/*
#include <libxml/xinclude.h>
*/
import "C"

func ProcessXInclude(doc *Document, flags ParserOption) int {
	return int(C.xmlXIncludeProcessFlags(doc.Ptr, C.int(flags)))
}
