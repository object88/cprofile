package cmd

import (
	"sort"

	"github.com/object88/cprofile"
)

func createImportsCommand(o *globalOptions) *astCmd {
	astSetup := &astSetup{
		"imports",
		"Print the imports.",
		"Gets the de-duplicated list of imports.",
		func(p *cprofile.Program) {
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

	importsCmd := createAstCommand(o, astSetup)
	return importsCmd
}
