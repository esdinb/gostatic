package builder

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gostatic/pkg/markup"
	"gostatic/pkg/transformer"
)

func ReadFile(srcPath string) *markup.Document {
	return markup.ReadFile(
		srcPath,
		"UTF-8",
		markup.XML_PARSE_RECOVER&
			markup.XML_PARSE_NONET&
			markup.XML_PARSE_PEDANTIC&
			markup.XML_PARSE_NOBLANKS&
			markup.XML_PARSE_XINCLUDE,
	)
}

func WriteFile(doc *markup.Document, destPath string) error {
	err := os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return err
	}
	options := markup.SaveOption(0)
	var length int
	if strings.HasSuffix(destPath, ".html") {
		options |= markup.XML_SAVE_NO_DECL
		options |= markup.XML_SAVE_NO_EMPTY
		options |= markup.XML_SAVE_AS_XML
	}
	saveCtx := markup.SaveToFilename(destPath, "utf-8", options)
	if saveCtx == nil {
		return errors.New("failed to create save context")
	}
	defer saveCtx.Free()
	length = saveCtx.SaveDoc(doc)
	if length == -1 {
		return errors.New("failed to write transformation to file")
	} else {
		return nil
	}
}

func CopyFile(srcPath string, destPath string) error {
	var err error
	err = os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return err
	}
	reader, err := os.OpenFile(srcPath, os.O_RDONLY, fs.ModeExclusive)
	if err != nil {
		return err
	}
	writer, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}
	return nil
}

type BuildCommand struct {
	Name string
	Args []string
}

type BuildTransformation string

type Pipeline []BuildTransformation

func (b *BuildTransformation) Parse() BuildCommand {
	parts := strings.SplitN(string(*b), ":", 3)
	return BuildCommand{parts[0], parts[1:]}
}

type BuildSection struct {
	In       string
	Out      string
	Pipeline Pipeline
}

func NewBuildSection(in string, out string, pipeline Pipeline) BuildSection {
	return BuildSection{in, out, pipeline}
}

func (p *Pipeline) Transform(ctx context.Context) (context.Context, error) {
	var status transformer.Status
	var err error
	for _, command := range *p {
		cmd := command.Parse()
		fn := transformer.Registry.Lookup(cmd.Name)
		if fn == nil {
			continue
		}
		ctx, status, err = fn(ctx, cmd.Args)
		if err != nil {
			return ctx, err
		}
		if status == transformer.Stop {
			break
		} else {
			continue
		}
	}
	return ctx, nil
}

func (b *BuildSection) ProcessFile(inPath string, outPath string, rootPath string) error {
	var doc *markup.Document
	var err error
	if inPath == "-" {
		return errors.New("cannot read from stdin")
	} else {
		doc = ReadFile(inPath)
	}

	if doc == nil {
		return errors.New("error processing input file")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, transformer.InPathContextKey, inPath)
	ctx = context.WithValue(ctx, transformer.OutPathContextKey, outPath)
	ctx = context.WithValue(ctx, transformer.RootPathContextKey, rootPath)
	ctx = context.WithValue(ctx, transformer.DocumentContextKey, doc)
	ctx = context.WithValue(ctx, transformer.ParamsContextKey, []string{})
	ctx = context.WithValue(ctx, transformer.StringParamsContextKey, []string{})

	ctx, err = b.Pipeline.Transform(ctx)
	if err != nil {
		return err
	}

	doc = ctx.Value(transformer.DocumentContextKey).(*markup.Document)
	defer doc.Free()

	outPath = ctx.Value(transformer.OutPathContextKey).(string)

	err = WriteFile(doc, outPath)
	if err != nil {
		return err
	}

	return nil
}

func (b *BuildSection) CopyFile(inPath string, outPath string) error {
	return CopyFile(inPath, outPath)
}

func (b *BuildSection) Build(rootPath string) error {
	var err error
	absoluteRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return err
	}

	absPath := absoluteRootPath
	outPath := b.Out
	outPathIsADir := strings.HasSuffix(outPath, string(os.PathSeparator))
	if outPath != "-" {
		absPath = filepath.Join(absoluteRootPath, outPath)

		fileInfo, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				newPath := absPath
				if !outPathIsADir {
					newPath = filepath.Dir(absPath)
				}
				err = os.MkdirAll(newPath, 0755)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			outPathIsADir = fileInfo.IsDir()
		}
		outPath = absPath
	} else {
		return errors.New("cannot write to stdout")
	}

	var matches []string
	matches, err = filepath.Glob(filepath.Join(rootPath, b.In))
	if err != nil {
		return err
	}

	for i := range matches {
		inPath := matches[i]
		fileInfo, err := os.Stat(inPath)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			continue
		}
		if outPathIsADir {
			relPath, err := filepath.Rel(rootPath, inPath)
			if err != nil {
				return err
			}
			outPath = filepath.Join(absPath, relPath)
		}
		if strings.HasSuffix(inPath, ".html") || strings.HasSuffix(inPath, ".xml") {
			err = b.ProcessFile(inPath, outPath, absoluteRootPath)
		} else {
			err = b.CopyFile(inPath, outPath)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
