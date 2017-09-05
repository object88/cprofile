package cmd

import (
	"sort"

	"github.com/object88/cprofile"
	"github.com/spf13/pflag"
)

const (
	Universal int = iota
)

func createFunctionsCommand(o *globalOptions) *astCmd {
	var scope string

	astSetup := &astSetup{
		"functions",
		"Returns list of public functions.",
		"Returns the list of public functions in the desired scope.",
		func(p *cprofile.Program) {
			functions := []string{}
			pkgs := p.Imports()

			if len(pkgs) == 0 {
				return
			}

			for _, pkg := range pkgs {
				fns := pkg.Functions(p.FileSet(), scope)
				for _, v := range fns {
					functions = append(functions, v)
				}
			}

			sort.Strings(functions)
			stdout := cprofile.Stdout()

			for _, v := range functions {
				stdout.Printf("%s\n", v)
			}
		},
		[]func(fs *pflag.FlagSet){
			func(fs *pflag.FlagSet) {
				fs.StringVarP(&scope, "scope", "s", "u", "Scope")
			},
		},
	}

	functionsCmd := createAstCommand(o, astSetup)
	return functionsCmd
}
