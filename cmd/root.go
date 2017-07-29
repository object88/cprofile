package cmd

import "github.com/spf13/cobra"

// Verbose describes whether output should be terse or verbose
var Verbose bool

func init() {
	RootCmd.AddCommand(importsCmd, runCmd, versionCmd)
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	runCmd.Flags().StringVarP(&binaryPath, "binaryPath", "b", "", "path to binary")
	runCmd.Flags().StringVarP(&profilePath, "profilePath", "p", "", "path to pprof output file")
}

// RootCmd is the main action taken by Cobra
var RootCmd = &cobra.Command{
	Use:   "cprofile",
	Short: "cprofile injests and processes profile information from pprof",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}
