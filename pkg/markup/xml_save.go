package markup

/*
#include <stdlib.h>
#include <libxml/xmlsave.h>
#include "xml_save.h"
*/
import "C"
import (
	"errors"
	"io"
	"os"
	"runtime/cgo"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES/STRUCTS
////////////////////////////////////////////////////////////////////////////////

type OutputWriteCallback = func(ctx interface{}, buffer *C.char, len C.int) C.int
type OutputCloseCallback = func(ctx interface{}) C.int

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

//export ioWrite
func ioWrite(ctx unsafe.Pointer, buffer *C.char, length C.int) C.int {
	handle := *(*cgo.Handle)(ctx)
	bytes := C.GoBytes(unsafe.Pointer(buffer), length)
	nn := 0
	if writer, ok := handle.Value().(io.Writer); !ok {
		panic("unable to cast writer context")
	} else {
		if n, err := writer.Write(bytes); err != nil {
			panic(err)
		} else {
			nn = n
		}
	}

	return C.int(nn)
}

//export ioClose
func ioClose(ctx unsafe.Pointer) C.int {
	handle := *(*cgo.Handle)(ctx)
	defer handle.Delete()
	file, ok := handle.Value().(*os.File)
	if ok {
		file.Close()
	}

	return C.int(0)
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE
////////////////////////////////////////////////////////////////////////////////

// xmlSaveClose
func (s *SaveContext) Free() int {
	return int(C.xmlSaveClose(s.Ptr))
}

// xmlSaveDoc
func (s *SaveContext) SaveDoc(doc *Document) error {
	// TODO: The function is not fully implemented yet as it does not return the byte count but 0 instead
	if length := int(C.xmlSaveDoc(s.Ptr, doc.Ptr)); length == -1 {
		return errors.New("failed to write document")
	}
	return nil
}

// xmlSaveFlush
func (s *SaveContext) SaveFlush() int {
	return int(C.xmlSaveFlush(s.Ptr))
}

// xmlSaveToFilename
func SaveToFilename(filename string, encoding string, options SaveOption) *SaveContext {
	ptrf := C.CString(filename)
	defer C.free(unsafe.Pointer(ptrf))
	ptre := C.CString(encoding)
	defer C.free(unsafe.Pointer(ptre))
	ptr := C.xmlSaveToFilename(ptrf, ptre, C.int(options))
	return makeSaveContext(ptr)
}

// xmlSaveToIO
func SaveToIO(ioCtx *os.File, encoding string, options SaveOption) *SaveContext {
	ptre := C.CString(encoding)
	defer C.free(unsafe.Pointer(ptre))
	handle := cgo.NewHandle(ioCtx)
	ptr := C.saveToIO(unsafe.Pointer(&handle), ptre, C.int(options))
	return makeSaveContext(ptr)
}

func SaveToIOCloser(ioCtx *os.File, encoding string, options SaveOption) *SaveContext {
	ptre := C.CString(encoding)
	defer C.free(unsafe.Pointer(ptre))
	handle := cgo.NewHandle(ioCtx)
	ptr := C.saveToIOCloser(unsafe.Pointer(&handle), ptre, C.int(options))
	return makeSaveContext(ptr)
}
