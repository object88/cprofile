package cprofile

import (
	"fmt"
	"strings"
)

// AstDepth directs the classifications of packages inspected
type AstDepth int

const (
	// Shallow addresses the contents of the package
	Shallow AstDepth = iota

	// Deep addresses the contents of the given package and all its referenced children
	Deep

	// Local addresses the contents of all packages in the same organization
	Local

	// Wide addresses the contents of all directly and indirectly imported non-stdlib packages
	Wide

	// Complete addresses the contents of every imported package, including stdlib
	Complete
)

// CheckAstDepth validates the AstDepth flag
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
