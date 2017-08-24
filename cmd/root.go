package cmd

import (
	"github.com/spf13/cobra"
)

// InitializeCommands sets up the cobra commands
func InitializeCommands() *cobra.Command {
	rootCmd := createRootCommand()

	functionsCmd := createFunctionsCommand(rootCmd.o)
	importsCmd := createImportsCommand(rootCmd.o)
	globalsCmd := createGlobalsCommand(rootCmd.o)
	versionCmd := createVersionCommand(rootCmd.o)

	rootCmd.cmd.AddCommand(functionsCmd.cmd, globalsCmd.cmd, importsCmd.cmd, versionCmd)
	return rootCmd.cmd
}

type rootCmd struct {
	cmd *cobra.Command
	o   *globalOptions
}

// RootCmd is the main action taken by Cobra
func createRootCommand() *rootCmd {
	o := &globalOptions{}

	cmd := &cobra.Command{
		Use:   "cprofile",
		Short: "cprofile injests and processes profile information from pprof",
		Long:  "",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			o.ProcessFlags(cmd, nil)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	o.AttachFlags(cmd)

	rootCmd := &rootCmd{cmd, o}
	return rootCmd
}
