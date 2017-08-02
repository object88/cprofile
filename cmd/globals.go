package cmd

import (
	"context"
	"fmt"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

var globalsCmd = &cobra.Command{
	Use:   "globals",
	Short: "Returns list of instances of global variables.",
	Long:  "Returns the list of global variables for a program, with file name and offsets.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		base := "."
		if len(args) > 0 {
			base = args[0]
		}

		l := cprofile.NewLoader()
		p, err := l.Load(ctx, base)
		if err != nil {
			fmt.Printf("Got error: %s\n", err.Error())
			return
		}

		pkg, err := p.Package()
		if err != nil {
			fmt.Printf("Got error: %s\n", err.Error())
			return
		}

		for k, v := range pkg.Globals() {
			fmt.Printf("%d: %s\n", k, v)
		}
	},
}
