package markup

/*
#include <libxml/tree.h>
*/
import "C"

import (
	"bufio"
	"strings"
)

var (
	voidElements     map[string]bool
	preElements      map[string]bool
	literalElements  map[string]bool
	scriptElements   map[string]bool
	quoteCharacters  *strings.Replacer
	escapeCharacters *strings.Replacer
)

func hasParent(node *Node, elements map[string]bool) bool {
	parent := makeNode((C.xmlNode)(*node.Ptr).parent)
	if _, ok := elements[parent.Name()]; ok {
		return true
	}
	return false
}

type HTML5Serializer struct {
	writer *bufio.Writer
}

func (s *HTML5Serializer) Serialize(doc *Document) error {
	err := s.serializeDocument(doc)
	if err != nil {
		return err
	}
	s.writer.Flush()
	return nil
}

func NewHTML5Serializer(writer *bufio.Writer) *HTML5Serializer {
	return &HTML5Serializer{writer}
}

func (s *HTML5Serializer) Write(data []byte) (int, error) {
	return s.writer.Write(data)
}

func (s *HTML5Serializer) WriteString(data string) (int, error) {
	return s.writer.WriteString(data)
}

func (s *HTML5Serializer) serializeDocument(doc *Document) error {
	s.serializeFragment(doc.Children())
	return nil
}

func (s *HTML5Serializer) serializeFragment(current *Node) error {
	for current != nil {
		switch current.Type() {
		case XML_ELEMENT_NODE:
			s.serializeElement(current)
		case XML_TEXT_NODE:
			s.serializeText(current)
		case XML_CDATA_SECTION_NODE:
			s.serializeText(current)
		case XML_COMMENT_NODE:
			s.serializeComment(current)
		case XML_PI_NODE:
			s.serializeProcessingInstruction(current)
		case XML_DOCUMENT_TYPE_NODE:
			s.serializeDocumentType(current)
		default:
			panic("invalid state")
		}
		current = current.Next()
	}
	return nil
}

func (s *HTML5Serializer) serializeElement(node *Node) error {
	name := node.Name()
	if _, err := s.Write([]byte{'<'}); err != nil {
		return err
	}
	if _, err := s.WriteString(name); err != nil {
		return err
	}
	if err := s.serializeAttributes(node.Attributes()); err != nil {
		return err
	}
	if _, err := s.Write([]byte{'>'}); err != nil {
		return err
	}
	if _, ok := voidElements[name]; ok {
		return nil
	}
	if _, ok := literalElements[name]; ok {
		s.Write([]byte{'\n'})
	}
	if err := s.serializeFragment(node.Children()); err != nil {
		return err
	}
	if _, err := s.Write([]byte{'<', '/'}); err != nil {
		return err
	}
	if _, err := s.WriteString(name); err != nil {
		return err
	}
	if _, err := s.Write([]byte{'>'}); err != nil {
		return err
	}
	return nil
}

func (s *HTML5Serializer) serializeText(node *Node) error {
	content := node.GetContent()
	if hasParent(node, scriptElements) {
		s.WriteString(content)
	} else {
		s.WriteString(escapeCharacters.Replace(content))
	}
	return nil
}

func (s *HTML5Serializer) serializeComment(node *Node) error {
	if _, err := s.Write([]byte{'<', '!', '-', '-'}); err != nil {
		return err
	}
	if _, err := s.WriteString(node.GetContent()); err != nil {
		return err
	}
	if _, err := s.Write([]byte{'-', '-', '>'}); err != nil {
		return err
	}
	return nil
}

func (s *HTML5Serializer) serializeProcessingInstruction(node *Node) error {
	if _, err := s.Write([]byte{'<', '?'}); err != nil {
		return err
	}
	//TODO write target attribute
	if _, err := s.WriteString(node.GetContent()); err != nil {
		return err
	}
	if _, err := s.Write([]byte{'>'}); err != nil {
		return err
	}
	return nil
}

func (s *HTML5Serializer) serializeDocumentType(node *Node) error {
	if _, err := s.Write([]byte{'<', '!', 'D', 'O', 'C', 'T', 'Y', 'P', 'E', ' '}); err != nil {
		return err
	}
	if _, err := s.WriteString(node.Name()); err != nil {
		return err
	}
	if _, err := s.Write([]byte{'>'}); err != nil {
		return err
	}
	return nil
}

func (s *HTML5Serializer) serializeAttributes(attr *Attribute) error {
	for attr != nil {
		if err := s.serializeAttribute(attr); err != nil {
			return err
		}
		attr = attr.Next()
	}
	return nil
}

func (s *HTML5Serializer) serializeAttribute(attr *Attribute) error {
	var (
		name    string
		content strings.Builder
	)
	name = attr.Name()
	for child := attr.Children(); child != nil; child = child.Next() {
		if child.Type() == XML_TEXT_NODE {
			content.WriteString(child.GetContent())
		} else {
			panic("invalid state")
		}
	}
	if _, err := s.Write([]byte{' '}); err != nil {
		return err
	}
	if _, err := s.WriteString(name); err != nil {
		return err
	}
	if content.Len() > 0 {
		if _, err := s.Write([]byte{'=', '"'}); err != nil {
			return err
		}
		if _, err := s.WriteString(quoteCharacters.Replace(content.String())); err != nil {
			return err
		}
		if _, err := s.Write([]byte{'"'}); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	voidElements = make(map[string]bool)
	for _, name := range strings.Split("area,base,basefont,bgsound,br,col,embed,frame,hr,img,input,link,meta,param,spacer,wbr", ",") {
		voidElements[name] = true
	}

	preElements = make(map[string]bool)
	for _, name := range strings.Split("pre,textarea,listing", ",") {
		preElements[name] = true
	}

	literalElements = make(map[string]bool)
	for _, name := range strings.Split("style,script,xmp,iframe,noembed,noframes,noscript,plaintext", ",") {
		literalElements[name] = true
	}

	scriptElements = make(map[string]bool)
	for _, name := range strings.Split("style,script,xmp,iframe,noembed,noframes,noscript,plaintext", ",") {
		scriptElements[name] = true
	}

	escapeCharacters = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;", "\xa0", "&nbsp;")
	quoteCharacters = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;", "\xa0", "&nbsp;", "\"", "&quote")
}
