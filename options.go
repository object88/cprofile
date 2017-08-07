package cprofile

import (
	"fmt"
	"strings"
)

// AstDepth directs the classifications of packages inspected
type AstDepth int

const (
	// Shallow only inspects the package in the given directory
	Shallow AstDepth = iota

	// Deep is the given package and all its children
	Deep

	// Local is all packages by the same owner
	Local

	// Wide is all non-stdlib packages
	Wide

	// Complete is every package, including stdlib
	Complete
)

func CheckAstDepth(v string) (AstDepth, error) {
	if len(v) == 1 {
		switch rune(strings.ToLower(v)[0]) {
		case 's':
			return Shallow, nil
		case 'd':
			return Deep, nil
		case 'l':
			return Local, nil
		case 'w':
			return Wide, nil
		case 'c':
			return Complete, nil
		}
	} else {

	}

	return 0, fmt.Errorf("Unknown ast depth value '%s'", v)
}
