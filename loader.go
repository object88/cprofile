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
	base         string
	organization string
	fset         *token.FileSet
	pkgs         map[string]*Package
	depth        AstDepth
	spacers      *[]string
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

func newLoaderState(pkgName string, depth AstDepth, maxDepth int) *loaderState {
	if maxDepth == -1 {
		maxDepth = 8
	}
	spacers := make([]string, maxDepth)
	for i := 0; i < maxDepth; i++ {
		spacers[i] = strings.Repeat("  ", i)
	}

	organization := ""
	s := strings.Split(pkgName, "/")
	if strings.Index(s[0], ".") != -1 {
		organization = fmt.Sprintf("%s/src/%s/%s", build.Default.GOPATH, s[0], s[1])
	}

	ls := &loaderState{
		pkgName,
		organization,
		token.NewFileSet(),
		map[string]*Package{},
		depth,
		&spacers,
	}

	return ls
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

	fi, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("Provided path '%s' must be a directory", base)
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
	l.stderr.Debugf("Source dirs:\n")
	for _, v := range l.srcDirs {
		l.stderr.Debugf("%s\n", v)
	}

	l.stderr.Debugf("Pkgname: %s\n", pkgName)

	ls := newLoaderState(pkgName, depth, 8)

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

	pkg, p, err := l.buildPackage(ls, fpath, base, depth)
	if err != nil {
		return nil, err
	}

	l.stderr.Verbosef("%sProcessing '%s' imports...\n", ls.getSpacer(depth), p.Path())

	err = l.visitImports(ctx, p, ls, depth)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func (l *Loader) buildAstPackages(buildP *build.Package, ls *loaderState) map[string]*ast.Package {
	astPkgs := map[string]*ast.Package{}

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

	return astPkgs
}

func (l *Loader) buildPackage(ls *loaderState, fpath, base string, depth int) (*Package, *types.Package, error) {
	buildP, err := build.ImportDir(fpath, 0)
	if err != nil {
		return nil, nil, err
	}

	pkg := newPkg(base)
	astPkgs := l.buildAstPackages(buildP, ls)

	for k, v := range astPkgs {
		if strings.HasSuffix(v.Name, "_test") {
			continue
		}

		pkg.asts = makeAsts(v)

		p, err := l.config.Check(k, ls.fset, *pkg.asts, pkg.info)
		if err != nil {
			l.stderr.Verbosef("Got error checking package '%s':\n%s\n", k, err.Error())
		}

		path, err := l.findSourcePath(base)
		if err != nil {
			return nil, nil, err
		}
		l.stderr.Verbosef("%s-- Adding key '%s' / '%s'\n", ls.getSpacer(depth), base, path)

		pkg.pkg = p
		ls.pkgs[path] = pkg
	}

	pkg, ok := ls.pkgs[fpath]
	if !ok {
		return nil, nil, errors.New("Failed to find package with directory-eponymous name")
	}

	return pkg, pkg.pkg, nil
}

func (l *Loader) visitImports(ctx context.Context, p *types.Package, ls *loaderState, depth int) error {
	imports := p.Imports()
	for _, v0 := range imports {
		id := v0.Path()
		path, err := l.findSourcePath(id)
		if err != nil {
			return err
		}

		if _, ok := ls.pkgs[path]; ok {
			l.stderr.Verbosef("%s** Checking for '%s' / '%s' already parsed; skipping...\n", ls.getSpacer(depth), id, path)
			continue
		}

		if !l.checkDepth(ls, id) {
			l.stderr.Verbosef("%s** Checking for '%s' / '%s' failed depth check.\n", ls.getSpacer(depth), id, path)
			continue
		}

		l.stderr.Verbosef("%s** Checking for '%s' / '%s' processing...\n", ls.getSpacer(depth), id, path)
		_, err = l.load(ctx, ls, path, id, depth+1)
		if err != nil {
			return err
		}
	}

	return nil
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

	if ls.depth == Shallow {
		// Only my direct imports
		return ls.base == pkgName
	} else if ls.depth == Deep {
		return strings.HasPrefix(p, ls.base)
	} else if ls.depth == Local {
		return strings.HasPrefix(p, ls.organization)
	}

	// Wide: everything that isn't in stdlib.
	return !strings.HasPrefix(p, runtime.GOROOT())
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

func (ls *loaderState) getSpacer(depth int) string {
	initialSize := len(*ls.spacers)
	if depth >= initialSize {
		// Expand the contents
		targetSize := nextPowerOfTwo(depth)
		target := make([]string, targetSize)
		for i := 0; i < initialSize; i++ {
			target[i] = (*ls.spacers)[i]
		}
		for i := initialSize; i < targetSize; i++ {
			target[i] = strings.Repeat("  ", i)
		}

		ls.spacers = &target
	}

	return (*ls.spacers)[depth]
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
