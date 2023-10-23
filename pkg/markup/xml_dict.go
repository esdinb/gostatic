package markup

/*
#include <libxml/xmlstring.h>
#include <libxml/dict.h>
*/
import "C"

type Dict struct {
	Ptr C.xmlDictPtr
}

func makeDict(dict C.xmlDictPtr) *Dict {
	if dict == nil {
		return nil
	}
	return &Dict{dict}
}

// xmlDictCreate
func CreateDict() *Dict {
	cdict := C.xmlDictCreate()
	return makeDict(cdict)
}

// xmlDictSetLimit
func (d *Dict) SetLimit(limit int) int {
	return int(C.xmlDictSetLimit(d.Ptr, C.size_t(limit)))
}

// xmlDictGetUsage
func (d *Dict) GetUsage() int {
	return int(C.xmlDictGetUsage(d.Ptr))
}

// xmlDictCreateSub
func (d *Dict) CreateSub() *Dict {
	cdict := C.xmlDictCreateSub(d.Ptr)
	return makeDict(cdict)
}

// xmlDictReference
func (d *Dict) Reference() int {
	return int(C.xmlDictReference(d.Ptr))
}

// xmlDictFree
func (d *Dict) Free() {
	C.xmlDictFree(d.Ptr)
}
