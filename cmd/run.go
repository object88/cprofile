package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/google/pprof/profile"
	"github.com/object88/cprofile"
	"github.com/spf13/cobra"
)

var binaryPath string
var profilePath string

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "Processes a profile",
	Long:    "Processes a profile.",
	Run: func(_ *cobra.Command, args []string) {
		// Start by trapping for an interrupt signal.
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)

		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		buffer, err := invokePprof(ctx)
		if err != nil {
			return
		}

		p, profileParseError := profile.ParseData(buffer)
		if profileParseError != nil {
			// Recieved an error loading the pprof data.
			return
		}

		fmt.Printf("Output:\n%s\n", p.String())

		cprofile.NewProcess(p)

		<-sigchan
	},
}

func validateRunFlags(cmd *cobra.Command) error {
	// Consider the binary path.
	return nil
}

func invokePprof(ctx context.Context) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "go", "tool", "pprof", "-proto", binaryPath, profilePath)

	buffer, pprofErr := cmd.Output()
	if pprofErr == context.Canceled {
		return nil, pprofErr
	}

	if pprofErr != nil {
		// Errored out for some other reason; report and bail.
		log.Fatalf("Attempt to get profile information from pprof failed\n%s\n", pprofErr.Error())
		return nil, pprofErr
	}

	return buffer, nil
}
