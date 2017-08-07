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
	"runtime"
	"strings"
)

// Loader will load code into an AST
type Loader struct {
	config  *types.Config
	srcDirs []string
	stderr  *Log
}

type loaderState struct {
	base  string
	owner string
	fset  *token.FileSet
	pkgs  map[string]*Package
	ps    map[string]*types.Package
	depth AstDepth
}

// NewLoader constructs a new Loader struct
func NewLoader() *Loader {
	l := Stderr()
	config := &types.Config{
		Error: func(e error) {
			l.Warnf("%s\n", e.Error())
		},
		Importer: importer.Default(),
	}
	srcDirs := build.Default.SrcDirs()
	return &Loader{config, srcDirs, l}
}

// Load reads in the AST
func (l *Loader) Load(ctx context.Context, base string, depth AstDepth) (*Program, error) {
	if l == nil {
		return nil, errors.New("No pointer receiver")
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
		return nil, fmt.Errorf("Failed to find '%s'", base)
	}

	l.stderr.Verbosef("pkgName: '%s'\n", pkgName)

	// [GOPATH]/src/github.com/object88/pkgA0/pkgA1/pkgA2
	fmt.Printf("Source dirs:\n")
	for _, v := range l.srcDirs {
		fmt.Printf("%s\n", v)
	}

	fmt.Printf("Pkgname: %s\n", pkgName)

	ls := &loaderState{pkgName, "", token.NewFileSet(), map[string]*Package{}, map[string]*types.Package{}, depth}

	pkg, err := l.load(ctx, ls, abs, pkgName, 0)
	if err != nil {
		return nil, err
	}

	program := newProgram(ls.fset, ls.pkgs, pkg)

	return program, nil
}

func (l *Loader) load(ctx context.Context, ls *loaderState, fpath, base string, depth int) (*Package, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	spacer := strings.Repeat("  ", depth)

	buildP, err := build.ImportDir(fpath, 0)
	if err != nil {
		return nil, err
	}

	pkg := newPkg(base)
	astPkgs := make(map[string]*ast.Package)
	for _, v := range buildP.GoFiles {
		fpath := path.Join(buildP.Dir, v)

		astf, err := parser.ParseFile(ls.fset, fpath, nil, 0)
		if err != nil {
			l.stderr.Verbosef("Got error while parsing file '%s':\n%s\n", fpath, err.Error())
		}

		name := astf.Name.Name
		astPkg, found := astPkgs[name]
		if !found {
			astPkg = &ast.Package{
				Name:  name,
				Files: make(map[string]*ast.File),
			}
			astPkgs[name] = astPkg
		}
		astPkg.Files[fpath] = astf
	}

	for k, v := range astPkgs {
		if !strings.HasSuffix(fpath, k) {
			l.stderr.Verbosef("fpath = '%s'\n", fpath)
			l.stderr.Verbosef("Skipping '%s'\n", k)
			continue
		}

		pkg.asts = makeAsts(v)

		p, err := l.config.Check(k, ls.fset, *pkg.asts, pkg.info)
		if err != nil {
			l.stderr.Verbosef("Got error checking package '%s':\n%s\n", k, err.Error())
			// return err
		}
		id := base
		path, err := l.findSourcePath(base)
		if err != nil {
			return nil, err
		}
		l.stderr.Verbosef("%s-- Adding key '%s' / '%s'\n", spacer, id, path)
		ls.pkgs[path] = pkg
		ls.ps[path] = p
	}

	p, ok := ls.ps[fpath]
	if !ok {
		return nil, errors.New("Failed to find package with directory-eponymous name")
	}
	l.stderr.Verbosef("%sProcessing '%s' imports...\n", spacer, p.Path())

	imports := p.Imports()
	for _, v0 := range imports {
		id := v0.Path()
		path, err := l.findSourcePath(id)
		if err != nil {
			return nil, err
		}
		l.stderr.Verbosef("%s** Checking for '%s' / '%s'...", spacer, id, path)
		if _, ok := ls.pkgs[path]; ok {
			l.stderr.Verbosef(" already parsed; skipping...\n")
			continue
		}

		if !l.checkDepth(ls, id) {
			l.stderr.Verbosef(" failed depth check...\n")
			continue
		}

		l.stderr.Verbosef(" processing...\n")
		_, err = l.load(ctx, ls, path, id, depth+1)
		if err != nil {
			return nil, err
		}
	}

	return pkg, nil
}

func (l *Loader) findSourcePath(pkgName string) (string, error) {
	if pkgName == "." {
		p, err := os.Getwd()
		if err != nil {
			return "", err
		}
		l.stderr.Verbosef("Got '.'; using '%s'\n", p)

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
			return fpath, nil
		}
	}

	return "", fmt.Errorf("Failed to locate package '%s'", pkgName)
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

func (l *Loader) checkDepth(ls *loaderState, pkgName string) bool {
	if ls.depth == Complete {
		// Complete includes everything
		return true
	}

	p, err := l.findSourcePath(pkgName)
	if err != nil {
		return false
	}

	fmt.Printf(" -> %s", p)

	if ls.depth == Shallow {
		// Only my direct imports
		return ls.base == pkgName
	} else if ls.depth == Deep {
		return strings.HasPrefix(p, ls.base)
	} else if ls.depth == Local {

	}

	// Wide: everything that isn't in stdlib.
	return strings.HasPrefix(p, runtime.GOROOT())
}
