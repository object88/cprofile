package cmd

import (
	"sort"

	"github.com/object88/cprofile"
)

func createGlobalsCommand(o *globalOptions) *astCmd {
	astSetup := &astSetup{
		"globals",
		"Returns list of instances of global variables.",
		"Returns the list of global variables for a program, with file name and offsets.",
		func(p *cprofile.Program) {
			globals := []string{}
			pkgs := p.Imports()

			if len(pkgs) == 0 {
				return
			}

			for _, pkg := range pkgs {
				gs := pkg.Globals(p.FileSet())
				for _, v := range gs {
					globals = append(globals, v)
				}
			}

			sort.Strings(globals)
			stdout := cprofile.Stdout()

			for _, v := range globals {
				stdout.Printf("%s\n", v)
			}
		},
		nil,
	}

	globalsCmd := createAstCommand(o, astSetup)
	return globalsCmd
}
