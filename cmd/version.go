package cmd

import (
	"context"
	"os/exec"
	"time"

	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version.",
	Long:  "The version of the application.",
	Run: func(_ *cobra.Command, args []string) {
		stdout := cprofile.Stdout()

		stdout.Printf("0.0.1\n")

		if Verbose {
			ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(time.Second))
			defer cancelFn()

			cmd := exec.CommandContext(ctx, "go", "version")

			out, err := cmd.Output()

			if ctx.Err() == context.DeadlineExceeded {
				stderr := cprofile.Stderr()
				stderr.Printf("Attempting to get go version, command timed out\n")
				return
			}

			if err != nil {
				stderr := cprofile.Stderr()
				stderr.Printf("Attempting to get go version, got non-zero exit code: %s\n", err.Error())
				return
			}

			// If there's no context error, we know the command completed (or errored).
			stdout.Printf("Found go: %s\n", string(out))
		}
	},
}
