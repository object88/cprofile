package cprofile

import (
	"fmt"
	"go/build"
	"go/token"
	"strings"
)

type loaderState struct {
	base    string
	org     string
	orgPath string
	fset    *token.FileSet
	pkgs    map[string]*Package
	spacers *[]string
	options *LoaderOptions
}

func newLoaderState(pkgName string, options *LoaderOptions) *loaderState {
	maxDepth := options.Depth
	if maxDepth == DefaultDepth {
		maxDepth = 8
	}
	spacers := make([]string, maxDepth)
	for i := 0; i < maxDepth; i++ {
		spacers[i] = strings.Repeat("  ", i)
	}

	org := ""
	orgPath := ""
	s := strings.Split(pkgName, "/")
	if strings.Index(s[0], ".") != -1 {
		org = fmt.Sprintf("%s/%s", s[0], s[1])
		orgPath = fmt.Sprintf("%s/src/%s/%s", build.Default.GOPATH, s[0], s[1])
	}

	ls := &loaderState{
		pkgName,
		org,
		orgPath,
		token.NewFileSet(),
		map[string]*Package{},
		&spacers,
		options,
	}

	return ls
}

func (ls *loaderState) getSpacer(depth int) string {
	initialSize := len(*ls.spacers)
	if depth >= initialSize {
		// Expand the contents
		targetSize := nextPowerOfTwo(depth + 1)
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
