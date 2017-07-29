package cprofile

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Program is not used
type Program struct {
}

// Loader will load code into an AST
type Loader struct {
	config  *types.Config
	pkgs    map[string]*pkg
	ps      map[string]*types.Package
	srcDirs []string
}

type pkg struct {
	asts     *[]*ast.File
	complete bool
	fset     *token.FileSet
	info     *types.Info
}

// NewLoader constructs a new Loader struct
func NewLoader() *Loader {
	config := &types.Config{
		Error: func(e error) {
			fmt.Println(e)
		},
		Importer: importer.Default(),
	}
	srcDirs := build.Default.SrcDirs()
	return &Loader{config, map[string]*pkg{}, map[string]*types.Package{}, srcDirs}
}

func newPkg() *pkg {
	fset := token.NewFileSet()
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	return &pkg{nil, false, fset, info}
}

// Load reads in the AST
func (l *Loader) Load(ctx context.Context, base string) (*Program, error) {
	if l == nil {
		return nil, errors.New("No pointer reciever")
	}

	abs, err := filepath.Abs(base)
	if err != nil {
		return nil, err
	}

	pkgName := ""
	for _, v := range l.srcDirs {
		if strings.HasPrefix(abs, v) {
			pkgName = abs[len(v)+1:]
		}
	}
	if pkgName == "" {
		msg := fmt.Sprintf("Failed to find '%s'", base)
		return nil, errors.New(msg)
	}

	fmt.Printf("pkgName: '%s'\n", pkgName)

	return l.load(ctx, pkgName, 0)
}

func (l *Loader) load(ctx context.Context, base string, depth int) (*Program, error) {
	spacer := strings.Repeat("  ", depth)
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	fpath, err := l.findSourcePath(base)
	if err != nil {
		return nil, err
	}

	buildP, err := build.ImportDir(fpath, 0)
	if err != nil {
		return nil, err
	}

	pkg := newPkg()
	astPkgs := make(map[string]*ast.Package)
	for _, v := range buildP.GoFiles {
		fpath := path.Join(buildP.Dir, v)

		astf, err := parser.ParseFile(pkg.fset, fpath, nil, 0)
		if err != nil {
			return nil, err
		}

		name := astf.Name.Name
		pkg, found := astPkgs[name]
		if !found {
			pkg = &ast.Package{
				Name:  name,
				Files: make(map[string]*ast.File),
			}
			astPkgs[name] = pkg
		}
		pkg.Files[fpath] = astf
	}

	for k, v := range astPkgs {
		if !strings.HasSuffix(fpath, k) {
			fmt.Printf("fpath = '%s'\n", fpath)
			fmt.Printf("Skipping '%s'\n", k)
			continue
		}

		pkg.asts = makeAsts(v)

		p, err := l.config.Check(k, pkg.fset, *pkg.asts, pkg.info)
		if err != nil {
			return nil, err
		}
		id := base
		fmt.Printf("%s-- Adding key '%s'\n", spacer, id)
		l.pkgs[id] = pkg
		l.ps[id] = p
	}

	p, ok := l.ps[base]
	if !ok {
		return nil, errors.New("Failed to find package with directory-eponymous name")
	}
	fmt.Printf("%sProcessing '%s' imports...\n", spacer, p.Path())

	imports := p.Imports()
	for _, v0 := range imports {
		id := v0.Path()
		fmt.Printf("%s** Checking for '%s'...", spacer, id)
		if _, ok := l.pkgs[id]; ok {
			fmt.Printf(" already parsed; skipping...\n")
			continue
		}

		fmt.Printf(" processing...\n")
		_, err := l.load(ctx, id, depth+1)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (l *Loader) findSourcePath(pkgName string) (string, error) {
	if pkgName == "." {
		p, err := os.Getwd()
		if err != nil {
			return "", err
		}
		fmt.Printf("Got '.'; using '%s'\n", p)

		return p, nil
	}

	for _, v := range l.srcDirs {
		fpath := path.Join(v, pkgName)
		isDir := false
		if build.Default.IsDir != nil {
			isDir = build.Default.IsDir(fpath)
		} else {
			s, err := os.Stat(fpath)
			if err != nil {
				continue
			}
			isDir = s.IsDir()
		}
		if isDir {
			// fmt.Printf("Have source path '%s'\n", fpath)
			return fpath, nil
		}
	}

	msg := fmt.Sprintf("Failed to locate package '%s'", pkgName)
	return "", errors.New(msg)
}

func makeAsts(pkg *ast.Package) *[]*ast.File {
	asts := make([]*ast.File, len(pkg.Files))
	i := 0
	for _, f := range pkg.Files {
		asts[i] = f
		i++
	}

	return &asts
}
