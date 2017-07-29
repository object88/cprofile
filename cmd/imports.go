package cmd

import (
	"context"
	"fmt"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

var importsCmd = &cobra.Command{
	Use:   "imports",
	Short: "Print the imports.",
	Long:  "Gets the de-duplicated list of imports.",
	Run: func(cmd *cobra.Command, args []string) {
		// i := cprofile.NewImports()

		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		base := "."
		if len(args) > 0 {
			base = args[0]
		}
		// err := i.Read(ctx, base)
		// if err != nil {
		// 	fmt.Printf("AWHELLNAW.\n%s\n", err.Error())
		// }

		// fmt.Printf("***\n")

		l := cprofile.NewLoader()
		_, err := l.Load(ctx, base)
		if err != nil {
			fmt.Printf("Got error: %s\n", err.Error())
		}

		// fmt.Printf("***\n")

		// for _, v := range i.Flatlist() {
		// 	fmt.Println(v)
		// }
	},
}
