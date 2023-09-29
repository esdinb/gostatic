package markdown

import (
	"errors"
	"unsafe"

	"gostatic/pkg/markup"

	"github.com/yuin/goldmark/util"
)

func NewTreeWriter(doc *markup.Document, node *markup.Node) TreeWriter {
	parser := markup.CreateHTML5Parser(doc, node)
	return TreeWriter{parser}
}

// implements the html.Writer interface
type HtmlWriter struct {
}

func (w HtmlWriter) Write(writer util.BufWriter, source []byte) {
	writer.Write(source)
}

func (w HtmlWriter) RawWrite(writer util.BufWriter, source []byte) {
	w.Write(writer, source)
}

func (w HtmlWriter) SecureWrite(writer util.BufWriter, source []byte) {
	w.Write(writer, source)
}

// implements the io.Writer interface
type TreeWriter struct {
	context *markup.HTML5Parser
}

func (w TreeWriter) Terminate() {
	w.context.Terminate()
}

func (w TreeWriter) Document() *markup.Document {
	return w.context.MyDoc()
}

func (w TreeWriter) Free() {
	w.context.Free()
}

func (w *TreeWriter) Write(p []byte) (n int, err error) {
	if res := w.context.ParseChunk(*(*string)(unsafe.Pointer(&p))); res < 0 {
		return 0, errors.New("TreeWriter error")
	}
	return len(p), nil
}
