package markdown

import (
	"errors"
	"unsafe"

	"github.com/jbussdieker/golibxml"
	"github.com/yuin/goldmark/util"
)

func NewTreeWriter(chunk string, filename string) TreeWriter {
	parser := golibxml.CreateHTMLPushParser(chunk, filename)
	parser.UseOptions(golibxml.HTML_PARSE_RECOVER | golibxml.HTML_PARSE_NOERROR)
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
	context *golibxml.HTMLParser
}

func (w TreeWriter) Terminate(chunk string) {
	w.context.Terminate(chunk)
}

func (w TreeWriter) Document() *golibxml.Document {
	return w.context.MyDoc()
}

func (w TreeWriter) Free() {
	w.context.Free()
}

func (w *TreeWriter) Write(p []byte) (n int, err error) {
	if res := w.context.ParseChunk(*(*string)(unsafe.Pointer(&p))); res != 0 {
		return 0, errors.New("TreeWriter error")
	}
	return len(p), nil
}
