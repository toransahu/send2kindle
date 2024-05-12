package cmd

import (
	"github.com/toransahu/send2kindle/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and build info",
	Long:  `Prints the version, platform and build date for send2kindle`,
	Run: func(cmd *cobra.Command, args []string) {
		util.PrintVersion()
	},
}
