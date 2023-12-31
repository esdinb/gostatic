package markup

/*
#include <stdlib.h>
#include <libxml/parser.h>
*/
import "C"
import "unsafe"

////////////////////////////////////////////////////////////////////////////////
// TYPES/STRUCTS
////////////////////////////////////////////////////////////////////////////////

type ParserOption int

const (
	XML_PARSE_RECOVER    ParserOption = C.XML_PARSE_RECOVER    //: recover on errors
	XML_PARSE_NOENT                   = C.XML_PARSE_NOENT      //: substitute entities
	XML_PARSE_DTDLOAD                 = C.XML_PARSE_DTDLOAD    //: load the external subset
	XML_PARSE_DTDATTR                 = C.XML_PARSE_DTDATTR    //: default DTD attributes
	XML_PARSE_DTDVALID                = C.XML_PARSE_DTDVALID   //: validate with the DTD
	XML_PARSE_NOERROR                 = C.XML_PARSE_NOERROR    //: suppress error reports
	XML_PARSE_NOWARNING               = C.XML_PARSE_NOWARNING  //: suppress warning reports
	XML_PARSE_PEDANTIC                = C.XML_PARSE_PEDANTIC   //: pedantic error reporting
	XML_PARSE_NOBLANKS                = C.XML_PARSE_NOBLANKS   //: remove blank nodes
	XML_PARSE_SAX1                    = C.XML_PARSE_SAX1       //: use the SAX1 interface internally
	XML_PARSE_XINCLUDE                = C.XML_PARSE_XINCLUDE   //: Implement XInclude substitition
	XML_PARSE_NONET                   = C.XML_PARSE_NONET      //: Forbid network access
	XML_PARSE_NODICT                  = C.XML_PARSE_NODICT     //: Do not reuse the context dictionnary
	XML_PARSE_NSCLEAN                 = C.XML_PARSE_NSCLEAN    //: remove redundant namespaces declarations
	XML_PARSE_NOCDATA                 = C.XML_PARSE_NOCDATA    //: merge CDATA as text nodes
	XML_PARSE_NOXINCNODE              = C.XML_PARSE_NOXINCNODE //: do not generate XINCLUDE START/END nodes
	XML_PARSE_COMPACT                 = C.XML_PARSE_COMPACT    //: compact small text nodes; no modification of the tree allowed afterwards (will possibly crash if you try to modify the tree)
	XML_PARSE_OLD10                   = C.XML_PARSE_OLD10      //: parse using XML-1.0 before update 5
	XML_PARSE_NOBASEFIX               = C.XML_PARSE_NOBASEFIX  //: do not fixup XINCLUDE xml:base uris
	XML_PARSE_HUGE                    = C.XML_PARSE_HUGE       //: relax any hardcoded limit from the parser
	XML_PARSE_OLDSAX                  = C.XML_PARSE_OLDSAX     //: parse using SAX2 interface from before 2.7.0
)

type Parser struct {
	Ptr C.xmlParserCtxtPtr
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func makeParser(parser C.xmlParserCtxtPtr) *Parser {
	if parser == nil {
		return nil
	}
	return &Parser{parser}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE
////////////////////////////////////////////////////////////////////////////////

// xmlByteConsumed
func (p *Parser) ByteConsumed() int {
	return int(C.xmlByteConsumed(p.Ptr))
}

// xmlCleanupParser
func CleanupParser() {
	C.xmlCleanupParser()
}

// xmlClearParserCtxt
func (p *Parser) Clear() {
	C.xmlClearParserCtxt(p.Ptr)
}

// xmlCreateDocParserCtxt
func CreateDocParser(cur string) *Parser {
	ptr := C.CString(cur)
	defer C.free(unsafe.Pointer(ptr))
	cparser := C.xmlCreateDocParserCtxt((*C.xmlChar)(unsafe.Pointer(ptr)))
	return makeParser(cparser)
}

// xmlCreatePushParserCtxt
func CreatePushParser(chunk string, filename string) *Parser {
	ptrc := C.CString(chunk)
	defer C.free(unsafe.Pointer(ptrc))
	ptrf := C.CString(filename)
	defer C.free(unsafe.Pointer(ptrf))
	// first two arguments is a SAX handler and a pointer to user data
	cparser := C.xmlCreatePushParserCtxt(nil, nil, (*C.char)(ptrc), C.int(len(chunk)), (*C.char)(ptrf))
	return makeParser(cparser)
}

// xmlCtxtReadDoc
func (p *Parser) ReadDoc(input string, url string, encoding string, options ParserOption) *Document {
	ptri := C.CString(input)
	defer C.free(unsafe.Pointer(ptri))
	ptru := C.CString(url)
	defer C.free(unsafe.Pointer(ptru))
	ptre := C.CString(encoding)
	defer C.free(unsafe.Pointer(ptre))
	doc := C.xmlCtxtReadDoc(p.Ptr, (*C.xmlChar)(unsafe.Pointer(ptri)), ptru, ptre, C.int(options))
	return makeDoc(doc)
}

// xmlCtxtReset
func (p *Parser) Reset() {
	C.xmlCtxtReset(p.Ptr)
}

// xmlCtxtUseOptions
func (p *Parser) UseOptions(options ParserOption) int {
	return int(C.xmlCtxtUseOptions(p.Ptr, C.int(options)))
}

// xmlFreeParserCtxt
func (p *Parser) Free() {
	C.xmlFreeParserCtxt(p.Ptr)
}

// xmlNewParserCtxt
func NewParser() *Parser {
	pctx := C.xmlNewParserCtxt()
	return makeParser(pctx)
}

// xmlParseDTD
func ParseDTD(ExternalID string, SystemID string) *Dtd {
	ptre := C.CString(ExternalID)
	defer C.free(unsafe.Pointer(ptre))
	ptrs := C.CString(SystemID)
	defer C.free(unsafe.Pointer(ptrs))
	cdtd := C.xmlParseDTD((*C.xmlChar)(unsafe.Pointer(ptre)), (*C.xmlChar)(unsafe.Pointer(ptrs)))
	return makeDtd(cdtd)
}

// xmlParseDocument
func (p *Parser) Parse() int {
	return int(C.xmlParseDocument(p.Ptr))
}

// xmlParseChunk
func (p *Parser) ParseChunk(chunk string) int {
	ptr := C.CString(chunk)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlParseChunk(p.Ptr, (*C.char)(ptr), C.int(len(chunk)), C.int(0)))
}

func (p *Parser) Terminate(chunk string) int {
	ptr := C.CString(chunk)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlParseChunk(p.Ptr, (*C.char)(ptr), C.int(len(chunk)), C.int(1)))
}

func (p *Parser) MyDoc() *Document {
	if docptr := p.Ptr.myDoc; docptr != nil {
		return makeDoc(docptr)
	}
	return nil
}

// xmlReadDoc
func ReadDoc(input string, url string, encoding string, options ParserOption) *Document {
	ptri := C.CString(input)
	defer C.free(unsafe.Pointer(ptri))
	ptru := C.CString(url)
	defer C.free(unsafe.Pointer(ptru))
	ptre := C.CString(encoding)
	defer C.free(unsafe.Pointer(ptre))
	doc := C.xmlReadDoc((*C.xmlChar)(unsafe.Pointer(ptri)), ptru, ptre, C.int(options))
	return makeDoc(doc)
}

// xmlReadFile
func ReadFile(filename string, encoding string, options ParserOption) *Document {
	ptrf := C.CString(filename)
	defer C.free(unsafe.Pointer(ptrf))
	ptre := C.CString(encoding)
	defer C.free(unsafe.Pointer(ptre))
	doc := C.xmlReadFile(ptrf, ptre, C.int(options))
	return makeDoc(doc)
}

// xmlReadMemory
func ReadMemory(buffer []byte, url string, encoding string, options ParserOption) *Document {
	ptru := C.CString(url)
	defer C.free(unsafe.Pointer(ptru))
	ptre := C.CString(encoding)
	defer C.free(unsafe.Pointer(ptre))
	doc := C.xmlReadMemory((*C.char)(unsafe.Pointer(&buffer[0])), C.int(len(buffer)), ptru, ptre, C.int(options))
	return makeDoc(doc)
}

// xmlStopParser
func (p *Parser) Stop() {
	C.xmlStopParser(p.Ptr)
}

func init() {
	C.xmlInitParser()
}
