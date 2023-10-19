package builder

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	inPath, ok := ctx.Value(transformer.InPathContextKey).(string)
	if !ok {
		return nil, errors.New("missing input path for xml loader")
	}

	doc := markup.ReadFile(inPath, "UTF-8", parseOptions)
	if doc == nil {
		return nil, errors.New("unable to load xml file")
	}

	return context.WithValue(ctx, transformer.DocumentContextKey, doc), nil
}

func htmlLoader(ctx context.Context) (context.Context, error) {
	inPath, ok := ctx.Value(transformer.InPathContextKey).(string)
	if !ok {
		return nil, errors.New("missing input path for html loader")
	}

	doc := markup.ReadHTMLFile(inPath, parseOptions)
	if doc == nil {
		return nil, errors.New("unable to load xml file")
	}

	return context.WithValue(ctx, transformer.DocumentContextKey, doc), nil
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

	doc, ok = ctx.Value(transformer.DocumentContextKey).(*markup.Document)
	if !ok {
		return errors.New("missing transformation document for xml formatter")
	}

	outFile, ok = ctx.Value(transformer.OutFileContextKey).(*os.File)
	if !ok {
		outPath, ok := ctx.Value(transformer.OutPathContextKey).(string)
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
	)

	options = markup.SaveOption(markup.XML_SAVE_NO_DECL | markup.XML_SAVE_NO_EMPTY | markup.XML_SAVE_AS_XML)

	doc, ok = ctx.Value(transformer.DocumentContextKey).(*markup.Document)
	if !ok {
		return errors.New("missing transformation document for html formatter")
	}

	outFile, ok = ctx.Value(transformer.OutFileContextKey).(*os.File)
	if !ok {
		outPath, ok := ctx.Value(transformer.OutPathContextKey).(string)
		if !ok {
			return errors.New("missing output file path for html formatter")
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

	inFile, ok = ctx.Value(transformer.InFileContextKey).(*os.File)
	if !ok {
		inPath, ok := ctx.Value(transformer.InPathContextKey).(string)
		if !ok {
			return errors.New("missing input file for copy")
		}
		inFile, err = openFileForReading(inPath)
		if err != nil {
			return errors.New("cannot open input file for copy")
		}
	}

	outFile, ok = ctx.Value(transformer.OutFileContextKey).(*os.File)
	if !ok {
		outPath, ok := ctx.Value(transformer.OutPathContextKey).(string)
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

func freeBuildContext(ctx context.Context) {
	var (
		doc     *markup.Document
		inFile  *os.File
		outFile *os.File
		ok      bool
	)

	doc, ok = ctx.Value(transformer.DocumentContextKey).(*markup.Document)
	if ok {
		doc.Free()
	}

	inFile, ok = ctx.Value(transformer.InFileContextKey).(*os.File)
	if ok {
		inFile.Close()
	}

	outFile, ok = ctx.Value(transformer.OutFileContextKey).(*os.File)
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

func (b *BuildSection) ProcessFile(ctx context.Context) error {

	ctx, err := b.Pipeline.Transform(ctx)
	if err != nil {
		return err
	}

	if format, ok := ctx.Value(transformer.FormatterContextKey).(FormatterFunc); ok {
		err := format(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BuildSection) Build(ctx context.Context, rootPath string) error {
	var (
		outFile *os.File = nil
		err     error
	)
	rootPathAbsolute, err := filepath.Abs(rootPath)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, transformer.InPathContextKey, b.In)
	ctx = context.WithValue(ctx, transformer.OutPathContextKey, b.Out)
	ctx = context.WithValue(ctx, transformer.RootPathContextKey, rootPath)
	ctx = context.WithValue(ctx, transformer.ParamsContextKey, []string{})
	ctx = context.WithValue(ctx, transformer.StringParamsContextKey, []string{})

	inPath := b.In
	outPath := b.Out
	outPathIsDir := strings.HasSuffix(outPath, string(os.PathSeparator))
	absPath := filepath.Join(rootPathAbsolute, outPath)

	if outPath == "-" {
		outFile = os.Stdout
		ctx = context.WithValue(ctx, transformer.OutFileContextKey, os.Stdout)
		ctx = context.WithValue(ctx, transformer.FormatterContextKey, lookupFormatter(filepath.Ext(inPath)))
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
		formatter := lookupFormatter(filepath.Ext(outPath))
		ctx = context.WithValue(ctx, transformer.FormatterContextKey, formatter)

		outPath = absPath
	}

	if b.In == "-" {
		var (
			reader   *bufio.Reader
			readNext bool
			bytes    []byte
			err      error
		)
		reader = bufio.NewReader(os.Stdin)
		readNext = true
		for readNext {
			buildCtx := ctx
			bytes, err = reader.ReadBytes(0)
			if err != nil {
				if errors.Is(err, io.EOF) {
					readNext = false
				} else {
					return err
				}
			}
			if len(bytes) > 0 {
				if doc, ok := buildCtx.Value(transformer.DocumentContextKey).(*markup.Document); ok {
					doc.Free()
					if outFile != nil {
						outFile.Write([]byte{0})
					}
				}
				buildCtx = context.WithValue(
					buildCtx,
					transformer.DocumentContextKey,
					markup.ReadMemory(bytes, b.In, "UTF-8", parseOptions),
				)
				err = b.ProcessFile(buildCtx)
				if err != nil {
					return err
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
			if outPathIsDir {
				if relPath, err := filepath.Rel(rootPath, inPath); err != nil {
					return err
				} else {
					outPath = filepath.Join(absPath, relPath)
				}
			}
			ctx = context.WithValue(ctx, transformer.InPathContextKey, inPath)
			ctx = context.WithValue(ctx, transformer.OutPathContextKey, outPath)
			loader := lookupLoader(filepath.Ext(inPath))
			ctx, err = loader(ctx)
			if err != nil {
				return err
			}

			err = b.ProcessFile(ctx)
			freeBuildContext(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
