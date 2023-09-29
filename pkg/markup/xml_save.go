package markup

/*
#cgo pkg-config: libxml-2.0
#include <libxml/xmlsave.h>

static inline void free_string(char* s) { free(s); }
static inline void free_xmlstring(xmlChar* s) { free(s); }
static inline xmlChar *to_xmlcharptr(const char *s) { return (xmlChar *)s; }
static inline char *to_charptr(const xmlChar *s) { return (char *)s; }
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES/STRUCTS
////////////////////////////////////////////////////////////////////////////////

type SaveContext struct {
	Ptr C.xmlSaveCtxtPtr
}

type SaveOption int

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS/ENUM
////////////////////////////////////////////////////////////////////////////////

const (
	XML_SAVE_FORMAT   SaveOption = C.XML_SAVE_FORMAT   /* format save output */
	XML_SAVE_NO_DECL             = C.XML_SAVE_NO_DECL  /* drop the xml declaration */
	XML_SAVE_NO_EMPTY            = C.XML_SAVE_NO_EMPTY /* no empty tags */
	XML_SAVE_NO_XHTML            = C.XML_SAVE_NO_XHTML /* disable XHTML1 specific rules */
	XML_SAVE_XHTML               = C.XML_SAVE_XHTML    /* force XHTML1 specific rules */
	XML_SAVE_AS_XML              = C.XML_SAVE_AS_XML   /* force XML serialization on HTML doc */
	XML_SAVE_AS_HTML             = C.XML_SAVE_AS_HTML  /* force HTML serialization on XML doc */
	XML_SAVE_WSNONSIG            = C.XML_SAVE_WSNONSIG /*  format with non-significant whitespace */
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func makeSaveContext(ctxt C.xmlSaveCtxtPtr) *SaveContext {
	if ctxt == nil {
		return nil
	}
	return &SaveContext{ctxt}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE
////////////////////////////////////////////////////////////////////////////////

// xmlSaveClose
func (s *SaveContext) Free() int {
	return int(C.xmlSaveClose(s.Ptr))
}

// xmlSaveDoc
func (s *SaveContext) SaveDoc(doc *Document) int {
	// TODO: The function is not fully implemented yet as it does not return the byte count but 0 instead
	return int(C.xmlSaveDoc(s.Ptr, doc.Ptr))
}

// xmlSaveFlush
func (s *SaveContext) SaveFlush() int {
	return int(C.xmlSaveFlush(s.Ptr))
}

// xmlSaveToFilename
func SaveToFilename(filename string, encoding string, options SaveOption) *SaveContext {
	ptrf := C.CString(filename)
	defer C.free_string(ptrf)
	ptre := C.CString(encoding)
	defer C.free_string(ptre)
	ptr := C.xmlSaveToFilename(ptrf, ptre, C.int(options))
	return makeSaveContext(ptr)
}
