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
	"runtime"
	"strings"
)

const (
	kPath = "./"
)

// Imports contains the import structure
type Imports struct {
	packages []string
}

type importReadData struct {
	// cwd      string
	// importer types.ImporterFrom
	packages map[string]bool
}

// NewImports creates a new Imports struct
func NewImports() *Imports {
	return &Imports{}
}

// Read loads the import chain.
func (i *Imports) Read(ctx context.Context, base string) error {
	if i == nil {
		return errors.New("No pointer receiver")
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, base, nil, parser.DeclarationErrors)
	if err != nil {
		fmt.Printf("Failed to do AST parsing: %s\n", err.Error())
		return err
	}

	astf := make([]*ast.File, 0)
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			astf = append(astf, f)
		}
	}

	config := &types.Config{
		Error: func(e error) {
			fmt.Println(e)
		},
		Importer: importer.Default(),
	}
	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	ird := &importReadData{
		// cwd:      trueCwd,
		// importer: impf,
		packages: map[string]bool{},
	}
	pkgfoo, err := config.Check(kPath, fset, astf, &info)
	if err != nil {
		return err
	}
	fmt.Printf("types.Config.Check got %v\n", pkgfoo.String())
	for _, v := range pkgfoo.Imports() {
		i.read(ctx, ird, v, 1)
		// fmt.Printf("%d -> %#v\n", k, v)
		// for k, v := range v.Imports() {
		// 	fmt.Printf("  %d -> %#v\n", k, v)
		// }
	}

	fmt.Printf("Go root: %s\n", runtime.GOROOT())
	// fmt.Printf("Tool dir: %s\n", build.ToolDir)
	// fmt.Printf("Default:\n%#v\n", build.Default)

	// cwd, err := os.Getwd()
	// if err != nil {
	// 	// Failed to get the current working directory.  Wat.
	// 	return err
	// }
	// trueCwd, err := filepath.EvalSymlinks(cwd)
	// if err != nil {
	// 	// Failed to eval any symlinks.  Wat.
	// 	return err
	// }

	// impf, ok := importer.Default().(types.ImporterFrom)
	// if !ok {
	// 	return errors.New("Crap")
	// }
	// select {
	// case <-ctx.Done():
	// 	return context.Canceled
	// default:

	// 	p, err := ird.importer.ImportFrom(base, ird.cwd, 0)
	// 	if err != nil {
	// 		fmt.Printf("Initial import failed with base '%s'\n", base)
	// 		return err
	// 	}

	// 	for _, v := range p.Imports() {
	// 		err := i.read(ctx, ird, v, 0)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// result := make([]string, len(ird.packages))
	// idx := 0
	// for k := range ird.packages {
	// 	result[idx] = k
	// 	idx++
	// }

	// sort.Strings(result)
	// i.packages = result

	return nil
}

func (i *Imports) read(ctx context.Context, ird *importReadData, pkg *types.Package, depth int) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}

	if !pkg.Complete() {
		fmt.Printf("Package %s is incomplete!\n", pkg.Name())
	}

	fmt.Printf("%s>> %s (is local: %t)\n", strings.Repeat("  ", depth), pkg.Path(), build.IsLocalImport(pkg.Name()))

	for _, v := range pkg.Imports() {
		// Ignore the "C" pseudo package and packages that we have
		// already seen.
		name := v.Path()
		if name == "C" {
			continue
		}
		if _, ok := ird.packages[name]; ok {
			fmt.Printf("Already have %s\n", name)
			continue
		}

		// Dig deeper.
		ird.packages[name] = true
		err := i.read(ctx, ird, v, depth+1)
		if err != nil {
			if iErr, ok := err.(ImportError); ok {
				// iErr.AddPackage(p)
				return iErr
			}
			// fmt.Printf("***\n%#v\n", p.Imports)
			// return NewImportError(v, p, err)
		}
	}

	return nil
}

// Flatlist returns a flat, alphabetically sorted list of package names
func (i *Imports) Flatlist() []string {
	if i == nil {
		return []string{}
	}

	return i.packages
}
