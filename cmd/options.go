package cmd

import (
	"strconv"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var optionFuncs []cprofile.LoaderOptionsFunc

type globalOptions struct {
	verbose bool
}

func (o *globalOptions) AttachFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&o.verbose, "verbose", "v", false, "verbose output")
}

func (o *globalOptions) ProcessFlags(cmd *cobra.Command, funcs []cprofile.LoaderOptionsFunc) {
	if o.verbose {
		// Adjusting the log level
		cprofile.Stdout().SetLevel(cprofile.Verbose)
		cprofile.Stderr().SetLevel(cprofile.Verbose)
	}
}

type astOptions struct {
	astDepth string
	depth    int
}

func (ao *astOptions) AttachFlags(cmd *cobra.Command, flagAdders ...func(fs *pflag.FlagSet)) {
	fs := cmd.Flags()
	fs.StringVarP(&ao.astDepth, "astDepth", "a", "s", "AST depth")
	fs.IntVarP(&ao.depth, "depth", "d", cprofile.DefaultDepth, "Depth")
	for _, v := range flagAdders {
		v(fs)
	}
}

func (ao *astOptions) ProcessFlags(cmd *cobra.Command, funcs []cprofile.LoaderOptionsFunc) []cprofile.LoaderOptionsFunc {
	if cmd.Flag("astDepth").Changed {
		d := cmd.Flag("astDepth").Value
		funcs = append(funcs, assignAstDepth(d.String()))
	}
	if cmd.Flag("depth").Changed {
		d := cmd.Flag("depth").Value
		funcs = append(funcs, assignDepth(d.String()))
	}

	return funcs
}

func assignAstDepth(astDepth string) cprofile.LoaderOptionsFunc {
	return func(lo *cprofile.LoaderOptions) (*cprofile.LoaderOptions, error) {
		d, err := cprofile.CheckAstDepth(astDepth)
		if err != nil {
			return lo, err
		}
		lo.AstDepth = d
		return lo, nil
	}
}

func assignDepth(depth string) cprofile.LoaderOptionsFunc {
	return func(lo *cprofile.LoaderOptions) (*cprofile.LoaderOptions, error) {
		d, err := strconv.Atoi(depth)
		if err != nil {
			return lo, err
		}
		lo.Depth = d
		return lo, nil
	}
}
