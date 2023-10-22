package markup

/*
#include <libxml/tree.h>

static inline void free_string(char* s) { free(s); }
static inline xmlChar *to_xmlcharptr(const char *s) { return (xmlChar *)s; }
static inline char *to_charptr(const xmlChar *s) { return (char *)s; }
*/
import "C"
import "unsafe"

func (doc *Document) String() string {
	buf := NewBuffer()
	defer buf.Free()
	doc.NodeDump(buf, doc.Node, 0, 0)
	return buf.Content()
}

func (node *Node) Document() *Document {
	return makeDoc(C.xmlDocPtr(unsafe.Pointer(node.Ptr.doc)))
}

func (node *Node) String() string {
	buf := NewBuffer()
	defer buf.Free()
	node.Document().NodeDump(buf, node, 0, 0)
	return buf.Content()
}

func (node *Node) Children() *Node {
	return makeNode(C.xmlNodePtr(unsafe.Pointer(node.Ptr.children)))
}

func (node *Node) Parent() *Node {
	return makeNode(C.xmlNodePtr(unsafe.Pointer(node.Ptr.parent)))
}

func (node *Node) Type() ElementType {
	return ElementType(node.Ptr._type)
}

func (node *Node) Name() string {
	return C.GoString(C.to_charptr(node.Ptr.name))
}

func (node *Node) Next() *Node {
	return makeNode(C.xmlNodePtr(unsafe.Pointer(node.Ptr.next)))
}

func (node *Node) Attributes() *Attribute {
	return makeAttribute(C.xmlAttrPtr(unsafe.Pointer(node.Ptr.properties)))
}

func (node *Node) Namespace() *Namespace {
	return makeNamespace(C.xmlNsPtr(unsafe.Pointer(node.Ptr.ns)))
}

func (attr *Attribute) Type() ElementType {
	return ElementType(attr.Ptr._type)
}

func (attr *Attribute) Name() string {
	return C.GoString(C.to_charptr(attr.Ptr.name))
}

func (attr *Attribute) Children() *Node {
	return makeNode(C.xmlNodePtr(unsafe.Pointer(attr.Ptr.children)))
}

func (attr *Attribute) Next() *Attribute {
	return makeAttribute(C.xmlAttrPtr(unsafe.Pointer(attr.Ptr.next)))
}

func (attr *Attribute) Namespace() *Namespace {
	return makeNamespace(C.xmlNsPtr(unsafe.Pointer(attr.Ptr.ns)))
}

func (ns *Namespace) Href() string {
	return C.GoString(C.to_charptr(ns.Ptr.href))
}

func (ns *Namespace) Prefix() string {
	return C.GoString(C.to_charptr(ns.Ptr.prefix))
}

func (elem ElementType) GoString() string {
	return elem.String()
}

func (elem ElementType) String() string {
	switch elem {
	case XML_ELEMENT_NODE:
		return "Element"
	case XML_ATTRIBUTE_NODE:
		return "Attribute"
	case XML_TEXT_NODE:
		return "Text"
	case XML_CDATA_SECTION_NODE:
		return "CDATA"
	case XML_ENTITY_REF_NODE:
		return "Entity Reference"
	case XML_ENTITY_NODE:
		return "Entity"
	case XML_PI_NODE:
		return "Processing Instruction"
	case XML_COMMENT_NODE:
		return "Comment"
	case XML_DOCUMENT_NODE:
		return "Document"
	case XML_DOCUMENT_TYPE_NODE:
		return "Document Type"
	case XML_DOCUMENT_FRAG_NODE:
		return "Document Fragment"
	case XML_NOTATION_NODE:
		return "Notation"
	case XML_HTML_DOCUMENT_NODE:
		return "HTML Document"
	case XML_DTD_NODE:
		return "DTD"
	case XML_ELEMENT_DECL:
		return "Element Declaration"
	case XML_ATTRIBUTE_DECL:
		return "Attribute Declaration"
	case XML_ENTITY_DECL:
		return "Entity Declaration"
	case XML_NAMESPACE_DECL:
		return "Namespace"
	case XML_XINCLUDE_START:
		return "XInclude Start"
	case XML_XINCLUDE_END:
		return "XInclude End"
	default:
		return "Unknown Type"
	}
}
