package markup

import "io"

type HTML5Serializer struct {
	writer io.Writer
}

func (s *HTML5Serializer) Serialize(doc *Document) error {
	return nil
}

func NewHTML5Serializer(writer io.Writer) *HTML5Serializer {
}
