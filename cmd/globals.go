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

	globalsCmd := createAstCommand(o, astSetup)
	return globalsCmd
}
