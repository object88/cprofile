package cprofile

const (
	// DefaultAstDepth specifies that the AST loader will only inspect the specified
	// package
	DefaultAstDepth AstDepth = Shallow

	// DefaultDepth specifies that there is no limit to the number of imports that
	// the AST loader will delve into
	DefaultDepth int = -1
)

// LoaderOptions are the options used by the Loader to limit package transversal
type LoaderOptions struct {
	AstDepth AstDepth
	Depth    int
}

// LoaderOptionsFunc allows functions to determine the state of the LoaderOptions
type LoaderOptionsFunc func(lo *LoaderOptions) (*LoaderOptions, error)

// NewLoaderOptions creates a new LoaderOptions struct and preconfigures it using
// any provided funcs.
func NewLoaderOptions(l *Log, funcs ...LoaderOptionsFunc) *LoaderOptions {
	lo := &LoaderOptions{
		AstDepth: DefaultAstDepth,
		Depth:    DefaultDepth,
	}

	var err error
	for _, lof := range funcs {
		lo, err = lof(lo)
		if err != nil {
			l.Warnf("Failed to apply option: %s", err.Error())
		}
	}

	return lo
}
