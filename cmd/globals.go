package cmd

import (
	"context"
	"sort"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

var globalsCmd = &cobra.Command{
	Use:   "globals",
	Short: "Returns list of instances of global variables.",
	Long:  "Returns the list of global variables for a program, with file name and offsets.",
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		if Verbose {
			// Adjusting the log level
			cprofile.Stdout().SetLevel(cprofile.Verbose)
			cprofile.Stderr().SetLevel(cprofile.Verbose)
		}
		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		base := "."
		if len(args) > 0 {
			base = args[0]
		}

		l := cprofile.NewLoader()
		p, err := l.Load(ctx, base, cprofile.Shallow)
		if err != nil {
			cprofile.Stderr().Printf("Got error: %s\n", err.Error())
			return
		}

		pkg, err := p.Package()
		if err != nil {
			cprofile.Stderr().Printf("Got error: %s\n", err.Error())
			return
		}

		globals := pkg.Globals(p.FileSet())
		sort.Strings(globals)

		stdout := cprofile.Stdout()

		for _, v := range globals {
			stdout.Printf("%s\n", v)
		}
	},
}
