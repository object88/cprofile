package cmd

import (
	"context"
	"sort"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

var astD cprofile.AstDepth

var astDepth string

var importsCmd = &cobra.Command{
	Use:   "imports",
	Short: "Print the imports.",
	Long:  "Gets the de-duplicated list of imports.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if Verbose {
			// Adjusting the log level
			cprofile.Stdout().SetLevel(cprofile.Verbose)
			cprofile.Stderr().SetLevel(cprofile.Verbose)
		}

		d := cmd.Flag("astDepth").Value
		a, err := cprofile.CheckAstDepth(d.String())
		if err != nil {
			return err
		}
		astD = a
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		base := "."
		if len(args) > 0 {
			base = args[0]
		}

		l := cprofile.NewLoader()
		p, err := l.Load(ctx, base, astD)
		if err != nil {
			cprofile.Stderr().Printf("Got error: %s\n", err.Error())
			return
		}

		pkgs := p.Imports()
		if len(pkgs) == 0 {
			return
		}

		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].Name() < pkgs[j].Name()
		})

		stdout := cprofile.Stdout()
		for _, v := range pkgs {
			stdout.Printf("%s\n", v.Name())
		}
	},
}
