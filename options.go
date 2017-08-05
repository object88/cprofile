package cprofile

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
