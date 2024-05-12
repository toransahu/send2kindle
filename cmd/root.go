package cmd

import (
	"fmt"
	"os"

	"github.com/toransahu/send2kindle/config"
	"github.com/toransahu/send2kindle/util"
	"github.com/spf13/cobra"
)

func init() {
	var configPath string
	configPath, err := config.DefaultConfigPath()
	if err != nil {
		util.Red.Println("Error setting default config path: ", err)
		os.Exit(1)
	}
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", configPath, "Path to config file")

}

var rootCmd = &cobra.Command{
	Use:   "send2kindle",
	Short: "send2kindle sends documents, webpages and books to your ereader",
	Long: `send2kindle is a CLI tool to send file (books/documents) and webpages to your ereader
It parses the webpage, optimizes it for reading on ereader, and then converts
into an ebook. Then it emails the ebook to the ereader.
Complete documentation is available at https://github.com/toransahu/send2kindle`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		_, err := config.Load(configPath)
		if err != nil {
			util.Red.Println(err)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
