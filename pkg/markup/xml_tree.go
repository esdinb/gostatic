// Package markup is a wrapper around libxml2, libxslt and lexbor html5 parser
//
// libxml2 wrappers are based on and mostly incorporated from jbussdiekers golibxml
//
// https://github.com/jbussdieker/golibxml
package markup

/*
#include <stdlib.h>
#include <libxml/tree.h>
#include <libxml/xinclude.h>
*/
import "C"
import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES/STRUCTS
////////////////////////////////////////////////////////////////////////////////

type AllocationScheme int

type Dtd struct {
	Ptr C.xmlDtdPtr
}

type Attribute struct {
	Ptr C.xmlAttrPtr
}

type Node struct {
	Ptr C.xmlNodePtr
}

type TextNode struct {
	*Node
}

type Document struct {
	*Node
	Ptr C.xmlDocPtr
}

type Namespace struct {
	Ptr C.xmlNsPtr
}

type Buffer struct {
	Ptr C.xmlBufferPtr
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS/ENUM
////////////////////////////////////////////////////////////////////////////////

const (
	XML_BUFFER_ALLOC_DOUBLEIT  AllocationScheme = 1 //: double each time one need to grow
	XML_BUFFER_ALLOC_EXACT                      = 2 //: grow only to the minimal size
	XML_BUFFER_ALLOC_IMMUTABLE                  = 3 //: immutable buffer
	XML_BUFFER_ALLOC_IO                         = 4 //: special allocation scheme used for I/O
)

type ElementType int

const (
	XML_ELEMENT_NODE       ElementType = 1
	XML_ATTRIBUTE_NODE                 = 2
	XML_TEXT_NODE                      = 3
	XML_CDATA_SECTION_NODE             = 4
	XML_ENTITY_REF_NODE                = 5
	XML_ENTITY_NODE                    = 6
	XML_PI_NODE                        = 7
	XML_COMMENT_NODE                   = 8
	XML_DOCUMENT_NODE                  = 9
	XML_DOCUMENT_TYPE_NODE             = 10
	XML_DOCUMENT_FRAG_NODE             = 11
	XML_NOTATION_NODE                  = 12
	XML_HTML_DOCUMENT_NODE             = 13
	XML_DTD_NODE                       = 14
	XML_ELEMENT_DECL                   = 15
	XML_ATTRIBUTE_DECL                 = 16
	XML_ENTITY_DECL                    = 17
	XML_NAMESPACE_DECL                 = 18
	XML_XINCLUDE_START                 = 19
	XML_XINCLUDE_END                   = 20
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func makeDoc(doc C.xmlDocPtr) *Document {
	if doc == nil {
		return nil
	}
	return &Document{
		Ptr:  doc,
		Node: makeNode(C.xmlNodePtr(unsafe.Pointer(doc))),
	}
}

func makeDtd(dtd C.xmlDtdPtr) *Dtd {
	if dtd == nil {
		return nil
	}
	return &Dtd{dtd}
}

func makeNamespace(ns C.xmlNsPtr) *Namespace {
	if ns == nil {
		return nil
	}
	return &Namespace{ns}
}

func makeNode(node C.xmlNodePtr) *Node {
	if node == nil {
		return nil
	}
	return &Node{node}
}

func makeTextNode(node C.xmlNodePtr) *TextNode {
	if node == nil {
		return nil
	}
	return &TextNode{makeNode(node)}
}

func makeAttribute(attr C.xmlAttrPtr) *Attribute {
	if attr == nil {
		return nil
	}
	return &Attribute{attr}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE
////////////////////////////////////////////////////////////////////////////////

// xmlAddChild
func (parent *Node) AddChild(cur *Node) *Node {
	return makeNode(C.xmlAddChild(parent.Ptr, cur.Ptr))
}

// xmlAddChildList
func (parent *Node) AddChildList(cur Node) *Node {
	return makeNode(C.xmlAddNextSibling(parent.Ptr, cur.Ptr))
}

// xmlAddNextSibling
func (cur *Node) AddNextSibling(elem *Node) *Node {
	return makeNode(C.xmlAddNextSibling(cur.Ptr, elem.Ptr))
}

// xmlAddPrevSibling
func (cur *Node) AddPrevSibling(elem *Node) *Node {
	return makeNode(C.xmlAddPrevSibling(cur.Ptr, elem.Ptr))
}

// xmlAddSibling
func (cur *Node) AddSibling(elem Node) *Node {
	return makeNode(C.xmlAddSibling(cur.Ptr, elem.Ptr))
}

// xmlBufferAdd
func (buffer *Buffer) Add(str []byte) int {
	return int(C.xmlBufferAdd(buffer.Ptr, (*C.xmlChar)(unsafe.Pointer(&str[0])), C.int(len(str))))
}

// xmlBufferAddHead
func (buffer *Buffer) AddHead(str []byte) int {
	return int(C.xmlBufferAddHead(buffer.Ptr, (*C.xmlChar)(unsafe.Pointer(&str[0])), C.int(len(str))))
}

// xmlBufferCat/xmlBufferCCat
func (buffer *Buffer) Cat(str string) int {
	ptr := C.CString(str)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlBufferCCat(buffer.Ptr, ptr))
}

// xmlBufferContent
func (buffer *Buffer) Content() string {
	return C.GoString((*C.char)(unsafe.Pointer(C.xmlBufferContent(buffer.Ptr))))
}

// xmlBufferCreate
func NewBuffer() *Buffer {
	return &Buffer{C.xmlBufferCreate()}
}

// xmlBufferCreateSize
func NewBufferSize(size int) *Buffer {
	return &Buffer{C.xmlBufferCreateSize(C.size_t(size))}
}

// xmlBufferEmpty
func (buffer *Buffer) Empty() {
	C.xmlBufferEmpty(buffer.Ptr)
}

// xmlBufferFree
func (buffer *Buffer) Free() {
	C.xmlBufferFree(buffer.Ptr)
	buffer.Ptr = nil
}

// xmlBufferGrow
func (buffer *Buffer) Grow(length int) int {
	return int(C.xmlBufferGrow(buffer.Ptr, C.uint(length)))
}

// xmlBufferLength
func (buffer *Buffer) Length() int {
	return int(C.xmlBufferLength(buffer.Ptr))
}

// xmlBufferResize
func (buffer *Buffer) Resize(size int) int {
	return int(C.xmlBufferResize(buffer.Ptr, C.uint(size)))
}

// xmlBufferSetAllocationScheme
func (buffer *Buffer) SetAllocationScheme(scheme AllocationScheme) {
	C.xmlBufferSetAllocationScheme(buffer.Ptr, C.xmlBufferAllocationScheme(scheme))
}

// xmlBufferShrink
func (buffer *Buffer) Shrink(length int) int {
	return int(C.xmlBufferShrink(buffer.Ptr, C.uint(length)))
}

// xmlBufferWriteChar/xmlBufferWriteCHAR
func (buffer *Buffer) WriteChar(str string) {
	ptr := C.CString(str)
	defer C.free(unsafe.Pointer(ptr))
	C.xmlBufferWriteChar(buffer.Ptr, ptr)
}

// xmlChildElementCount
func (node *Node) ChildCount() int {
	return int(C.xmlChildElementCount(node.Ptr))
}

// xmlCopyDoc
func (doc *Document) Copy(recursive int) *Document {
	cdoc := C.xmlCopyDoc(doc.Ptr, C.int(recursive))
	return makeDoc(cdoc)
}

// xmlCopyDtd
func (dtd *Dtd) Copy() *Dtd {
	return makeDtd(C.xmlCopyDtd(dtd.Ptr))
}

// xmlCopyNamespace
func (ns *Namespace) Copy(extended int) *Namespace {
	return makeNamespace(C.xmlCopyNamespace(ns.Ptr))
}

// xmlCopyNamespaceList
func (ns *Namespace) CopyList(extended int) *Namespace {
	return makeNamespace(C.xmlCopyNamespaceList(ns.Ptr))
}

// xmlCopyNode
func (node *Node) Copy(extended int) *Node {
	cnode := C.xmlCopyNode(node.Ptr, C.int(extended))
	return makeNode(cnode)
}

// xmlCopyNodeList
func (node *Node) CopyList() *Node {
	cnode := C.xmlCopyNodeList(node.Ptr)
	return makeNode(cnode)
}

// xmlGetProp
func (node *Node) GetAttribute(name string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cstr := C.xmlGetProp(node.Ptr, (*C.xmlChar)(unsafe.Pointer(cname)))
	if cstr == nil {
		return ""
	} else {
		return C.GoString((*C.char)(unsafe.Pointer(cstr)))
	}
}

// xmlHasProp
func (node *Node) HasAttribute(name string) *Attribute {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cattr := C.xmlHasProp(node.Ptr, (*C.xmlChar)(unsafe.Pointer(cname)))
	return makeAttribute(cattr)
}

// xmlCopyProp
func (attr *Attribute) Copy(target *Node) *Attribute {
	cattr := C.xmlCopyProp(target.Ptr, attr.Ptr)
	return makeAttribute(cattr)
}

// xmlCopyPropList
func (attr *Attribute) CopyList(target *Node) *Attribute {
	cattr := C.xmlCopyPropList(target.Ptr, attr.Ptr)
	return makeAttribute(cattr)
}

// xmlDocGetRootElement
func (doc *Document) Root() *Node {
	cnode := C.xmlDocGetRootElement(doc.Ptr)
	return makeNode(cnode)
}

// xmlDocSetRootElement
func (doc *Document) SetRoot(root *Node) *Node {
	cnode := C.xmlDocSetRootElement(doc.Ptr, root.Ptr)
	return makeNode(cnode)
}

// xmlFirstElementChild
func (node *Node) FirstChild() *Node {
	cnode := C.xmlFirstElementChild(node.Ptr)
	return makeNode(cnode)
}

func (node *Node) FirstChildNode() *Node {
	return makeNode(node.Ptr.children)
}

// xmlFreeDoc
func (doc *Document) Free() {
	C.xmlFreeDoc(doc.Ptr)
	doc.Ptr = nil
	doc.Node = nil
}

// xmlFreeDtd
func (dtd *Dtd) Free() {
	C.xmlFreeDtd(dtd.Ptr)
	dtd.Ptr = nil
}

// xmlFreeNode
func (node *Node) Free() {
	C.xmlFreeNode(node.Ptr)
	node.Ptr = nil
}

// xmlFreeNodeList
func (node *Node) FreeList() {
	C.xmlFreeNodeList(node.Ptr)
	node.Ptr = nil
}

// xmlFreeNs
func (ns *Namespace) Free() {
	C.xmlFreeNs(ns.Ptr)
	ns.Ptr = nil
}

// xmlFreeNsList
func (ns *Namespace) FreeList() {
	C.xmlFreeNsList(ns.Ptr)
	ns.Ptr = nil
}

// xmlFreeProp
func (attr *Attribute) Free() {
	C.xmlFreeProp(attr.Ptr)
	attr.Ptr = nil
}

// xmlFreePropList
func (attr *Attribute) FreeList() {
	C.xmlFreePropList(attr.Ptr)
	attr.Ptr = nil
}

// xmlGetLastChild
func (node *Node) LastChild() *Node {
	return makeNode(C.xmlGetLastChild(node.Ptr))
}

// xmlGetNodePath
func (node *Node) Path() string {
	cstr := C.xmlGetNodePath(node.Ptr)
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString((*C.char)(unsafe.Pointer(cstr)))
}

// xmlLastElementChild
func (node *Node) LastElementChild() *Node {
	return makeNode(C.xmlLastElementChild(node.Ptr))
}

// xmlNewChild
func (node *Node) NewChild(ns *Namespace, name string, content string) *Node {
	ptrn := C.CString(name)
	defer C.free(unsafe.Pointer(ptrn))
	ptrc := C.CString(content)
	defer C.free(unsafe.Pointer(ptrc))
	if ns != nil {
		return makeNode(C.xmlNewChild(node.Ptr, ns.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrc))))
	}
	return makeNode(C.xmlNewChild(node.Ptr, nil, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrc))))
}

// xmlNewComment
func NewComment(content string) *Node {
	ptr := C.CString(content)
	defer C.free(unsafe.Pointer(ptr))
	cnode := C.xmlNewComment((*C.xmlChar)(unsafe.Pointer(ptr)))
	return makeNode(cnode)
}

// xmlNewDoc
func NewDoc(version string) *Document {
	ptr := C.CString(version)
	defer C.free(unsafe.Pointer(ptr))
	cdoc := C.xmlNewDoc((*C.xmlChar)(unsafe.Pointer(ptr)))
	return makeDoc(cdoc)
}

// xmlNewDocComment
func (doc *Document) NewComment(content string) *Node {
	ptr := C.CString(content)
	defer C.free(unsafe.Pointer(ptr))
	return makeNode(C.xmlNewDocComment(doc.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr))))
}

// xmlNewDocFragment
func (doc *Document) NewFragment() *Node {
	return makeNode(C.xmlNewDocFragment(doc.Ptr))
}

// xmlNewDocNode
func (doc *Document) NewNode(ns *Namespace, name string, content string) *Node {
	ptrn := C.CString(name)
	defer C.free(unsafe.Pointer(ptrn))
	ptrc := C.CString(content)
	defer C.free(unsafe.Pointer(ptrc))
	if ns != nil {
		return makeNode(C.xmlNewDocNode(doc.Ptr, ns.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrc))))
	}
	return makeNode(C.xmlNewDocNode(doc.Ptr, nil, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrc))))
}

// xmlNewDocProp
func (doc *Document) NewProp(name string, value string) *Attribute {
	ptrn := C.CString(name)
	defer C.free(unsafe.Pointer(ptrn))
	ptrv := C.CString(value)
	defer C.free(unsafe.Pointer(ptrv))
	cattr := C.xmlNewDocProp(doc.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrv)))
	return makeAttribute(cattr)
}

// xmlNewDocRawNode
func (doc *Document) NewRawNode(ns *Namespace, name string, content string) *Node {
	ptrn := C.CString(name)
	defer C.free(unsafe.Pointer(ptrn))
	ptrc := C.CString(content)
	defer C.free(unsafe.Pointer(ptrc))
	if ns != nil {
		return makeNode(C.xmlNewDocRawNode(doc.Ptr, ns.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrc))))
	}
	return makeNode(C.xmlNewDocRawNode(doc.Ptr, nil, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrc))))
}

// xmlNewDocText
func (doc *Document) NewText(content string) *TextNode {
	ptr := C.CString(content)
	defer C.free(unsafe.Pointer(ptr))
	return makeTextNode(C.xmlNewDocText(doc.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr))))
}

// xmlNewDtd
func (doc *Document) NewDtd(name string, ExternalID string, SystemID string) *Dtd {
	ptrn := C.CString(name)
	defer C.free(unsafe.Pointer(ptrn))
	ptre := C.CString(ExternalID)
	defer C.free(unsafe.Pointer(ptre))
	ptrs := C.CString(SystemID)
	defer C.free(unsafe.Pointer(ptrs))
	cdtd := C.xmlNewDtd(doc.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptre)), (*C.xmlChar)(unsafe.Pointer(ptrs)))
	return makeDtd(cdtd)
}

// xmlNewNode
func NewNode(ns *Namespace, name string) *Node {
	ptr := C.CString(name)
	defer C.free(unsafe.Pointer(ptr))
	if ns != nil {
		return makeNode(C.xmlNewNode(ns.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr))))
	}
	return makeNode(C.xmlNewNode(nil, (*C.xmlChar)(unsafe.Pointer(ptr))))
}

// xmlNewNs
func (node *Node) NewNs(href string, prefix string) *Namespace {
	ptrh := C.CString(href)
	defer C.free(unsafe.Pointer(ptrh))
	ptrp := C.CString(prefix)
	defer C.free(unsafe.Pointer(ptrp))
	return makeNamespace(C.xmlNewNs(node.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrh)), (*C.xmlChar)(unsafe.Pointer(ptrp))))
}

// xmlNewProp
func (node *Node) NewAttribute(name string, value string) *Attribute {
	ptrn := C.CString(name)
	defer C.free(unsafe.Pointer(ptrn))
	ptrv := C.CString(value)
	defer C.free(unsafe.Pointer(ptrv))
	cattr := C.xmlNewProp(node.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrv)))
	return makeAttribute(cattr)
}

// xmlNewText
func NewText(content string) *TextNode {
	ptr := C.CString(content)
	defer C.free(unsafe.Pointer(ptr))
	return makeTextNode(C.xmlNewText((*C.xmlChar)(unsafe.Pointer(ptr))))
}

// xmlNewTextChild
func (node *Node) NewTextChild(ns *Namespace, name string, content string) *TextNode {
	ptrn := C.CString(name)
	defer C.free(unsafe.Pointer(ptrn))
	ptrc := C.CString(content)
	defer C.free(unsafe.Pointer(ptrc))
	if ns == nil {
		return makeTextNode(C.xmlNewTextChild(node.Ptr, nil, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrc))))
	}
	return makeTextNode(C.xmlNewTextChild(node.Ptr, ns.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrc))))
}

// xmlNextElementSibling
func (node *Node) NextSibling() *Node {
	return makeNode(C.xmlNextElementSibling(node.Ptr))
}

// xmlNodeAddContent
func (node *Node) AddContent(content string) {
	ptr := C.CString(content)
	defer C.free(unsafe.Pointer(ptr))
	C.xmlNodeAddContent(node.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr)))
}

// xmlNodeDump
func (doc *Document) NodeDump(buf *Buffer, cur *Node, level int, format int) int {
	return int(C.xmlNodeDump(buf.Ptr, doc.Ptr, cur.Ptr, C.int(level), C.int(format)))
}

// xmlNodeGetContent
func (node *Node) GetContent() string {
	content := (*C.char)(unsafe.Pointer(C.xmlNodeGetContent(node.Ptr)))
	defer C.free(unsafe.Pointer(content))
	return C.GoString(content)
}

// xmlNodeListGetString
func (node *Node) ListGetString(inLine bool) string {
	ptr := node.Ptr
	docptr := C.xmlDocPtr((*C.xmlDoc)(ptr.doc))
	cInLine := C.int(0)
	if inLine {
		cInLine = C.int(1)
	}
	str := (*C.char)(unsafe.Pointer(C.xmlNodeListGetString(docptr, ptr, cInLine)))
	defer C.free(unsafe.Pointer(str))
	return C.GoString(str)
}

// xmlNodeSetContent
func (node *Node) SetContent(content string) {
	ptr := C.CString(content)
	defer C.free(unsafe.Pointer(ptr))
	C.xmlNodeSetContent(node.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr)))
}

// xmlNodeSetName
func (node *Node) SetName(name string) {
	ptr := C.CString(name)
	defer C.free(unsafe.Pointer(ptr))
	C.xmlNodeSetName(node.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr)))
}

// xmlPreviousElementSibling
func (node *Node) PreviousSibling() *Node {
	return makeNode(C.xmlPreviousElementSibling(node.Ptr))
}

// xmlRemoveProp
func RemoveAttribute(attr *Attribute) int {
	return int(C.xmlRemoveProp(attr.Ptr))
}

// xmlReplaceNode
func (old *Node) Replace(cur *Node) *Node {
	return makeNode(C.xmlReplaceNode(old.Ptr, cur.Ptr))
}

// xmlXIncludeProcessFlags
func (doc *Document) ProcessXInclude(flags ParserOption) int {
	return int(C.xmlXIncludeProcessFlags(doc.Ptr, C.int(flags)))
}

// xmlSaveFile
func (doc *Document) SaveFile(filename string) int {
	ptr := C.CString(filename)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlSaveFile(ptr, doc.Ptr))
}

// xmlSaveFileEnc
func (doc *Document) SaveFileEnc(filename string, encoding string) int {
	ptrf := C.CString(filename)
	defer C.free(unsafe.Pointer(ptrf))
	ptre := C.CString(encoding)
	defer C.free(unsafe.Pointer(ptre))
	return int(C.xmlSaveFileEnc(ptrf, doc.Ptr, ptre))
}

// xmlSaveFormatFile
func (doc *Document) SaveFormatFile(filename string, format int) int {
	ptr := C.CString(filename)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlSaveFormatFile(ptr, doc.Ptr, C.int(format)))
}

// xmlSaveFormatFileEnc
func (doc *Document) SaveFormatFileEnc(filename string, encoding string, format int) int {
	ptrf := C.CString(filename)
	defer C.free(unsafe.Pointer(ptrf))
	ptre := C.CString(encoding)
	defer C.free(unsafe.Pointer(ptre))
	return int(C.xmlSaveFormatFileEnc(ptrf, doc.Ptr, ptre, C.int(format)))
}

// xmlSetBufferAllocationScheme
func SetAllocationScheme(scheme AllocationScheme) {
	C.xmlSetBufferAllocationScheme(C.xmlBufferAllocationScheme(scheme))
}

// xmlSetCompressMode
func SetCompressionLevel(level int) {
	C.xmlSetCompressMode(C.int(level))
}

// xmlSetDocCompressMode
func (doc *Document) SetCompressionLevel(level int) {
	C.xmlSetDocCompressMode(doc.Ptr, C.int(level))
}

// xmlSetProp
func (node *Node) SetAttribute(name string, value string) *Attribute {
	ptrn := C.CString(name)
	defer C.free(unsafe.Pointer(ptrn))
	ptrv := C.CString(value)
	defer C.free(unsafe.Pointer(ptrv))
	cattr := C.xmlSetProp(node.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrn)), (*C.xmlChar)(unsafe.Pointer(ptrv)))
	return makeAttribute(cattr)
}

// xmlUnSetProp
func (node *Node) UnsetAttribute(name string) bool {
	ptrn := C.CString(name)
	defer C.free(unsafe.Pointer(ptrn))
	res := C.xmlUnsetProp(node.Ptr, (*C.xmlChar)(unsafe.Pointer(ptrn)))
	return res == 0
}

// xmlTextConcat
func (node *TextNode) Concat(content string) int {
	ptr := C.CString(content)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlTextConcat(node.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr)), C.int(len(content))))
}

// xmlTextMerge
func (first *TextNode) Merge(second *Node) *Node {
	return makeNode(C.xmlTextMerge(first.Ptr, second.Ptr))
}

// xmlUnlinkNode
func (node *Node) Unlink() {
	C.xmlUnlinkNode(node.Ptr)
}

// xmlUnsetNsProp
func (node *Node) UnsetNsProp(ns *Namespace, name string) int {
	ptr := C.CString(name)
	defer C.free(unsafe.Pointer(ptr))
	cint := C.xmlUnsetNsProp(node.Ptr, ns.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr)))
	return int(cint)
}

// xmlUnsetProp
func (node *Node) UnsetProp(name string) int {
	ptr := C.CString(name)
	defer C.free(unsafe.Pointer(ptr))
	cint := C.xmlUnsetProp(node.Ptr, (*C.xmlChar)(unsafe.Pointer(ptr)))
	return int(cint)
}

// xmlValidateNCName
func ValidateNCName(value string, space int) int {
	ptr := C.CString(value)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlValidateNCName((*C.xmlChar)(unsafe.Pointer(ptr)), C.int(space)))
}

// xmlValidateNMToken
func ValidateNMToken(value string, space int) int {
	ptr := C.CString(value)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlValidateNMToken((*C.xmlChar)(unsafe.Pointer(ptr)), C.int(space)))
}

// xmlValidateName
func ValidateName(value string, space int) int {
	ptr := C.CString(value)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlValidateName((*C.xmlChar)(unsafe.Pointer(ptr)), C.int(space)))
}

// xmlValidateQName
func ValidateQName(value string, space int) int {
	ptr := C.CString(value)
	defer C.free(unsafe.Pointer(ptr))
	return int(C.xmlValidateQName((*C.xmlChar)(unsafe.Pointer(ptr)), C.int(space)))
}
