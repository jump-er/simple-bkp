package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of your CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", version)
	},
}
