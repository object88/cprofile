package cprofile

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

// Package represents a Go package
type Package struct {
	asts *[]*ast.File
	info *types.Info
	name string
	pkg  *types.Package
}

func newPkg(name string) *Package {
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	return &Package{nil, info, name, nil}
}

// Globals returns a list of global variables in the package
func (p *Package) Globals(fset *token.FileSet) []string {
	if p == nil {
		return []string{}
	}

	globals := map[string]string{}
	for k, def := range p.info.Defs {
		if def == nil {
			continue
		}

		switch v := def.(type) {
		case *types.Var:
			parent := v.Parent()
			if parent == nil {
				continue
			}
			grandparent := parent.Parent()
			if grandparent != types.Universe {
				continue
			}
			globals[k.Name] = fset.Position(k.Pos()).String()
		default:
			continue
		}
	}

	results := make([]string, len(globals))
	i := 0
	for k, v := range globals {
		results[i] = fmt.Sprintf("%s: %s", k, v)
		i++
	}

	return results
}

// Imports returns the list of packages imported by this package
func (p *Package) Imports() ([]*types.Object, error) {
	if p == nil {
		return []*types.Object{}, nil
	}

	results := make([]*types.Object, len(p.info.Uses))
	i := 0
	for _, v := range p.info.Uses {
		results[i] = &v
		i++
	}
	return results, nil
}

// Name returns the name of this package.
func (p *Package) Name() string {
	return p.name
}
