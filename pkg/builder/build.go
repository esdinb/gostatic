package builder

import (
    "os"
    "path/filepath"
    "errors"
    "strings"

    "gostatic/pkg/transformer"

    "github.com/jbussdieker/golibxml"
)

func ReadFile(path string) *golibxml.Document {
    return golibxml.ReadFile(
        path,
        "UTF-8",
        golibxml.XML_PARSE_RECOVER &
        golibxml.XML_PARSE_NONET &
        golibxml.XML_PARSE_PEDANTIC &
        golibxml.XML_PARSE_NOBLANKS &
        golibxml.XML_PARSE_XINCLUDE,
    )
}

func WriteFile(doc *golibxml.Document, path string) error {
    err := os.MkdirAll(filepath.Dir(path), 0755)
    if err != nil {
        return err
    }
    length := doc.SaveFile(path)
    if length == -1 {
        return errors.New("failed to write transformation to file")
    } else {
        return nil
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
    In string
    Out string
    Pipeline Pipeline
}

func NewBuildSection(in string, out string, pipeline Pipeline) BuildSection {
    return BuildSection{in, out, pipeline}
}

func (p *Pipeline) Transform(context *transformer.Context) (*transformer.Context, error) {
    var status transformer.Status
    var err error
    for _, command := range *p {
        cmd := command.Parse()
        fn := transformer.Registry.Lookup(cmd.Name)
        if fn == nil {
            continue
        }
        context, status, err = fn(context, cmd.Args)
        if err != nil {
            return context, err
        }
        if status == transformer.Stop {
            break
        } else {
            continue
        }
    }
    return context, nil
}

func (b *BuildSection) ProcessFile(inPath string, outPath string, rootPath string) error {
    var doc *golibxml.Document
    var err error
    if inPath == "-" {
        return errors.New("cannot read from stdin")
    } else {
        doc = ReadFile(inPath)
    }
    defer doc.Free()

    context := &transformer.Context{inPath, outPath, rootPath, doc}

    context, err = b.Pipeline.Transform(context)
    if err != nil {
        return err
    }

    err = WriteFile(context.Document, context.OutPath)
    if err != nil {
        return err
    }
    return nil
}

func (b *BuildSection) Build(rootPath string) error {
    var err error
    absPath := rootPath
    outPath := b.Out
    outPathIsADir := strings.HasSuffix(outPath, string(os.PathSeparator))
    if outPath != "-" {
        absPath, err = filepath.Abs(filepath.Join(rootPath, outPath))
        if err != nil {
            return err
        }

        fileInfo, err := os.Stat(absPath)
        if err != nil {
            if os.IsNotExist(err) {
                newPath := absPath
                if !outPathIsADir {
                    newPath = filepath.Dir(absPath)
                }
                err = os.MkdirAll(newPath, 0750)
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
    }

    var matches []string
    matches, err = filepath.Glob(filepath.Join(rootPath, b.In))
    if err != nil {
        return err
    }

    for i := range matches {
        inPath := matches[i]
        if outPathIsADir {
            relPath, err := filepath.Rel(rootPath, inPath)
            if err != nil {
                return err
            }
            outPath = filepath.Join(absPath, relPath)
        }
        err = b.ProcessFile(inPath, outPath, rootPath)
        if err != nil{
            return err
        }
    }
    return nil
}

