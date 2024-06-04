package cmd

import (
	"fmt"

	"github.com/casimir/freon/buildinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version string",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", buildinfo.Version)
	},
}
