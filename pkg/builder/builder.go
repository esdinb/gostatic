package builder

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	builder_context "gostatic/pkg/builder/context"
	"gostatic/pkg/markup"
	"gostatic/pkg/transformer"
)

type FormatterFunc = func(ctx context.Context) error
type LoaderFunc = func(ctx context.Context) (context.Context, error)
type LookupFunc = interface{ FormatterFunc | LoaderFunc }

const (
	parseOptions markup.ParserOption = markup.XML_PARSE_RECOVER &
		markup.XML_PARSE_NONET &
		markup.XML_PARSE_PEDANTIC &
		markup.XML_PARSE_NOBLANKS &
		markup.XML_PARSE_XINCLUDE &
		markup.XML_PARSE_HUGE
)

func lookupFunc[T LookupFunc](fs map[string]T, none T) func(string) T {
	return func(ext string) T {
		f, ok := fs[ext]
		if ok {
			return f
		} else {
			return none
		}
	}
}

var lookupFormatter = lookupFunc(map[string]FormatterFunc{
	".xml":  xmlFormatter,
	".html": htmlFormatter,
}, copyFormatter)

var lookupLoader = lookupFunc(map[string]LoaderFunc{
	".xml":  xmlLoader,
	".html": htmlLoader,
}, nopLoader)

func xmlLoader(ctx context.Context) (context.Context, error) {
	inPath, ok := ctx.Value(builder_context.InPathContextKey).(string)
	if !ok {
		return nil, errors.New("missing input path for xml loader")
	}

	doc := markup.ReadFile(inPath, "UTF-8", parseOptions)
	if doc == nil {
		return nil, errors.New("unable to load xml file")
	}

	return context.WithValue(ctx, builder_context.DocumentContextKey, doc), nil
}

func htmlLoader(ctx context.Context) (context.Context, error) {
	inPath, ok := ctx.Value(builder_context.InPathContextKey).(string)
	if !ok {
		return nil, errors.New("missing input path for html loader")
	}

	doc := markup.ReadHTMLFile(inPath, parseOptions)
	if doc == nil {
		return nil, errors.New("unable to load xml file")
	}

	return context.WithValue(ctx, builder_context.DocumentContextKey, doc), nil
}

func nopLoader(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func xmlFormatter(ctx context.Context) error {
	var (
		options markup.SaveOption
		saveCtx *markup.SaveContext
		doc     *markup.Document
		outFile *os.File
		ok      bool
	)

	options = markup.SaveOption(0)

	doc, ok = ctx.Value(builder_context.DocumentContextKey).(*markup.Document)
	if !ok {
		return errors.New("missing transformation document for xml formatter")
	}

	outFile, ok = ctx.Value(builder_context.OutFileContextKey).(*os.File)
	if !ok {
		outPath, ok := ctx.Value(builder_context.OutPathContextKey).(string)
		if !ok {
			return errors.New("missing output file path for xml formatter")
		} else {
			var err error
			outFile, err = openFileForWriting(outPath)
			if err != nil {
				return err
			}
		}
	}

	saveCtx = markup.SaveToIO(outFile, "UTF-8", options)
	if saveCtx == nil {
		return errors.New("failed to create save xml formatter context")
	}
	defer saveCtx.Free()

	if err := saveCtx.SaveDoc(doc); err != nil {
		return err
	}

	return nil
}

func htmlFormatter(ctx context.Context) error {
	var (
		options markup.SaveOption
		saveCtx *markup.SaveContext
		doc     *markup.Document
		outFile *os.File
		ok      bool
		err     error
	)

	options = markup.SaveOption(markup.XML_SAVE_NO_DECL | markup.XML_SAVE_NO_EMPTY | markup.XML_SAVE_AS_XML)

	doc, ok = ctx.Value(builder_context.DocumentContextKey).(*markup.Document)
	if !ok {
		return errors.New("missing transformation document for html formatter")
	}

	outFile, ok = ctx.Value(builder_context.OutFileContextKey).(*os.File)
	if !ok {
		outPath, ok := ctx.Value(builder_context.OutPathContextKey).(string)
		if !ok {
			return errors.New("missing output file path for html formatter")
		} else {
			outFile, err = openFileForWriting(outPath)
			if err != nil {
				return err
			}
		}
	}

	saveCtx = markup.SaveToIO(outFile, "UTF-8", options)
	if saveCtx == nil {
		return errors.New("failed to create html formatter context")
	}
	defer saveCtx.Free()

	if err := saveCtx.SaveDoc(doc); err != nil {
		return err
	}

	return nil
}

func copyFormatter(ctx context.Context) error {
	var (
		inFile  *os.File
		outFile *os.File
		ok      bool
		err     error
	)

	inFile, ok = ctx.Value(builder_context.InFileContextKey).(*os.File)
	if !ok {
		inPath, ok := ctx.Value(builder_context.InPathContextKey).(string)
		if !ok {
			return errors.New("missing input file for copy")
		}
		inFile, err = openFileForReading(inPath)
		if err != nil {
			return errors.New("cannot open input file for copy")
		}
	}

	outFile, ok = ctx.Value(builder_context.OutFileContextKey).(*os.File)
	if !ok {
		outPath, ok := ctx.Value(builder_context.OutPathContextKey).(string)
		if !ok {
			return errors.New("missing output file for copy")
		}
		outFile, err = openFileForWriting(outPath)
		if err != nil {
			return errors.New("cannot open output file for copy")
		}
	}

	if _, err := io.Copy(outFile, inFile); err != nil {
		return err
	}

	return nil
}

func openFileForReading(inPath string) (*os.File, error) {
	if reader, err := os.OpenFile(inPath, os.O_RDONLY, fs.ModeExclusive); err != nil {
		return nil, err
	} else {
		return reader, nil
	}
}

func openFileForWriting(outPath string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return nil, err
	}
	if writer, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return nil, err
	} else {
		return writer, nil
	}
}

func freeContextDocument(ctx context.Context) {
	doc, ok := ctx.Value(builder_context.DocumentContextKey).(*markup.Document)
	if ok {
		doc.Free()
	}
}

func freeBuildContext(ctx context.Context) {
	var (
		inFile  *os.File
		outFile *os.File
		ok      bool
	)

	freeContextDocument(ctx)

	inFile, ok = ctx.Value(builder_context.InFileContextKey).(*os.File)
	if ok {
		inFile.Close()
	}

	outFile, ok = ctx.Value(builder_context.OutFileContextKey).(*os.File)
	if ok {
		outFile.Close()
	}
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
			return ctx, errors.New(fmt.Sprintf("unknown transform name: %s", cmd.Name))
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

func (b *BuildSection) ProcessFile(ctx context.Context) error {

	ctx, err := b.Pipeline.Transform(ctx)
	if err != nil {
		return err
	}

	if format, ok := ctx.Value(builder_context.FormatterContextKey).(FormatterFunc); ok {
		err := format(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BuildSection) Build(ctx context.Context, rootPath string) error {
	var (
		outFile   *os.File      = nil
		formatter FormatterFunc = nil
		err       error
	)
	rootPathAbsolute, err := filepath.Abs(rootPath)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, builder_context.InPathContextKey, b.In)
	ctx = context.WithValue(ctx, builder_context.OutPathContextKey, b.Out)
	ctx = context.WithValue(ctx, builder_context.RootPathContextKey, rootPath)
	ctx = context.WithValue(ctx, builder_context.ParamsContextKey, []string{})
	ctx = context.WithValue(ctx, builder_context.StringParamsContextKey, []string{})

	if logger, ok := ctx.Value(builder_context.LoggerContextKey).(*log.Logger); ok {
		markup.SetErrorReporting(logger)
	}

	inPath := b.In
	outPath := b.Out
	outPathIsDir := strings.HasSuffix(outPath, string(os.PathSeparator))
	absPath := filepath.Join(rootPathAbsolute, outPath)

	if outPath == "-" {
		outFile = os.Stdout
		ctx = context.WithValue(ctx, builder_context.OutFileContextKey, os.Stdout)
		formatter = lookupFormatter(filepath.Ext(inPath))
	} else {
		if info, err := os.Stat(absPath); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				newPath := absPath
				if !outPathIsDir {
					newPath = filepath.Dir(absPath)
				}

				if err := os.MkdirAll(newPath, 0755); err != nil {
					return err
				}
			} else {
				return err
			}
		} else if info.IsDir() {
			outPathIsDir = true
		}

		if !outPathIsDir {
			formatter = lookupFormatter(filepath.Ext(outPath))
		}

		outPath = absPath
	}

	if b.In == "-" {
		var (
			reader   *bufio.Reader
			readNext bool
			bytes    []byte
		)
		if outPathIsDir {
			return errors.New("cannot pipe to a directory")
		}
		if outFile != nil {
			formatter = lookupFormatter(".xml")
		}
		reader = bufio.NewReader(os.Stdin)
		readNext = true
		for readNext {
			bytes, err = reader.ReadBytes(0)
			if err != nil {
				if errors.Is(err, io.EOF) {
					readNext = false
				} else {
					return err
				}
			}
			if len(bytes) > 0 {
				buildCtx := context.WithValue(
					ctx,
					builder_context.DocumentContextKey,
					markup.ReadMemory(bytes, b.In, "UTF-8", parseOptions),
				)
				buildCtx = context.WithValue(buildCtx, builder_context.FormatterContextKey, formatter)
				err = b.ProcessFile(buildCtx)
				freeContextDocument(buildCtx)
				if err != nil {
					return err
				}
				if outFile != nil && readNext {
					outFile.Write([]byte{0})
				}
			}
		}
		freeBuildContext(ctx)
	} else {
		matches, err := filepath.Glob(filepath.Join(rootPath, b.In))
		if err != nil {
			return err
		}

		for i := range matches {
			inPath = matches[i]
			if info, err := os.Stat(inPath); err != nil {
				return err
			} else if info.IsDir() {
				continue
			}
			if outFile != nil {
				formatter = lookupFormatter(filepath.Ext(inPath))
			} else if outPathIsDir {
				if relPath, err := filepath.Rel(rootPath, inPath); err != nil {
					return err
				} else {
					outPath = filepath.Join(absPath, relPath)
				}
				formatter = lookupFormatter(filepath.Ext(outPath))
			}
			buildCtx := context.WithValue(ctx, builder_context.InPathContextKey, inPath)
			buildCtx = context.WithValue(buildCtx, builder_context.OutPathContextKey, outPath)
			buildCtx = context.WithValue(buildCtx, builder_context.FormatterContextKey, formatter)
			loader := lookupLoader(filepath.Ext(inPath))
			buildCtx, err = loader(buildCtx)
			if err != nil {
				return err
			}

			err = b.ProcessFile(buildCtx)
			freeContextDocument(buildCtx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
