package main

import (
	"fmt"
	"os"

	"github.com/object88/cprofile/cmd"
)

// Temporary while working out non-empoymous package issue
var errExitCode = -1

func main() {
	rootCmd := cmd.InitializeCommands()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(errExitCode)
	}
}
