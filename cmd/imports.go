package cmd

import (
	"context"
	"sort"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

type importsCmd struct {
	cmd *cobra.Command
	o   *globalOptions
	ao  *astOptions
	l   []cprofile.LoaderOptionsFunc
}

func createImportsCommand(o *globalOptions) *importsCmd {
	importsCmd := &importsCmd{
		nil,
		o,
		&astOptions{},
		[]cprofile.LoaderOptionsFunc{},
	}

	c := &cobra.Command{
		Use:     "imports",
		Short:   "Print the imports.",
		Long:    "Gets the de-duplicated list of imports.",
		PreRunE: importsCmd.PreRunE,
		Run:     importsCmd.Run,
	}

	importsCmd.cmd = c
	importsCmd.ao.AttachFlags(c)

	return importsCmd
}

func (c *importsCmd) PreRunE(cmd *cobra.Command, args []string) error {
	c.l = c.ao.ProcessFlags(cmd, c.l)
	return nil
}

func (c *importsCmd) Run(cmd *cobra.Command, args []string) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	base := "."
	if len(args) > 0 {
		base = args[0]
	}

	l := cprofile.NewLoader()
	p, err := l.Load(ctx, base, c.l...)
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
}
