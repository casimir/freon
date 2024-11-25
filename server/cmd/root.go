package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// https://github.com/spf13/cobra/blob/main/site/content/user_guide.md

var rootCmd = &cobra.Command{Use: "freon"}

func init() {
	viper.SetEnvPrefix("freon")
	viper.AutomaticEnv()

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(updatePasswordCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
