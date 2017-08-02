package cprofile

import (
	"go/ast"
	"go/token"
	"go/types"
)

// Package represents a Go package
type Package struct {
	asts *[]*ast.File
	fset *token.FileSet
	info *types.Info
	name string
}

func newPkg(name string) *Package {
	fset := token.NewFileSet()
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	return &Package{nil, fset, info, name}
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
