package cmd

import (
	"context"
	"sort"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

var astDepth string

var importsCmd = &cobra.Command{
	Use:   "imports",
	Short: "Print the imports.",
	Long:  "Gets the de-duplicated list of imports.",
	Run: func(cmd *cobra.Command, args []string) {
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
