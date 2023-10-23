package markup

/*
#include <stdlib.h>
#include <libxml/xpath.h>

xmlNode* fetchNode(xmlNodeSet *nodeset, int index) {
  	return nodeset->nodeTab[index];
}

*/
import "C"
import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES/STRUCTS
////////////////////////////////////////////////////////////////////////////////

type XPath struct {
	Ptr C.xmlXPathCompExprPtr
}

type XPathContext struct {
	Ptr C.xmlXPathContextPtr
}

type XPathObject struct {
	Ptr C.xmlXPathObjectPtr
}

type XPathFunction = func(C.xmlXPathParserContextPtr, C.int)

type NodeSet struct {
	Ptr C.xmlNodeSetPtr
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS/ENUM
////////////////////////////////////////////////////////////////////////////////

type XpathObjectType int

const (
	XPATH_UNDEFINED   XpathObjectType = 0
	XPATH_NODESET                     = 1
	XPATH_BOOLEAN                     = 2
	XPATH_NUMBER                      = 3
	XPATH_STRING                      = 4
	XPATH_POINT                       = 5
	XPATH_RANGE                       = 6
	XPATH_LOCATIONSET                 = 7
	XPATH_USERS                       = 8
	XPATH_XSLT_TREE                   = 9 //: An XSLT value tree, non modifiable
)

// //////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS
// //////////////////////////////////////////////////////////////////////////////
func makeXpath(ptr C.xmlXPathCompExprPtr) *XPath {
	if ptr == nil {
		return nil
	}
	return &XPath{ptr}
}

func makeXpathObj(ptr C.xmlXPathObjectPtr) *XPathObject {
	if ptr == nil {
		return nil
	}
	return &XPathObject{ptr}
}

func makeNodeSet(ptr C.xmlNodeSetPtr) *NodeSet {
	if ptr == nil {
		return nil
	}
	return &NodeSet{ptr}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE
////////////////////////////////////////////////////////////////////////////////

// xmlXPathCastToNumber
func (obj *XPathObject) CastToNumber() float32 {
	cdbl := C.xmlXPathCastToNumber(obj.Ptr)
	return float32(cdbl)
}

// xmlXPathCastToString
func (obj *XPathObject) String() string {
	cstr := C.xmlXPathCastToString(obj.Ptr)
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString((*C.char)(unsafe.Pointer(cstr)))
}

// xmlXPathCompile
func CompileXPath(str string) *XPath {
	ptr := C.CString(str)
	defer C.free(unsafe.Pointer(ptr))
	return makeXpath(C.xmlXPathCompile((*C.xmlChar)(unsafe.Pointer(ptr))))
}

// xmlXPathConvertBoolean
func (obj *XPathObject) ConvertBoolean() *XPathObject {
	return makeXpathObj(C.xmlXPathConvertBoolean(obj.Ptr))
}

// xmlXPathConvertNumber
func (obj *XPathObject) ConvertNumber() *XPathObject {
	return makeXpathObj(C.xmlXPathConvertNumber(obj.Ptr))
}

// xmlXPathConvertString
func (obj *XPathObject) ConvertString() *XPathObject {
	return makeXpathObj(C.xmlXPathConvertString(obj.Ptr))
}

func (obj *XPathObject) ResultsChannel() chan *Node {
	channel := make(chan *Node)
	go func(obj *XPathObject, channel chan *Node) {
		if obj.Ptr._type != 1 {
			close(channel)
			return
		}
		for i := 0; i < int(obj.Ptr.nodesetval.nodeNr); i++ {
			channel <- makeNode(C.fetchNode(obj.Ptr.nodesetval, C.int(i)))
		}
		close(channel)
	}(obj, channel)
	return channel
}

func (obj *XPathObject) Results() []*Node {
	var (
		results []*Node
		length  int
	)
	if obj.Ptr._type != 1 {
		return []*Node{}
	}
	length = int(obj.Ptr.nodesetval.nodeNr)
	results = make([]*Node, length)
	for i := 0; i < length; i++ {
		results[i] = makeNode(C.fetchNode(obj.Ptr.nodesetval, C.int(i)))
	}
	return results
}

// xmlXPathCtxtCompile
func (ctx *XPathContext) Compile(str string) *XPath {
	ptr := C.CString(str)
	defer C.free(unsafe.Pointer(ptr))
	return makeXpath(C.xmlXPathCtxtCompile(ctx.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr))))
}

// xmlXPathEval
func (ctx *XPathContext) Eval(str string) *XPathObject {
	ptr := C.CString(str)
	defer C.free(unsafe.Pointer(ptr))
	return makeXpathObj(C.xmlXPathEval((*C.xmlChar)(unsafe.Pointer(ptr)), ctx.Ptr))
}

// xmlXPathEvalPredicate
func (ctx *XPathContext) EvalPredicate(res *XPathObject) bool {
	result := C.xmlXPathEvalPredicate(ctx.Ptr, res.Ptr)
	return int(result) != 0
}

// xmlXPathEvalExpression
func (ctx *XPathContext) EvalExpression(str string) *XPathObject {
	ptr := C.CString(str)
	defer C.free(unsafe.Pointer(ptr))
	return makeXpathObj(C.xmlXPathEvalExpression((*C.xmlChar)(unsafe.Pointer(ptr)), ctx.Ptr))
}

// xmlXPathFreeCompExpr
func (xpath *XPath) Free() {
	C.xmlXPathFreeCompExpr(xpath.Ptr)
	xpath.Ptr = nil
}

// xmlXPathFreeContext
func (ctx *XPathContext) Free() {
	C.xmlXPathFreeContext(ctx.Ptr)
	ctx.Ptr = nil
}

// xmlXPathFreeNodeSet
func (nodeset *NodeSet) Free() {
	C.xmlXPathFreeNodeSet(nodeset.Ptr)
	nodeset.Ptr = nil
}

// xmlXPathFreeNodeSetList
func (obj *XPathObject) FreeList() {
	C.xmlXPathFreeNodeSetList(obj.Ptr)
	obj.Ptr = nil
}

// xmlXPathFreeObject
func (obj *XPathObject) Free() {
	C.xmlXPathFreeObject(obj.Ptr)
	obj.Ptr = nil
}

// xmlXPathIsInf
func xmlXPathIsInf(val float32) int {
	return int(C.xmlXPathIsInf(C.double(val)))
}

// xmlXPathIsNaN
func XpathIsNaN(val float32) int {
	return int(C.xmlXPathIsNaN(C.double(val)))
}

// xmlXPathNewContext
func NewXPathContext(doc *Document) *XPathContext {
	return &XPathContext{C.xmlXPathNewContext(doc.Ptr)}
}

// xmlXPathNodeSetCreate
func NodeSetCreate(node *Node) *NodeSet {
	return makeNodeSet(C.xmlXPathNodeSetCreate(node.Ptr))
}

// xmlXPathObjectCopy
func (obj *XPathObject) Copy() *XPathObject {
	return makeXpathObj(C.xmlXPathObjectCopy(obj.Ptr))
}

// xmlXPathOrderDocElems
func (doc *Document) OrderDocElems() int {
	return int(C.xmlXPathOrderDocElems(doc.Ptr))
}
