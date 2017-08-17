package cmd

import (
	"context"
	"sort"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

type globalsCmd struct {
	cmd *cobra.Command
	o   *globalOptions
	ao  *astOptions
	l   []cprofile.LoaderOptionsFunc
}

func createGlobalsCommand(o *globalOptions) *globalsCmd {
	globalsCmd := &globalsCmd{
		nil,
		o,
		&astOptions{},
		[]cprofile.LoaderOptionsFunc{},
	}

	c := &cobra.Command{
		Use:     "globals",
		Short:   "Returns list of instances of global variables.",
		Long:    "Returns the list of global variables for a program, with file name and offsets.",
		PreRunE: globalsCmd.PreRunE,
		Run:     globalsCmd.Run,
	}

	globalsCmd.cmd = c
	globalsCmd.ao.AttachFlags(c)

	return globalsCmd
}

func (c *globalsCmd) PreRunE(cmd *cobra.Command, _ []string) error {
	c.l = c.ao.ProcessFlags(cmd, c.l)
	return nil
}

func (c *globalsCmd) Run(_ *cobra.Command, args []string) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	base := "."
	if len(args) > 0 {
		base = args[0]
	}

	l := cprofile.NewLoader()
	p, err := l.Load(ctx, base)
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
}
