package cprofile

import (
	"bytes"
	"go/build"
)

// ImportError is an error that comes from reading the import chain
type ImportError interface {
	error
	AddPackage(*build.Package)
}

type importError struct {
	base        error
	packageName string
	packages    []*build.Package
}

// NewImportError returns a new import error
func NewImportError(packageName string, currentPackage *build.Package, base error) ImportError {
	return &importError{base, packageName, []*build.Package{currentPackage}}
}

func (ie *importError) AddPackage(p *build.Package) {
	ie.packages = append(ie.packages, p)
}

// Error returns the error string
func (ie *importError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString(ie.base.Error())
	buffer.WriteString("\n")
	buffer.WriteString(ie.packageName)
	for _, v := range ie.packages {
		buffer.WriteString("\n")
		buffer.WriteString(v.Name)
	}
	return buffer.String()
}
