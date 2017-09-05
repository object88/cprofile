package cmd

import (
	"context"

	"github.com/spf13/pflag"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

type astCmd struct {
	cmd *cobra.Command
	o   *globalOptions
	ao  *astOptions
	r   func(p *cprofile.Program)
	l   []cprofile.LoaderOptionsFunc
}

type astSetup struct {
	useText   string
	shortText string
	longText  string
	runner    func(p *cprofile.Program)
	flags     []func(fs *pflag.FlagSet)
}

func createAstCommand(o *globalOptions, setup *astSetup) *astCmd {
	astCmd := &astCmd{
		nil,
		o,
		&astOptions{},
		setup.runner,
		[]cprofile.LoaderOptionsFunc{},
	}

	c := &cobra.Command{
		Use:    setup.useText,
		Short:  setup.shortText,
		Long:   setup.longText,
		PreRun: astCmd.PreRun,
		Run:    astCmd.Run,
	}

	astCmd.cmd = c
	astCmd.ao.AttachFlags(c, setup.flags...)

	return astCmd
}

func (c *astCmd) PreRun(cmd *cobra.Command, _ []string) {
	c.l = c.ao.ProcessFlags(cmd, c.l)
}

func (c *astCmd) Run(_ *cobra.Command, args []string) {
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

	// Do something...
	c.r(p)
}
