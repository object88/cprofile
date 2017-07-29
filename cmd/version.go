package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version.",
	Long:  "The version of the application.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.0.1")

		if Verbose {
			ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(time.Second))
			defer cancelFn()

			cmd := exec.CommandContext(ctx, "go", "version")

			out, err := cmd.Output()

			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("Attempting to get go version, command timed out")
				return
			}

			if err != nil {
				fmt.Println("Attempting to get go version, got non-zero exit code:", err)
				return
			}

			// If there's no context error, we know the command completed (or errored).
			fmt.Println("Found go:", string(out))
		}
	},
}
