package main

import (
	"fmt"
	"os"

	"github.com/object88/cprofile/cmd"
)

// Temporary while working out non-empoymous package issue
var errExitCode = -1

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(errExitCode)
	}
}
